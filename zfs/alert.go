package zfs

import (
    "github.com/pkg/errors"
)

// IsAlerted returns true if device name is found in the alert state
func (z Zfs) isAlerted(name, state string) bool {
	s, ok := z.AlertState[name]
	if ok {
		// return false if the state has changed
		if state != s {
			return false
		}
	}
	return ok
}

func (z Zfs) sendAlert(pool Zpool, healthy bool) error {
	msg := "host: " + z.Hostname + " | zpool " + pool.Name + " is not in a healthy state, got status: " + pool.State.String()
	if healthy {
		msg = "host: " + z.Hostname + ": zpool " + pool.Name + " is back to a healthy state, got status: " + pool.State.String()
	}

    z.Log.Info(msg)

	if z.AlertConfig.NoAlert {
		z.Log.Info("skipping alert, 'no alert' set to true")
		return nil
	}

	if err := z.Alert.Message(msg); err != nil {
		return errors.Wrap(err, "failed to send alert")
	}
	return nil
}
