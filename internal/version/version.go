package version

import "fmt"

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) Compare(other Version) int {
	if v.Major > other.Major {
		return 1
	}
	if v.Major < other.Major {
		return -1
	}
	if v.Minor > other.Minor {
		return 1
	}
	if v.Minor < other.Minor {
		return -1
	}
	if v.Patch > other.Patch {
		return 1
	}
	if v.Patch < other.Patch {
		return -1
	}
	return 0
}

func (v Version) IsGreaterThan(other Version) bool {
	return v.Compare(other) > 0
}

func (v Version) IsLessThan(other Version) bool {
	return v.Compare(other) < 0
}

func (v Version) IsEqualTo(other Version) bool {
	return v.Compare(other) == 0
}

func (v Version) IsGreaterThanOrEqualTo(other Version) bool {
	return v.Compare(other) >= 0
}

func (v Version) IsLessThanOrEqualTo(other Version) bool {
	return v.Compare(other) <= 0
}

func (v Version) IsBetween(min, max Version) bool {
	return v.IsGreaterThanOrEqualTo(min) && v.IsLessThanOrEqualTo(max)
}

func (v Version) IsNotBeetween(min, max Version) bool {
	return !v.IsBetween(min, max)
}

func (v Version) IsBetweenExclusive(min, max Version) bool {
	return v.IsGreaterThan(min) && v.IsLessThan(max)
}

func (v Version) IsNotBetweenExclusive(min, max Version) bool {
	return !v.IsBetweenExclusive(min, max)
}

func (v Version) ToString() string {
	if v.Major >= 11 {
		// adobe series
		return fmt.Sprintf("Adobe Director %d.%d.%d", v.Major, v.Minor, v.Patch)
	}

	if v.Major == 10 {
		// they added "2004" to the version string to differentiate it from version 9
		return fmt.Sprintf("Macromedia Director MX 2004 (%d.%d.%d)", v.Major, v.Minor, v.Patch)
	}

	if v.Major == 9 {
		return fmt.Sprintf("Macromedia Director MX (%d.%d.%d)", v.Major, v.Minor, v.Patch)
	}

	return fmt.Sprintf("Macromedia Director %d.%d.%d", v.Major, v.Minor, v.Patch)
}

var Director_12_0_0 = Version{12, 0, 0}
var Director_11_5_0 = Version{11, 5, 0}
var Director_11_0_0 = Version{11, 0, 0}
var Director_10_0_0 = Version{10, 0, 0}
var Director_8_5_0 = Version{8, 5, 0}
var Director_8_0_0 = Version{8, 0, 0}
var Director_7_0_0 = Version{7, 0, 0}
var Director_6_0_0 = Version{6, 0, 0}
var Director_5_0_0 = Version{5, 0, 0}
var Director_4_0_4 = Version{4, 0, 4}
var Director_4_0_0 = Version{4, 0, 0}
var Director_3_1_0 = Version{3, 1, 0}
var Director_3_0_0 = Version{3, 0, 0}
var Director_2_0_0 = Version{2, 0, 0}
var Director_1_0_0 = Version{1, 0, 0}

func ParseVersion(ver int32) Version {

	if ver >= 1951 {
		return Director_12_0_0
	}
	if ver >= 1922 {
		return Director_11_5_0
	}
	if ver >= 1921 {
		return Director_11_0_0
	}
	if ver >= 1851 {
		return Director_10_0_0
	}
	if ver >= 1700 {
		return Director_8_5_0
	}
	if ver >= 1410 {
		return Director_8_0_0
	}
	if ver >= 1224 {
		return Director_7_0_0
	}
	if ver >= 1218 {
		return Director_6_0_0
	}
	if ver >= 1201 {
		return Director_5_0_0
	}
	if ver >= 1117 {
		return Director_4_0_4
	}
	if ver >= 1115 {
		return Director_4_0_0
	}
	if ver >= 1029 {
		return Director_3_1_0
	}
	if ver >= 1028 {
		return Director_3_0_0
	}

	return Director_2_0_0
}
