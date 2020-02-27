package models

import (
	"fmt"
	"regexp"
	"strconv"
)

type ConfigVersion struct {
	major, minor uint64
}

var versionRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile(`^"([0-9]+)\.([0-9]+)"$`)
}

func (v *ConfigVersion) String() string {
	return fmt.Sprintf(`%d.%d`, v.major, v.minor)
}

func (v *ConfigVersion) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, v.String())), nil
}

func (v *ConfigVersion) UnmarshalJSON(data []byte) error {
	strVersion := string(data)
	m := versionRegex.FindStringSubmatch(strVersion)
	if m == nil {
		return ErrInvalidVersionFormat
	}

	v.major, _ = strconv.ParseUint(m[1], 10, 64)
	if m[2] != "" {
		v.minor, _ = strconv.ParseUint(m[2], 10, 64)
	}

	return nil
}

func (v *ConfigVersion) MustEqualTo(v2Str string) (bool, error) {
	v2, err := parseVersion(v2Str)
	if err != nil {
		return false, err
	}
	return v.major == v2.major && v.minor == v2.minor, nil
}

func (v *ConfigVersion) MustGreaterThan(v2Str string) (bool, error) {
	v2, err := parseVersion(v2Str)
	if err != nil {
		return false, err
	}
	return v.major > v2.major || (v.major == v2.major && v.minor > v2.minor), nil
}

func (v *ConfigVersion) MustLessThan(v2Str string) (bool, error) {
	v2, err := parseVersion(v2Str)
	if err != nil {
		return false, err
	}
	return v.major < v2.major || (v.major == v2.major && v.minor < v2.minor), nil
}

func (v *ConfigVersion) MustGreaterThanOrEqualTo(v2Str string) (bool, error) {
	v2, err := parseVersion(v2Str)
	if err != nil {
		return false, err
	}
	return v.major > v2.major || (v.major == v2.major && v.minor > v2.minor) || (v.major == v2.major && v.minor == v2.minor), nil
}

func (v *ConfigVersion) MustLessThanOrEqualTo(v2Str string) (bool, error) {
	v2, err := parseVersion(v2Str)
	if err != nil {
		return false, err
	}
	return v.major < v2.major || (v.major == v2.major && v.minor < v2.minor) || (v.major == v2.major && v.minor == v2.minor), nil
}

func parseVersion(version string) (*ConfigVersion, error) {
	v := &ConfigVersion{}
	err := v.UnmarshalJSON([]byte(version))
	if err != nil {
		return nil, err
	}

	return v, nil
}
