package models

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestVersion struct {
	Version ConfigVersion `json:"version"`
}

func TestConfigVersion_MarshalJSON(t *testing.T) {
	for _, testcase := range []struct {
		version     *ConfigVersion
		expectedStr string
	}{
		{
			version:     &ConfigVersion{major: 1, minor: 8},
			expectedStr: "1.8",
		},
		{
			version:     &ConfigVersion{major: 1, minor: 0},
			expectedStr: "1.0",
		},
		{
			version:     &ConfigVersion{},
			expectedStr: "0.0",
		},
	} {
		version := &TestVersion{Version: *testcase.version}

		result, err := json.Marshal(version)
		if assert.NoError(t, err) {
			assert.Equal(t, fmt.Sprintf(`{"version":"%s"}`, testcase.expectedStr), string(result))
		}
	}
}

func TestConfigVersion_UnmarshalJSON(t *testing.T) {
	for _, testcase := range []struct {
		strVersion      string
		expectedVersion *ConfigVersion
		expectedError   error
	}{
		{
			strVersion:      "1.0",
			expectedVersion: &ConfigVersion{major: 1, minor: 0},
		},
		{
			strVersion:      "2.3",
			expectedVersion: &ConfigVersion{major: 2, minor: 3},
		},
		{
			strVersion:      "0.0",
			expectedVersion: &ConfigVersion{},
		},
		{
			strVersion:      "18.3956",
			expectedVersion: &ConfigVersion{major: 18, minor: 3956},
		},
		{
			strVersion:    "1",
			expectedError: ErrInvalidVersionFormat,
		},
		{
			strVersion:    "0.0.1",
			expectedError: ErrInvalidVersionFormat,
		},
		{
			strVersion:    "test",
			expectedError: ErrInvalidVersionFormat,
		},
	} {
		version := &TestVersion{}
		err := json.Unmarshal([]byte(fmt.Sprintf(`{"version":"%s"}`, testcase.strVersion)), &version)
		if testcase.expectedError != nil {
			assert.Equal(t, testcase.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, *testcase.expectedVersion, version.Version)
		}
	}
}

func TestConfigVersion_MustEqualTo(t *testing.T) {
	for _, testcase := range []struct {
		v1, v2 string
		equal  bool
	}{
		{v1: "1.0", v2: "1.0", equal: true},
		{v1: "1.0", v2: "1.1", equal: false},
	} {
		version, _ := parseVersion(testcase.v1)
		result, _ := version.MustEqualTo(testcase.v2)
		assert.Equal(t, testcase.equal, result)
	}
}

func TestConfigVersion_MustGreaterThan(t *testing.T) {
	for _, testcase := range []struct {
		v1, v2  string
		greater bool
	}{
		{v1: "1.0", v2: "1.0", greater: false},
		{v1: "1.0", v2: "1.1", greater: false},
		{v1: "1.0", v2: "2.0", greater: false},
		{v1: "1.0", v2: "0.8", greater: true},
		{v1: "1.1", v2: "1.0", greater: true},
	} {
		version, _ := parseVersion(testcase.v1)
		result, _ := version.MustGreaterThan(testcase.v2)
		assert.Equal(t, testcase.greater, result)
	}
}

func TestConfigVersion_MustGreaterThanOrEqualTo(t *testing.T) {
	for _, testcase := range []struct {
		v1, v2         string
		greaterOrEqual bool
	}{
		{v1: "1.0", v2: "1.0", greaterOrEqual: true},
		{v1: "1.0", v2: "1.1", greaterOrEqual: false},
		{v1: "1.0", v2: "2.0", greaterOrEqual: false},
		{v1: "1.0", v2: "0.8", greaterOrEqual: true},
		{v1: "1.1", v2: "1.0", greaterOrEqual: true},
	} {
		version, _ := parseVersion(testcase.v1)
		result, _ := version.MustGreaterThanOrEqualTo(testcase.v2)
		assert.Equal(t, testcase.greaterOrEqual, result)
	}
}

func TestConfigVersion_MustLessThan(t *testing.T) {
	for _, testcase := range []struct {
		v1, v2 string
		less   bool
	}{
		{v1: "1.0", v2: "1.0", less: false},
		{v1: "1.0", v2: "0.9", less: false},
		{v1: "2.0", v2: "1.0", less: false},
		{v1: "1.0", v2: "1.1", less: true},
		{v1: "1.0", v2: "2.0", less: true},
	} {
		version, _ := parseVersion(testcase.v1)
		result, _ := version.MustLessThan(testcase.v2)
		assert.Equal(t, testcase.less, result)
	}
}

func TestConfigVersion_MustLessThanOrEqualTo(t *testing.T) {
	for _, testcase := range []struct {
		v1, v2      string
		lessOrEqual bool
	}{
		{v1: "1.0", v2: "1.0", lessOrEqual: true},
		{v1: "1.0", v2: "0.9", lessOrEqual: false},
		{v1: "2.0", v2: "1.0", lessOrEqual: false},
		{v1: "1.0", v2: "1.1", lessOrEqual: true},
		{v1: "1.0", v2: "2.0", lessOrEqual: true},
	} {
		version, _ := parseVersion(testcase.v1)
		result, _ := version.MustLessThanOrEqualTo(testcase.v2)
		assert.Equal(t, testcase.lessOrEqual, result)
	}
}

func TestConfigVersion_parseVersion(t *testing.T) {
	version, _ := parseVersion(`"1.8"`)
	assert.Equal(t, uint64(1), version.major)
	assert.Equal(t, uint64(8), version.minor)

	_, err := parseVersion("test")
	assert.Error(t, err)
}
