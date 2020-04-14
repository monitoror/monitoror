package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type (
	RawVersion string

	ConfigVersion struct {
		major, minor uint64
	}
)

var versionRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile(`^"([0-9]+)\.([0-9]+)"$`)
}

func (v *ConfigVersion) ToVersion() RawVersion {
	return RawVersion(fmt.Sprintf(`%d.%d`, v.major, v.minor))
}

func (v *ConfigVersion) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%q`, v.ToVersion())), nil
}

func (v *ConfigVersion) UnmarshalJSON(data []byte) error {
	strVersion := string(data)
	m := versionRegex.FindStringSubmatch(strVersion)
	if m == nil {
		return &ConfigVersionFormatError{WrongVersion: strVersion}
	}

	v.major, _ = strconv.ParseUint(m[1], 10, 64)
	if m[2] != "" {
		v.minor, _ = strconv.ParseUint(m[2], 10, 64)
	}

	return nil
}

func (v *ConfigVersion) IsEqualTo(v2Str RawVersion) bool {
	v2 := ParseVersion(v2Str)
	return v.major == v2.major && v.minor == v2.minor
}

func (v *ConfigVersion) IsGreaterThan(v2Str RawVersion) bool {
	v2 := ParseVersion(v2Str)
	return v.major > v2.major || (v.major == v2.major && v.minor > v2.minor)
}

func (v *ConfigVersion) IsLessThan(v2Str RawVersion) bool {
	v2 := ParseVersion(v2Str)
	return v.major < v2.major || (v.major == v2.major && v.minor < v2.minor)
}

func (v *ConfigVersion) IsGreaterThanOrEqualTo(v2Str RawVersion) bool {
	v2 := ParseVersion(v2Str)
	return v.major > v2.major || (v.major == v2.major && v.minor > v2.minor) || (v.major == v2.major && v.minor == v2.minor)
}

func (v *ConfigVersion) IsLessThanOrEqualTo(v2Str RawVersion) bool {
	v2 := ParseVersion(v2Str)
	return v.major < v2.major || (v.major == v2.major && v.minor < v2.minor) || (v.major == v2.major && v.minor == v2.minor)
}

func ParseVersion(version RawVersion) *ConfigVersion {
	// Hack to use "X.Y" in test or code instead of "\"X.Y"\"
	version = RawVersion(strings.ReplaceAll(fmt.Sprintf(`"%s"`, version), `""`, `"`))
	v := &ConfigVersion{}
	err := v.UnmarshalJSON([]byte(version))
	if err != nil {
		panic(fmt.Sprintf(`Invalid version format used in configversion comparaison function %q, %v`, version, err))
	}

	return v
}
