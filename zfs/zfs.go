package zfs

import (
	"fmt"

	"github.com/jsirianni/zfsmon/util/alert"
	"github.com/jsirianni/zfsmon/util/file"

	multierror "github.com/hashicorp/go-multierror"
	libzfs "github.com/jsirianni/go-libzfs"
)

// Zfs type holds the global configuration for the zfs package
type Zfs struct {
	HookURL      string
	SlackChannel string
	AlertFile    string
	NoAlert      bool

	JSONOutput  bool

	Pools []Zpool
}

// ZFSMon builds an array of zpool objects and performs health checks on them
func (z *Zfs) ZFSMon() error {
	var err error
	z.Pools, err = MakeSystemReport()
	if err != nil {
		return err
	}

	for _, pool := range z.Pools {
		if err := pool.Print(z.JSONOutput); err != nil {
			return err
		}
	}

	//return z.checkPools(zpools)
	return z.checkPools()
}

// checkPools takes an array of zpool objects and sends alert to slack for
// every pool that is in a bad state
func (z Zfs) checkPools() error {
	// all errors will be collected with 'go-multierror' and returned at the
	// end of this function
	var e error

	for _, p := range z.Pools {

		// if zpool is not healthy, send an alert and write the pool name to
		// the alert file
		if p.State != libzfs.VDevStateHealthy {
			p.Print(z.JSONOutput)
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
func (z Zfs) sendAlert(pool Zpool) error {
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
