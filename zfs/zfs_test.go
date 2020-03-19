package zfs

import (
    "testing"
)

func TestIsAlerted(t *testing.T) {
    z := Zfs{}
    z.AlertState = make(map[string]string)
    z.AlertState["a"] = "bad"
    if z.isAlerted("a", "bad") != true {
        t.Errorf("expected isAlerted(a, bad) to return true, got false")
    }

    if z.isAlerted("b", "bad") != false {
        t.Errorf("expected isAlerted(b, bad) to return false, got true")
    }
}

func TestIsAlertedChanged(t *testing.T) {
    z := Zfs{}
    z.AlertState = make(map[string]string)
    z.AlertState["a"] = "bad"
    if z.isAlerted("a", "good") != false {
        t.Errorf("expected isAlerted(a, good) to return false, got true")
    }
}
