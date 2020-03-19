package terminal

import (
    "testing"
)

func TestStdout(t *testing.T) {
    a := Terminal{}
    if err := a.Message("hello"); err != nil {
        t.Errorf("expected terminal.Message to return a nil error, got: " + err.Error())
    }
}
