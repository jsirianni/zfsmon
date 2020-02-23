package zfs

import (
    "os"
    "fmt"
    "encoding/json"

    "github.com/jsirianni/zfsmon/zpool"
    "github.com/jsirianni/zfsmon/util/file"
)

func (z *Zfs) ReadState() error {
    var err error
    var lastSavedState Zfs

    z.Pools, err = zpool.RunningPools()
	if err != nil {
		return err
	}

    lastSavedState, err = z.ReadStateFile()
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return nil
    }

    s := false
    s, err = zpool.ComparePools(z.Pools, lastSavedState.Pools)
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

func (z *Zfs) SaveStateFile() error {
    z.State.lock.Lock()
    defer z.State.lock.Unlock()

    b, err := json.MarshalIndent(z, " ", " ")
    if err != nil {
        return err
    }
    return file.WriteFile(b, z.State.File)
}

func (z *Zfs) ReadStateFile() (Zfs, error) {
    z.State.lock.Lock()
    defer z.State.lock.Unlock()

    newZfs := Zfs{}
    b, err := file.ReadFile(z.State.File)
    if err != nil {
        return newZfs, err
    }

    err = json.Unmarshal(b, &newZfs)
    return newZfs, err
}
