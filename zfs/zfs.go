package zfs

import (
	"os"
	"sync"
	"time"
	"strconv"
	"encoding/json"

	"github.com/jsirianni/zfsmon/alert"
	"github.com/jsirianni/zfsmon/util/logger"

	"github.com/pkg/errors"
)

// Zfs type holds the global configuration for the zfs package
type Zfs struct {
	Hostname string `json:"hostname"`

	DaemonMode bool `json:"daemon_mode"`

	State struct {
		File string `json:"file"`
		lock sync.Mutex `json:"-"`
	} `json:"state"`

	Pools []Zpool `json:"pools"`

	// Alert is a pluggable interface that
	// can accept different systems for notifying
	// users. See alert/alert.go
	Alert       alert.Alert `json:"-"`
	AlertConfig struct {
		NoAlert bool `json:"no_alert"`
	} `json:"alert_config"`

	// devices in this slice have had a sucessful alert triggered
	AlertState map[string]string `json:"alert_state"`

	Log logger.Logger `json:"-"`
}

// Init initilizes the type
func (z *Zfs) Init() error {
	// TODO: validate all params
	z.AlertState = make(map[string]string)

	if z.Log.Configured() == false {
		z.Log.Configure("error")
		z.Log.Info("logging level set to error")
	}

	// purposely ignore errors
    if s, _ := strconv.ParseBool(os.Getenv("ZFSMON_TEST_ALERT")); s {
		z.Log.Configure("trace")
		z.Log.Info("Environment set to ZFSMON_TEST_ALERT true, testing alert and then exiting.")
		return z.Alert.Message("zfsmon test alert")
    }

	return nil
}

// ZFSMon builds an array of zpool objects and performs health checks on them
func (z Zfs) ZFSMon() error {
	if err := z.readState(); err != nil {
		return err
	}

	if z.Log.Level() == "trace" {
		config, err := json.Marshal(z)
		if err != nil {
			return err
		}
		z.Log.Trace("zfsmon config: " + string(config))
	}

	if z.DaemonMode {
		for {
			if err := z.checkPools(); err != nil {
				z.Log.Error(err)
			}

			if err := z.saveStateFile(); err != nil {
				z.Log.Error(err)
			}

			time.Sleep(time.Second * time.Duration(10))
		}
	}

	if err := z.checkPools(); err != nil {
		if e := z.saveStateFile(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	return z.saveStateFile()
}
