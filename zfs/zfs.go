package zfs

import (
	"log"
	"sync"
	"time"

	"github.com/jsirianni/zfsmon/alert"
	"github.com/jsirianni/zfsmon/zpool"

	multierror "github.com/hashicorp/go-multierror"
	libzfs "github.com/jsirianni/go-libzfs"
	"github.com/pkg/errors"
)

// Zfs type holds the global configuration for the zfs package
type Zfs struct {
	Hostname string `json:"-"`

	DaemonMode bool `json:"-"`
	Verbose    bool `json:"-"`

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

	if z.DaemonMode {
		for {
			if err := z.checkPools(); err != nil {
				log.Println(err)
			}

			if err := z.SaveStateFile(); err != nil {
				log.Println(err)
			}

			time.Sleep(time.Second * time.Duration(10))
		}
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
	for _, p := range z.Pools {
		if err := z.checkDevices(p); err != nil {
			e = multierror.Append(e, err)
		}
	}
	return e
}

func (z Zfs) checkDevices(p zpool.Zpool) (e error) {
	for _, d := range p.Devices {
		t := string(d.Type)
		if ( t == "raidz" || t == "mirror" ) {
			for _, d := range d.Devices {
				if err := z.checkDevice(p, d); err != nil {
					e = multierror.Append(e, err)
				}
			}
		} else {
			if err := z.checkDevice(p, d); err != nil {
				e = multierror.Append(e, err)
			}
		}
	}
	return e
}

func (z Zfs) checkDevice(p zpool.Zpool, d zpool.Device) error {
	if z.Verbose {
		log.Println("checking device in pool: " + p.Name + " " + d.Name + " " + d.State.String())
	}

	if d.State == libzfs.VDevStateHealthy {
		if z.IsAlerted(d.Name, d.State.String()) {
			if err := z.sendAlert(p, true); err != nil {
				return err
			}
			// assume the key exists because z.IsAlerted returned true
			delete(z.AlertState, d.Name)
		}
		return nil
	}

	// if not healthy, if not alerted, send alert else return
	if z.IsAlerted(d.Name, d.State.String()) == false {
		if err := z.sendAlert(p, false); err != nil {
			return err
		}
		z.AlertState[d.Name] = d.State.String()
	}
	return nil
}

func (z Zfs) sendAlert(pool zpool.Zpool, healthy bool) error {
	msg := "host: " + z.Hostname + ": zpool " + pool.Name + " is not in a healthy state, got status: " + pool.State.String()
	if healthy {
		msg = "host: " + z.Hostname + ": zpool " + pool.Name + " is back to a healthy state, got status: " + pool.State.String()
	}

	if z.DaemonMode {
		log.Println(msg)
	}

	if z.AlertConfig.NoAlert == true {
		log.Println("skipping alert, --no-alert passed.")
		return nil
	}

	if err := z.Alert.Message(msg); err != nil {
		return errors.Wrap(err, "failed to send alert")
	}
	return nil
}
