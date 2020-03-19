package zfs

import (
	"log"
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

	Pools []zpool.Zpool `json:"pools"`

	// Alert is a pluggable interface that
	// can accept different systems for notifying
	// users. See alert/alert.go
	Alert       alert.Alert `json:"-"`
	AlertConfig struct {
		NoAlert bool `json:"-"`
	} `json:"-"`

	// devices in this slice have had a sucessful alert triggered
	AlertState map[string]string `json:"alert_state"`
}

// Init initilizes the type
func (z *Zfs) Init() error {
	// TODO: validate all params
	z.AlertState = make(map[string]string)
	return nil
}

// ZFSMon builds an array of zpool objects and performs health checks on them
func (z Zfs) ZFSMon() error {
	if err := z.ReadState(); err != nil {
		return err
	}

	if err := z.checkPools(); err != nil {
		if e := z.SaveStateFile(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	return z.SaveStateFile()
}

// IsAlerted returns true if device name is found in the alert state
func (z Zfs) IsAlerted(name, state string) bool {
	s, ok := z.AlertState[name]
	if ok {
		// return false if the state has changed
		if state != s {
			return false
		}
	}
	return ok
}

// checkPools takes an array of zpool objects and sends alert to slack for
// every pool that is in a bad state
func (z Zfs) checkPools() (e error) {
	for poolIndex, p := range z.Pools {
		for i, d := range p.Devices {
			if d.State != libzfs.VDevStateHealthy {
				if z.IsAlerted(d.Name, d.State.String()) == false {
					if err := z.sendAlert(p, false); err != nil {
						e = multierror.Append(e, err)
					} else {
						dev := z.Pools[poolIndex].Devices[i]
						z.AlertState[dev.Name] = dev.State.String()
					}
				}
				continue
			}

			if d.State == libzfs.VDevStateHealthy {
				if z.IsAlerted(d.Name, d.State.String()) {
					if err := z.sendAlert(p, true); err != nil {
						e = multierror.Append(e, err)
					} else {
						dev := z.Pools[poolIndex].Devices[i]
						_, ok := z.AlertState[dev.Name]
						if ok {
							delete(z.AlertState, dev.Name)
						}
					}
				}
				continue
			}
		}
	}
	return e
}

func (z Zfs) sendAlert(pool zpool.Zpool, healthy bool) error {
	msg := "host: " + z.Hostname + ": zpool " + pool.Name + " is not in a healthy state, got status: " + pool.State.String()
	if healthy {
		msg = "host: " + z.Hostname + ": zpool " + pool.Name + " is back to a healthy state, got status: " + pool.State.String()
	}

	if z.AlertConfig.NoAlert == true {
		log.Println("skipping alert, --no-alert passed.")
		return nil
	}

	if err := z.Alert.Message(msg); err != nil {
		log.Println("DEBUG: SLACK ERROR: " + err.Error())
		return errors.Wrap(err, "failed to send alert")
	}
	return nil
}
