package terminal

import (
    "github.com/jsirianni/zfsmon/util/logger"
)

type Terminal struct {

}

func (t Terminal) Message(message string) error {
    log := logger.Logger{}
    if err := log.Configure(logger.InfoLVL); err != nil {
        return err
    }
    log.Info(message)
    return nil
}
