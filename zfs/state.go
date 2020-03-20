package zfs

import (
	"encoding/json"

	"github.com/jsirianni/zfsmon/util/file"
)

func (z *Zfs) readState() error {
	var err error
	var lastSavedState Zfs

	z.Pools, err = runningPools()
	if err != nil {
		return err
	}

	lastSavedState, err = z.readStateFile()
	if err != nil {
		z.Log.Error(err)
		return nil
	}

	z.AlertState = lastSavedState.AlertState
	return nil
}

func (z *Zfs) saveStateFile() error {
	z.State.lock.Lock()
	defer z.State.lock.Unlock()

	b, err := json.MarshalIndent(z, " ", " ")
	if err != nil {
		return err
	}
	return file.WriteFile(b, z.State.File)
}

func (z *Zfs) readStateFile() (Zfs, error) {
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
