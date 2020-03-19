package zfs

import (
    "log"

    multierror "github.com/hashicorp/go-multierror"
    libzfs "github.com/jsirianni/go-libzfs"
)

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

func (z Zfs) checkDevices(p Zpool) (e error) {
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

func (z Zfs) checkDevice(p Zpool, d Device) error {
	if z.Verbose {
		log.Println("checking device in pool: " + p.Name + " " + d.Name + " " + d.State.String())
	}

	if d.State == libzfs.VDevStateHealthy {
		if z.isAlerted(d.Name, d.State.String()) {
			if err := z.sendAlert(p, true); err != nil {
				return err
			}
			// assume the key exists because z.isAlerted returned true
			delete(z.AlertState, d.Name)
		}
		return nil
	}

	// if not healthy, if not alerted, send alert else return
	if z.isAlerted(d.Name, d.State.String()) == false {
		if err := z.sendAlert(p, false); err != nil {
			return err
		}
		z.AlertState[d.Name] = d.State.String()
	}
	return nil
}
