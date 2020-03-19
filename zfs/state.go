package zfs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jsirianni/zfsmon/util/file"
	"github.com/jsirianni/zfsmon/zpool"
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

	z.AlertState = lastSavedState.AlertState
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
		// TODO check for file open error
		return newZfs, err
	}

	err = json.Unmarshal(b, &newZfs)
	return newZfs, err
}
