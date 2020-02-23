package zfs

import (
	"errors"
	"fmt"
	"encoding/json"

	libzfs "github.com/jsirianni/go-libzfs"
)

type Zpool struct {
	Name    string
	State   libzfs.VDevState
	Devices []Device
}

type Device struct {
	Name    string
	Type    libzfs.VDevType
	State   libzfs.VDevState
	Devices []Device `json:",omitempty"`
}

// MakeSystemReport returns an array of Zpool objects
func MakeSystemReport() ([]Zpool, error) {
	globalPools, err := libzfs.PoolOpenAll()
	defer libzfs.PoolCloseAll(globalPools)
	if err != nil {
		return nil, err
	}

	report := make([]Zpool, len(globalPools))

	for t, pool := range globalPools {
		// get the root vdev (rootDevice.Name will be the zpool name)
		zpool, err := pool.VDevTree()
		if err != nil {
			return report, err
		} else if zpool.Type != libzfs.VDevTypeRoot {
			return report, errors.New("ERROR: Expected zpool type to be 'root', got: " + string(zpool.Type))
		}

		// print results of the top level vdev (zpool)
		report[t].Name = zpool.Name
		report[t].State = zpool.Stat.State

		// iterate each vdev and display results
		report[t].Devices = make([]Device, len(zpool.Devices))
		for i, vdev := range zpool.Devices {
			report[t].Devices[i].Name = vdev.Name
			report[t].Devices[i].Type = vdev.Type
			report[t].Devices[i].State = vdev.Stat.State

			if len(vdev.Devices) > 0 {
				report[t].Devices[i].Devices = make([]Device, len(vdev.Devices))
				for n, disk := range vdev.Devices {
					report[t].Devices[i].Devices[n].Name = disk.Name
					report[t].Devices[i].Devices[n].Type = disk.Type
					report[t].Devices[i].Devices[n].State = disk.Stat.State
				}
			}
		}
	}
	return report, nil
}

// Print displays the zpool health report to standard out
func (zpool *Zpool) Print(jsonFmt bool) error {
	if jsonFmt {
		b, err := json.MarshalIndent(zpool, " ", " ")
		if err != nil {
			return err
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
