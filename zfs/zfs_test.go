package zfs

import (
    "testing"
)

func TestIsAlerted(t *testing.T) {
    z := Zfs{}
    z.AlertState = make(map[string]string)
    z.AlertState["a"] = "bad"
    if z.IsAlerted("a", "bad") != true {
        t.Errorf("expected IsAlerted(a, bad) to return true, got false")
    }

    if z.IsAlerted("b", "bad") != false {
        t.Errorf("expected IsAlerted(b, bad) to return false, got true")
    }
}

func TestIsAlertedChanged(t *testing.T) {
    z := Zfs{}
    z.AlertState = make(map[string]string)
    z.AlertState["a"] = "bad"
    if z.IsAlerted("a", "good") != false {
        t.Errorf("expected IsAlerted(a, good) to return false, got true")
    }
}
