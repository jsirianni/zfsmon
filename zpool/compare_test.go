package zpool

import (
	"testing"
)

func testPoolA() Zpool {
	p := Zpool{
		Name:  "a",
		State: 7,
	}
	p.Devices = append(p.Devices, testDeviceA())
	return p
}

func testDeviceA() Device {
	return Device{
		Name:  "a",
		Type:  "disk",
		State: 7,
	}
}

func TestComparePoolsEqual(t *testing.T) {
	a := []Zpool{}
	a = append(a, testPoolA())
	a = append(a, testPoolA())
	b := a

	x, err := ComparePools(a, b)
	if err != nil {
		t.Errorf("Expected ComparePools() to return a nil error, got: " + err.Error())
	}

	if x != true {
		t.Errorf("Expected ComparePools() to return true when given identical []Zpool slices")
	}
}

func TestComparePoolsUnEqualLength(t *testing.T) {
	a := []Zpool{}
	a = append(a, testPoolA())
	a = append(a, testPoolA())
	b := a
	b = append(b, testPoolA())

	x, err := ComparePools(a, b)
	if err != nil {
		t.Errorf("Expected ComparePools() to return a nil error, got: " + err.Error())
	}

	if x != false {
		t.Errorf("Expected ComparePools() to return false when given different amount of zpools")
	}
}

func TestComparePoolEqual(t *testing.T) {
	a := testPoolA()
	b := a

	x, err := ComparePool(a, b)
	if err != nil {
		t.Errorf("Expected ComparePool() to return a nil error, got: " + err.Error())
	}

	if x != true {
		t.Errorf("Expected CompareCool() to return true when given identical pools")
	}
}

func TestComparePoolUnEqualName(t *testing.T) {
	a := testPoolA()
	b := a
	b.Name = "b"

	x, err := ComparePool(a, b)
	if err != nil {
		t.Errorf("Expected ComparePool() to return a nil error, got: " + err.Error())
	}

	if x != false {
		t.Errorf("Expected ComparePool() to return false when given pools with different names")
	}
}

func TestComparePoolUnEqualDevices(t *testing.T) {
	a := testPoolA()
	b := a
	b.Devices = append(b.Devices, testDeviceA())

	x, err := ComparePool(a, b)
	if err != nil {
		t.Errorf("Expected ComparePool() to return a nil error, got: " + err.Error())
	}

	if x != false {
		t.Errorf("Expected ComparePool() to return false when given pools with different devices")
	}
}

func TestCompareDevicesEqual(t *testing.T) {
	a := []Device{}
	a = append(a, testDeviceA())
	a = append(a, testDeviceA())
	b := a

	x, err := CompareDevices(a, b)
	if err != nil {
		t.Errorf("Expected CompareDevices() to return a nil error, got: " + err.Error())
		return
	}

	if x != true {
		t.Errorf("Expected CompareDevices() to return true when given identical []Device objects")
	}
}

func TestCompareDevicesNotEqualLength(t *testing.T) {
	a := []Device{}
	b := []Device{}

	a = append(a, testDeviceA())
	a = append(a, testDeviceA())
	b = append(b, testDeviceA())

	x, err := CompareDevices(a, b)
	if err != nil {
		t.Errorf("Expected CompareDevices() to return a nil error, got: " + err.Error())
		return
	}

	if x != false {
		t.Errorf("Expected CompareDevices() to return false when given slices of Device objects that have different lengths")
	}
}

func TestCompareDeviceEqual(t *testing.T) {
	a := testDeviceA()
	b := a

	x, err := CompareDevice(a, b)
	if err != nil {
		t.Errorf("Expected CompareDevice() to return a nil error, got: " + err.Error())
		return
	}

	if x != true {
		t.Errorf("Expected CompareDevices() to return true when given identical devices")
	}
}

func TestCompareDeviceNotEqualName(t *testing.T) {
	a := testDeviceA()
	b := a
	b.Name = "r"

	x, err := CompareDevice(a, b)
	if err != nil {
		t.Errorf("Expected CompareDevice() to return a nil error, got: " + err.Error())
		return
	}

	if x != false {
		t.Errorf("Expected CompareDevices() to return false when given devices with different names")
	}
}

func TestCompareDeviceNotEqualType(t *testing.T) {
	a := testDeviceA()
	b := a
	b.Type = "notadisk"

	x, err := CompareDevice(a, b)
	if err != nil {
		t.Errorf("Expected CompareDevice() to return a nil error, got: " + err.Error())
		return
	}

	if x != false {
		t.Errorf("Expected CompareDevices() to return false when given devices with different types")
	}
}
