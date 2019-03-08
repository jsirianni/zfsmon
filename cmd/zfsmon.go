package cmd

import (
    "os"
    "fmt"
    "errors"
    "bufio"

    "zfsmon/alert"

    zfs "github.com/jsirianni/go-libzfs"
)

var hook_url string
var channel string
var printReport bool
var noAlert bool
var alertFile string

type ZpoolReport struct {
    Name string
    State zfs.VDevState
    Devices []Device
}

type Device struct {
    Name string
    Type zfs.VDevType
    State zfs.VDevState
    Devices []Device
}

func zfsmon() error {
     err := checkFlags()
     if err != nil {
        return err
    }

    report, err := makeSystemReport()
    if err != nil {
        return err
    }

    if printReport == true {
        for _, zpool := range report {
            zpool.print()
        }
    }

    // iterate all zpools
    var zpoolErrors []error
    for _, zpool := range report {

        // if zpool is not healthy, send an alert and write alert to file
        if zpool.State != zfs.VDevStateHealthy {
             if poolAlerted(zpool.Name) == true {
                fmt.Println("already alerted")
            } else {
                err := zpool.zfsAlert()
                if err != nil {
                    zpoolErrors = append(zpoolErrors, err)
                // if alert success
                } else {
                    // write alert to file
                    if err := manageAlertFile(zpool.Name, 0); err != nil {
                        // if write to file fails, add error and alert again next time
                        fmt.Println(err.Error())  // TODO: TEMP
                        zpoolErrors = append(zpoolErrors, err)
                    }
                }
            }

        // if zpool is healthy, remove from file if present, alert if removed
        } else {
            // alert that pool is now healthy if found in alert file
            if poolAlerted(zpool.Name) == true {
                err := zpool.zfsAlert()
                if err != nil {
                    zpoolErrors = append(zpoolErrors, err)
                // alert success, remove from local file
                } else {
                    err := manageAlertFile(zpool.Name, 1)
                    if err != nil {
                        // if remove from file fails, add error and alert again next time
                        zpoolErrors = append(zpoolErrors, err)
                    }
                }
            }
        }
    }

    // if errors, make a big error and return it
    if len(zpoolErrors) != 0 {
        var err string
        for _, e := range zpoolErrors {
            err = e.Error() + "\n"
        }
        return errors.New(err)
    }
    return nil
}

// send an alert to slack
func (zpool *ZpoolReport) zfsAlert() error {
    var a alert.Slack
    a.HookURL = hook_url
    a.Post.Channel = channel
    a.Post.Text = ("zpool " + zpool.Name + " is not in a healthy state, got: " + string(zpool.State.String()))

    if len(a.Post.Text) == 0 {
        return nil
    }

    if noAlert == true {
        fmt.Println("skipping alert, --no-alert passed.")
        return nil
    }

    return a.BasicMessage()
}

// Print displays the zpool health report to standard out
func (zpool *ZpoolReport) print() {
    fmt.Println("zpool:", zpool.Name, zpool.State.String())
    for _, d := range zpool.Devices {
        fmt.Println("vdev:", d.Name, d.Type, d.State.String())
        for _, s := range d.Devices {
            fmt.Println("vdev:", s.Name, s.Type, s.State.String())
        }
    }
}

// makeSystemReport builds an array of ZpooLReports
func makeSystemReport() ([]ZpoolReport, error) {
    globalPools, err := zfs.PoolOpenAll()
    defer zfs.PoolCloseAll(globalPools)
    if err != nil {
        return nil, err
    }

    report := make([]ZpoolReport, len(globalPools))

    for t, pool := range globalPools {
        // get the root vdev (rootDevice.Name will be the zpool name)
        zpool, err := pool.VDevTree()
        if err != nil {
            return report, err
        } else if zpool.Type != zfs.VDevTypeRoot {
            return report, errors.New("ERROR: Expected type to be 'root', got: " + string(zpool.Type))
        }

        // print results of the top level vdev (zpool)
        report[t].Name = zpool.Name
        report[t].State = zpool.Stat.State

        // iterate each vdev and display results
        report[t].Devices = make([]Device, len(zpool.Devices))
        for i, vdev := range zpool.Devices {
            report[t].Devices[i].Name = vdev.Name
            report[t].Devices[i].Type = vdev.Type
            report[t].Devices[i].State = vdev.Stat.State

            if len(vdev.Devices) > 0 {
                report[t].Devices[i].Devices = make([]Device, len(vdev.Devices))
                for n, disk := range vdev.Devices {
                    report[t].Devices[i].Devices[n].Name = disk.Name
                    report[t].Devices[i].Devices[n].Type = disk.Type
                    report[t].Devices[i].Devices[n].State = disk.Stat.State
                }
            }
        }
    }
    return report, nil
}

func poolAlerted(name string) bool {
    alertedZpools, err := readAlertFile()
    if err != nil {
        fmt.Println("failed to read alert file")
        return false
    }

    for _, z := range alertedZpools {
        if z == name {
            return true
        }
    }
    return false
}

// action = 0 adds to file
// action = 1 removes from file
func manageAlertFile(zpoolName string, action int) error {
    // get []string of alerted pool names
    alertedZpools, err := readAlertFile()
    if err != nil {
        return err
    }

    found := -1
    for i, z := range alertedZpools {
        if z == zpoolName {
            found = i
        }
    }

    // if found and action is to add, do nothing
    if found != -1 && action == 0 {
        return nil
    }
    // if not found and action is to add, add to array
    if found == -1 && action == 0 {
        alertedZpools = append(alertedZpools, zpoolName)
    }

    // if not found and action to remove, do nothing
    if found == -1 && action == 1 {
        return nil
    }
    // if found and action is to remove, remove from array
    if found != -1 && action == 1 {
        alertedZpools = removeFromArray(alertedZpools, found)
    }

    // write array to file here
    if err := writeAlertFile(alertedZpools); err != nil {
        return err
    }


    return nil
}

func removeFromArray(s []string, i int) []string {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func readAlertFile() ([]string, error) {
    var a []string

    f, err := os.OpenFile(alertFile, os.O_RDONLY|os.O_CREATE, 0600)
    if err != nil {
        return a, err
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        a = append(a, scanner.Text())
    }
    return a, nil
}

func writeAlertFile(a []string) error {
    f, err := os.Create(alertFile)
    if err != nil {
        return err
    }
    defer f.Close()

    w := bufio.NewWriter(f)
    for _, alert := range a {
        fmt.Fprintln(w, alert)
    }
    return w.Flush()
}

func checkFlags() error {
    if noAlert == true {
        return nil
    }
    if len(channel) == 0 {
        return errors.New("You must pass a channel '--channel <channel_name>'")
    }
    if len(hook_url) == 0 {
        return errors.New("You must pass a slack hook url '--url <hook url>'")
    }
    return nil
}
