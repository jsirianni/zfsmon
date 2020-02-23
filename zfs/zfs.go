package zfs

import (
    "fmt"

    "github.com/jsirianni/zfsmon/util/alert"
    "github.com/jsirianni/zfsmon/util/file"
    "github.com/jsirianni/zfsmon/zfs/zpool"

    libzfs "github.com/jsirianni/go-libzfs"
    multierror "github.com/hashicorp/go-multierror"
)

/// Zfs type holds the global configuration for the zfs package
type Zfs struct {
    HookURL      string
    SlackChannel string
    AlertFile    string
    NoAlert      bool
}

// ZFSMon builds an array of zpool objects and performs health checks on them
func (z Zfs) ZFSMon() error {
    // discover zpools found on the system
    zpools, err := zpool.MakeSystemReport()
    if err != nil {
        return err
    }
    return z.checkPools(zpools)
}

// checkPools takes an array of zpool objects and sends alert to slack for
// every pool that is in a bad state
func (z Zfs) checkPools(zpools []zpool.Zpool) error {
    // all errors will be collected with 'go-multierror' and returned at the
    // end of this function
    var e error

    for _, p := range zpools {

        // if zpool is not healthy, send an alert and write the pool name to
        // the alert file
        if p.State != libzfs.VDevStateHealthy {
            p.Print()
            if file.PoolAlerted(p.Name, z.AlertFile) == true {
                fmt.Println(p.Name, "already alerted")
            } else {
                err := z.sendAlert(p)
                if err != nil {
                    e = multierror.Append(e, err)
                } else {
                    err := file.ManageFile(p.Name, 0, z.AlertFile)
                    if err != nil {
                        e = multierror.Append(e, err)
                    }
                }
            }

        // if zpool is healthy, check to see if the pool exists in the alert file.
        // If it does, send an alert notifying that the pool is now healthy and
        // then remove the pool from the alert file
        } else {
            if file.PoolAlerted(p.Name, z.AlertFile) == true {
                err := z.sendAlert(p)
                if err != nil {
                    e = multierror.Append(e, err)
                } else {
                    err := file.ManageFile(p.Name, 1, z.AlertFile)
                    if err != nil {
                        e = multierror.Append(e, err)
                    }
                }
            }
        }
    }

    return e
}

// sendAlert sends a slack alert for a specific zpool, returns nil if z.NoAlert
// is set to true
func (z Zfs) sendAlert(pool zpool.Zpool ) error {
    if z.NoAlert == true {
        fmt.Println("skipping alert, --no-alert passed.")
        return nil
    }

    var a alert.Slack
    a.HookURL = z.HookURL
    a.Post.Channel = z.SlackChannel
    a.Post.Text = ("zpool " + pool.Name + " is not in a healthy state, got: " + string(pool.State.String()))
    return a.BasicMessage()
}
