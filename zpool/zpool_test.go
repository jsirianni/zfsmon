package zpool

import (
    "testing"
)

func TestSortedVdevs(t *testing.T) {
    z := Zpool{}
    z.Devices = append(z.Devices, Device{Name: "b"})
    z.Devices = append(z.Devices, Device{Name: "z"})
    z.Devices = append(z.Devices, Device{Name: "a"})

    sorted := z.SortedDevices()
    if sorted[0] != "a" {
        t.Errorf("expected first device in sorted slice to be named 'a'")
    }

    if sorted[1] != "b" {
        t.Errorf("expected second device in sorted slice to be named 'b'")
    }

    // last
    if sorted[2] != "z" {
        t.Errorf("expected last device in sorted slize to be named 'z'")
    }
}

func TestSortedDevices(t *testing.T) {
    d := Device{}
    d.Devices = append(d.Devices, Device{Name: "b"})
    d.Devices = append(d.Devices, Device{Name: "z"})
    d.Devices = append(d.Devices, Device{Name: "a"})

    sorted := d.SortedDevices()
    if sorted[0] != "a" {
        t.Errorf("expected first device in sorted slice to be named 'a'")
    }

    if sorted[1] != "b" {
        t.Errorf("expected second device in sorted slice to be named 'b'")
    }

    // last
    if sorted[2] != "z" {
        t.Errorf("expected last device in sorted slize to be named 'z'")
    }
}
