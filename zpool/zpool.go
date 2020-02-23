package zpool

import (
	"fmt"
	"encoding/json"

	libzfs "github.com/jsirianni/go-libzfs"
	"github.com/pkg/errors"
)

type Zpool struct {
	Name    string
	State   libzfs.VDevState
	Devices []Device

	// set to true when an alert is triggered
	// for this pool
	Alerted bool
}

type Device struct {
	Name    string
	Type    libzfs.VDevType
	State   libzfs.VDevState
	Devices []Device `json:",omitempty"`
}

// Print a Device object in json or standard output
func (device Device) Print(jsonFmt bool) error {
	if jsonFmt {
		b, err := json.MarshalIndent(device, " ", " ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal Device object for printing in json format")
		}
		fmt.Println(string(b))
		return nil
	}

	fmt.Println("device:", device.Name, device.Type, device.State.String())
	for _, d := range device.Devices {
		fmt.Println("vdev:", d.Name, d.Type, d.State.String())
		for _, s := range d.Devices {
			fmt.Println("vdev:", s.Name, s.Type, s.State.String())
		}
	}
	return nil
}

// Print a Zpool object in json or standard output
func (zpool Zpool) Print(jsonFmt bool) error {
	if jsonFmt {
		b, err := json.MarshalIndent(zpool, " ", " ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal Zpool object for printing in json format")
		}
		fmt.Println(string(b))
		return nil
	}

	fmt.Println("zpool:", zpool.Name, zpool.State.String())
	for _, d := range zpool.Devices {
		fmt.Println("vdev:", d.Name, d.Type, d.State.String())
		for _, s := range d.Devices {
			fmt.Println("vdev:", s.Name, s.Type, s.State.String())
		}
	}
	return nil
}

// RunningPools returns a slice of Zpool objects that are detected
// on the running system
func RunningPools() ([]Zpool, error) {
	globalPools, err := libzfs.PoolOpenAll()
	defer libzfs.PoolCloseAll(globalPools)
	if err != nil {
		return nil, err
	}

	pools := make([]Zpool, len(globalPools))

	for t, pool := range globalPools {
		// get the root vdev (rootDevice.Name will be the zpool name)
		zpool, err := pool.VDevTree()
		if err != nil {
			return pools, err
		} else if zpool.Type != libzfs.VDevTypeRoot {
			return pools, errors.New("ERROR: Expected zpool type to be 'root', got: " + string(zpool.Type))
		}

		// print results of the top level vdev (zpool)
		pools[t].Name = zpool.Name
		pools[t].State = zpool.Stat.State

		// iterate each vdev and display results
		pools[t].Devices = make([]Device, len(zpool.Devices))
		for i, vdev := range zpool.Devices {
			pools[t].Devices[i].Name = vdev.Name
			pools[t].Devices[i].Type = vdev.Type
			pools[t].Devices[i].State = vdev.Stat.State

			if len(vdev.Devices) > 0 {
				pools[t].Devices[i].Devices = make([]Device, len(vdev.Devices))
				for n, disk := range vdev.Devices {
					pools[t].Devices[i].Devices[n].Name = disk.Name
					pools[t].Devices[i].Devices[n].Type = disk.Type
					pools[t].Devices[i].Devices[n].State = disk.Stat.State
				}
			}
		}
	}
	return pools, nil
}
