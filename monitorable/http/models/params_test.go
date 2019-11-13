package models

import (
	"encoding/json"
	"reflect"
	"regexp"
	"runtime"
	"testing"

	. "github.com/monitoror/monitoror/pkg/monitoror/utils"
	"gopkg.in/yaml.v2"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestHTTPParams_IsValid(t *testing.T) {
	for _, testcase := range []struct {
		params   Validator
		expected bool
	}{
		{&HTTPAnyParams{}, false},
		{&HTTPAnyParams{URL: "toto"}, true},
		{&HTTPAnyParams{URL: "toto", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPAnyParams{URL: "toto", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},

		{&HTTPRawParams{}, false},
		{&HTTPRawParams{URL: "toto"}, true},
		{&HTTPRawParams{URL: "toto", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPRawParams{URL: "toto", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HTTPRawParams{URL: "toto", Regex: "("}, false},
		{&HTTPRawParams{URL: "toto", Regex: "(.*)"}, true},

		{&HTTPJsonParams{}, false},
		{&HTTPJsonParams{URL: "toto"}, false},
		{&HTTPJsonParams{URL: "toto", Key: "."}, false},
		{&HTTPJsonParams{URL: "toto", Key: ".key"}, true},
		{&HTTPJsonParams{URL: "toto", Key: ".key", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPJsonParams{URL: "toto", Key: ".key", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HTTPJsonParams{URL: "toto", Key: ".key", Regex: "("}, false},
		{&HTTPJsonParams{URL: "toto", Key: ".key", Regex: "(.*)"}, true},

		{&HTTPYamlParams{}, false},
		{&HTTPYamlParams{URL: "toto"}, false},
		{&HTTPYamlParams{URL: "toto", Key: "."}, false},
		{&HTTPYamlParams{URL: "toto", Key: ".key"}, true},
		{&HTTPYamlParams{URL: "toto", Key: ".key", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPYamlParams{URL: "toto", Key: ".key", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HTTPYamlParams{URL: "toto", Key: ".key", Regex: "("}, false},
		{&HTTPYamlParams{URL: "toto", Key: ".key", Regex: "(.*)"}, true},
	} {
		assert.Equal(t, testcase.expected, testcase.params.IsValid())
	}
}

func TestHTTPParams_GetRegex(t *testing.T) {
	for _, testcase := range []struct {
		params         RegexProvider
		expectedRegex  string
		expectedRegexp *regexp.Regexp
	}{
		{&HTTPRawParams{}, "", nil},
		{&HTTPRawParams{Regex: ""}, "", nil},
		{&HTTPRawParams{Regex: "("}, "(", nil},
		{&HTTPRawParams{Regex: "(.*)"}, "(.*)", regexp.MustCompile("(.*)")},

		{&HTTPJsonParams{}, "", nil},
		{&HTTPJsonParams{Regex: ""}, "", nil},
		{&HTTPJsonParams{Regex: "("}, "(", nil},
		{&HTTPJsonParams{Regex: "(.*)"}, "(.*)", regexp.MustCompile("(.*)")},

		{&HTTPYamlParams{}, "", nil},
		{&HTTPYamlParams{Regex: ""}, "", nil},
		{&HTTPYamlParams{Regex: "("}, "(", nil},
		{&HTTPYamlParams{Regex: "(.*)"}, "(.*)", regexp.MustCompile("(.*)")},
	} {
		assert.Equal(t, testcase.expectedRegex, testcase.params.GetRegex())
		if isValidRegex(testcase.params) {
			assert.Equal(t, testcase.expectedRegexp, testcase.params.GetRegexp())
		}
	}
}

func TestHTTPSerializedDataFileParams_FormatedDataProvider(t *testing.T) {
	for _, testcase := range []struct {
		params               FormatedDataProvider
		expectedKey          string
		expectedUnmarshaller func(data []byte, v interface{}) error
	}{
		{&HTTPJsonParams{}, "", json.Unmarshal},
		{&HTTPJsonParams{Key: ".key"}, ".key", json.Unmarshal},

		{&HTTPYamlParams{}, "", yaml.Unmarshal},
		{&HTTPYamlParams{Key: ".key"}, ".key", yaml.Unmarshal},
	} {
		assert.Equal(t, testcase.expectedKey, testcase.params.GetKey())

		// Tricks for testing 2 functions. See : https://github.com/stretchr/testify/issues/182#issuecomment-495359313
		funcName1 := runtime.FuncForPC(reflect.ValueOf(testcase.expectedUnmarshaller).Pointer()).Name()
		funcName2 := runtime.FuncForPC(reflect.ValueOf(testcase.params.GetUnmarshaller()).Pointer()).Name()
		assert.Equal(t, funcName1, funcName2)
	}
}
