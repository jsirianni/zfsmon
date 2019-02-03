package cmd

import (
    "fmt"
    "errors"

    zfs "github.com/bicomsystems/go-libzfs-0.2"
    "github.com/jsirianni/slacklib/slacklib"
)

var hook_url string
var channel string
var daemon bool
var printReport bool
var noAlert bool

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
            zpool.Print()
        }
    }

    for _, zpool := range report {
        zpool.Alert()
    }

    return nil
}

func (zpool *ZpoolReport) Alert() {
    var alert slacklib.SlackPost
    alert.Channel = channel

    if zpool.State != zfs.VDevStateHealthy {
        alert.Text = ("zpool " + zpool.Name + " is not in a healthy state, got: " + string(zpool.State.String()))
    }

    if len(alert.Text) != 0 {
        if noAlert == true {
            fmt.Println("skipping alert, --no-alert passed.")
        } else {
            slacklib.BasicMessage(alert, hook_url)
        }
    }
}

func (zpool *ZpoolReport) Print() {
    fmt.Println("zpool:", zpool.Name, zpool.State.String())
    for _, d := range zpool.Devices {
        fmt.Println("vdev:", d.Name, d.Type, d.State.String())
        for _, s := range d.Devices {
            fmt.Println("vdev:", s.Name, s.Type, s.State.String())
        }
    }
}

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
            // print the vdev, which could be a disk or a raidz object
            //fmt.Println("vdev:", vdev.Name, vdev.Type, vdev.Stat.State)
            report[t].Devices[i].Name = vdev.Name
            report[t].Devices[i].Type = vdev.Type
            report[t].Devices[i].State = vdev.Stat.State


            // if vdev is a raidz object
            if vdev.Type == zfs.VDevTypeRaidz {
                report[t].Devices[i].Devices = make([]Device, len(vdev.Devices))
                for n, disk := range vdev.Devices {
                    //fmt.Println("vdev:", disk.Name, disk.Type, disk.Stat.State)
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
