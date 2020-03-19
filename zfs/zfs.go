package zfs

import (
	"log"
	"sync"
	"time"

	"github.com/jsirianni/zfsmon/alert"

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

	Pools []Zpool `json:"pools"`

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
