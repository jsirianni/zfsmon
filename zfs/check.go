package zfs

import (
    multierror "github.com/hashicorp/go-multierror"
    libzfs "github.com/jsirianni/go-libzfs"
)

// checkPools takes an array of zpool objects and sends alert to slack for
// every pool that is in a bad state
func (z Zfs) checkPools() (e error) {
    z.Log.Trace("checking pools")

	for _, p := range z.Pools {
		if err := z.checkDevices(p); err != nil {
			e = multierror.Append(e, err)
		}
	}
	return e
}

func (z Zfs) checkDevices(p Zpool) (e error) {
    z.Log.Trace("checking pool '" + p.Name + "'")

	for _, d := range p.Devices {
		t := string(d.Type)
        z.Log.Trace("device '" + d.Name + "' has type '" + t + "'")

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

func (z Zfs) checkDevice(p Zpool, d Device) error {
    z.Log.Trace("checking device '" + d.Name + "' in pool '" + p.Name + "'")
	if d.State == libzfs.VDevStateHealthy {
        z.Log.Info("device '" + d.Name + "' in pool '" + p.Name + "' is healthy. Status: " + d.State.String())

		if z.isAlerted(d.Name, d.State.String()) {
			if err := z.sendAlert(p, true); err != nil {
				return err
			}
			// assume the key exists because z.isAlerted returned true
			delete(z.AlertState, d.Name)
		}
		return nil
	}


    z.Log.Warning("device '" + d.Name + "' in pool '" + p.Name + "' is not healthy. Status: " + d.State.String())

	// if not healthy, if not alerted, send alert else return
	if z.isAlerted(d.Name, d.State.String()) == false {

		if err := z.sendAlert(p, false); err != nil {
			return err
		}
		z.AlertState[d.Name] = d.State.String()
	}
	return nil
}
