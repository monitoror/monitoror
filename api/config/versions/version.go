package versions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ----------------------------------------------------------------
// ---------------------- AVAILABLE VERSIONS ----------------------
const (
	CurrentVersion = Version2001
	MinimalVersion = Version2000

	Version2000 RawVersion = "2.0" // Initial version
	Version2001 RawVersion = "2.1" // Youtrack
)

// ----------------------------------------------------------------
// ----------------------------------------------------------------

var versionRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile(`^"([0-9]+)\.([0-9]+)"$`)
}

type (
	// RawVersion used to write version in string format
	RawVersion string

	// ConfigVersion used in Config. Store version in structured format to compare version easily
	ConfigVersion struct {
		major, minor uint64
	}
)

func (v RawVersion) ToConfigVersion() *ConfigVersion {
	return parseVersion(&v)
}

func (v ConfigVersion) ToRawVersion() RawVersion {
	return RawVersion(fmt.Sprintf(`%d.%d`, v.major, v.minor))
}

func (v *ConfigVersion) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%q`, v.ToRawVersion())), nil
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
	v2 := parseVersion(&v2Str)
	return v.major == v2.major && v.minor == v2.minor
}

func (v *ConfigVersion) IsGreaterThan(v2Str RawVersion) bool {
	v2 := parseVersion(&v2Str)
	return v.major > v2.major || (v.major == v2.major && v.minor > v2.minor)
}

func (v *ConfigVersion) IsLessThan(v2Str RawVersion) bool {
	v2 := parseVersion(&v2Str)
	return v.major < v2.major || (v.major == v2.major && v.minor < v2.minor)
}

func (v *ConfigVersion) IsGreaterThanOrEqualTo(v2Str RawVersion) bool {
	v2 := parseVersion(&v2Str)
	return v.major > v2.major || (v.major == v2.major && v.minor > v2.minor) || (v.major == v2.major && v.minor == v2.minor)
}

func (v *ConfigVersion) IsLessThanOrEqualTo(v2Str RawVersion) bool {
	v2 := parseVersion(&v2Str)
	return v.major < v2.major || (v.major == v2.major && v.minor < v2.minor) || (v.major == v2.major && v.minor == v2.minor)
}

func parseVersion(version *RawVersion) *ConfigVersion {
	// Hack to use "X.Y" in test or code instead of "\"X.Y"\"
	versionStr := strings.ReplaceAll(fmt.Sprintf(`"%s"`, *version), `""`, `"`)
	v := &ConfigVersion{}
	if err := v.UnmarshalJSON([]byte(versionStr)); err != nil {
		panic(fmt.Sprintf(`Invalid version format used in configversion comparaison function %q, %v`, *version, err))
	}

	return v
}
