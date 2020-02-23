package zfs

import (
    "os"
    "fmt"
    //"reflect"
    "encoding/json"

    "github.com/jsirianni/zfsmon/util/file"
)

func (z *Zfs) ReadState() error {
    var err error
    var lastSavedState Zfs

    z.Pools, err = RunningPools()
	if err != nil {
		return err
	}

    lastSavedState, err = z.ReadStateFile()
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return nil
    }

    s := false
    s, err = comparePools(z.Pools, lastSavedState.Pools)
    if err != nil {
        return err
    }
    if s == false {
        fmt.Println("last saved state file has different pools than the running system, ignoring the state and not merging old alert status")
        return nil
    }

    // merge alert status from last known state to the running configuration
    // to prevent re-alerting
    for _, lastSavedStatePool := range lastSavedState.Pools {
        for i, pool := range z.Pools {
            if lastSavedStatePool.Name == pool.Name {
                z.Pools[i].Alerted = lastSavedStatePool.Alerted
            }
        }

    }

    return nil
}

func (z Zfs) SaveStateFile() error {
    b, err := json.MarshalIndent(z, " ", " ")
    if err != nil {
        return err
    }
    return file.WriteFile(b, z.AlertFile)
}

func (z Zfs) ReadStateFile() (Zfs, error) {
    newZfs := Zfs{}

    b, err := file.ReadFile(z.AlertFile)
    if err != nil {
        return newZfs, err
    }

    err = json.Unmarshal(b, &newZfs)
    return newZfs, err
}

func comparePools(a, b []Zpool) (bool, error) {
    for i, _ := range a {
        if a[i].Name != b[i].Name {
            return false, nil
        }
    }
    return true, nil
}
