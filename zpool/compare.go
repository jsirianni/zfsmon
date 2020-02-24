package zpool

func ComparePools(a, b []Zpool) (bool, error) {
	if len(a) != len(b) {
		return false, nil
	}

	for i, _ := range a {
		x, err := ComparePool(a[i], b[i])
		if err != nil {
			return false, err
		}
		if x == false {
			return false, nil
		}
	}

	return true, nil
}

func ComparePool(a, b Zpool) (bool, error) {
	if a.Name != b.Name {
		return false, nil
	}

	return CompareDevices(a.Devices, b.Devices)
}

func CompareDevices(a, b []Device) (bool, error) {
	if len(a) != len(b) {
		return false, nil
	}

	for i, _ := range a {
		x, err := CompareDevice(a[i], b[i])
		if err != nil {
			return false, err
		}
		if x == false {
			return false, err
		}
	}

	return true, nil
}

func CompareDevice(a, b Device) (bool, error) {
	if a.Name != b.Name {
		return false, nil
	}

	if a.Type != b.Type {
		return false, nil
	}

	// return if number of sub devices is not eqaul
	if len(a.Devices) != len(b.Devices) {
		return false, nil
	}

	// we blindly assume that the index will work for both
	// a.Devices and b.Devices because the previous length check
	// passed above ^ ^ ^
	for i, _ := range a.Devices {
		if a.Devices[i].Name != b.Devices[i].Name {
			return false, nil
		}

		if a.Devices[i].Type != b.Devices[i].Type {
			return false, nil
		}
	}

	return true, nil
}
