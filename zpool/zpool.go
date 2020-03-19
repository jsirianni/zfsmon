package zpool

import (
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
