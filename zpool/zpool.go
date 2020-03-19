package zpool

import (
	"encoding/json"
	"fmt"
	"sort"

	libzfs "github.com/jsirianni/go-libzfs"
	"github.com/pkg/errors"
)

type Zpool struct {
	Name    string `json:"name"`
	State   libzfs.VDevState `json:"state"`
	Devices []Device `json:"devices"`
}

type Device struct {
	Name    string `json:"name"`
	Type    libzfs.VDevType `json:"type"`
	State   libzfs.VDevState `json:"state"`
	Devices []Device `json:",omitempty"`
}

// Devices returns a slice of device names for a given pool
// sorted by name
func (z Zpool) SortedDevices() []string {
	d := []string{}
	for i := range z.Devices {
		d = append(d, z.Devices[i].Name)
	}
	sort.Sort(sort.StringSlice(d))
	return d
}

func (d Device) SortedDevices() []string {
	dv := []string{}
	for i := range d.Devices {
		dv = append(dv, d.Devices[i].Name)
	}
	sort.Sort(sort.StringSlice(dv))
	return dv
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
