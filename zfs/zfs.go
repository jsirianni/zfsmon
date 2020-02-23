package zfs

import (
	"fmt"

	"github.com/jsirianni/zfsmon/util/alert"

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
	if err := z.ReadState(); err != nil {
		return err
	}

	for _, pool := range z.Pools {
		if err := pool.Print(z.JSONOutput); err != nil {
			return err
		}
	}

	if err := z.checkPools(); err != nil {
		return err
	}

	return z.SaveStateFile()
}

// checkPools takes an array of zpool objects and sends alert to slack for
// every pool that is in a bad state
func (z *Zfs) checkPools() (e error) {
	for i, p := range z.Pools {
		if p.State != libzfs.VDevStateHealthy {
			if p.Alerted == false {
				if err := z.sendAlert(p, false); err != nil {
					e = multierror.Append(e, err)
				} else {
					z.Pools[i].Alerted = true
				}
			}
			continue
		}

		if p.State == libzfs.VDevStateHealthy {
			if p.Alerted == true {
				if err := z.sendAlert(p, true); err != nil {
					e = multierror.Append(e, err)
				} else {
					z.Pools[i].Alerted = false
				}
			}
			continue
		}
	}

	return e
}

// sendAlert sends a slack alert for a specific zpool, returns nil if z.NoAlert
// is set to true
func (z Zfs) sendAlert(pool Zpool, healthy bool) error {
	msg := "zpool " + pool.Name + " is not in a healthy state, got: " + pool.State.String()
	if healthy {
		msg = "zpool " + pool.Name + " is back to a healthy state, got: " + pool.State.String()
	}

	if z.NoAlert == true {
		fmt.Println(msg)
		fmt.Println("skipping alert, --no-alert passed.")
		return nil
	}

	var a alert.Slack
	a.HookURL = z.HookURL
	a.Post.Channel = z.SlackChannel
	a.Post.Text = ("zpool " + pool.Name + " is not in a healthy state, got: " + string(pool.State.String()))
	return a.BasicMessage()
}
