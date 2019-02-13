package cmd

import (
    "fmt"
    "errors"
    "time"

    "zfsmon/alert"

    zfs "github.com/bicomsystems/go-libzfs-0.2"
)

var hook_url string
var channel string
var daemon bool
var printReport bool
var noAlert bool
var checkInt int

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


func run() error {
    for {
        /* zfsmon will
            - build report
            - print if --print is passed
            - alert if --no-alert is NOT passed

            if an error is returned in daemon mode, print the error
            and do not exit
        */
        err := zfsmon()
        if err != nil {
            if daemon == true {
                fmt.Println(err.Error())
            } else {
                return err
            }
        }

        // if daemon mode, sleep and then run again
        // if not daemon mode, return nil
        if daemon == true {
            time.Sleep(time.Duration(checkInt * 60) * time.Second)
        } else {
            return nil
        }
    }
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
            zpool.Print()
        }
    }

    // iterate all zpools
    var zpoolErrors []error
    for _, zpool := range report {
        // if zpool is not healthy
        if zpool.State != zfs.VDevStateHealthy {
            err := zpool.zfsAlert()
            if err != nil {
                zpoolErrors = append(zpoolErrors, err)
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
func (zpool *ZpoolReport) Print() {
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
