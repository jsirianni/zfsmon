package zfs

import (
	"fmt"
	"sync"

	"github.com/jsirianni/zfsmon/alert"
	"github.com/jsirianni/zfsmon/zpool"

	multierror "github.com/hashicorp/go-multierror"
	libzfs "github.com/jsirianni/go-libzfs"
	"github.com/pkg/errors"
)

// Zfs type holds the global configuration for the zfs package
type Zfs struct {
	Hostname string `json:"-"`

	State struct {
		File string `json:"-"`
		lock sync.Mutex `json:"-"`
	} `json:"-"`

	JSONOutput bool `json:"-"`

	Pools []zpool.Zpool `json:"pools,omitempty"`

	// Alert is a pluggable interface that
	// can accept different systems for notifying
	// users. See alert/alert.go
	Alert       alert.Alert `json:"-"`
	AlertConfig struct {
		NoAlert bool `json:"-"`
	} `json:"-"`
}

// ZFSMon builds an array of zpool objects and performs health checks on them
func (z Zfs) ZFSMon() error {
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
func (z Zfs) checkPools() (e error) {
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

func (z Zfs) sendAlert(pool zpool.Zpool, healthy bool) error {
	msg := "host: " + z.Hostname + ": zpool " + pool.Name + " is not in a healthy state, got: " + pool.State.String()
	if healthy {
		msg = "host: " + z.Hostname + ": zpool " + pool.Name + " is back to a healthy state, got: " + pool.State.String()
	}

	if z.AlertConfig.NoAlert == true {
		fmt.Println(msg)
		fmt.Println("skipping alert, --no-alert passed.")
		return nil
	}

	if err := z.Alert.Message(msg); err != nil {
		return errors.Wrap(err, "failed to send alert")
	}
	return nil
}
