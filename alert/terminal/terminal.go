package terminal

import (
    "fmt"
)

type Terminal struct {

}

func (t Terminal) Print() {
    fmt.Println("standard out notifier")
}

func (t Terminal) Message(message string) error {
    fmt.Println("ZFSMON ALERT:", message)
    return nil
}
