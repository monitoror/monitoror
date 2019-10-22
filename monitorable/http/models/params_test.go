package models

import (
	"encoding/json"
	"reflect"
	"regexp"
	"runtime"
	"testing"

	. "github.com/monitoror/monitoror/pkg/monitoror/utils"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestHttpParams_IsValid(t *testing.T) {
	for _, testcase := range []struct {
		params   Validator
		expected bool
	}{
		{&HttpAnyParams{}, false},
		{&HttpAnyParams{Url: "toto"}, true},
		{&HttpAnyParams{Url: "toto", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HttpAnyParams{Url: "toto", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},

		{&HttpRawParams{}, false},
		{&HttpRawParams{Url: "toto"}, true},
		{&HttpRawParams{Url: "toto", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HttpRawParams{Url: "toto", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HttpRawParams{Url: "toto", Regex: "("}, false},
		{&HttpRawParams{Url: "toto", Regex: "(.*)"}, true},

		{&HttpJsonParams{}, false},
		{&HttpJsonParams{Url: "toto"}, false},
		{&HttpJsonParams{Url: "toto", Key: "."}, false},
		{&HttpJsonParams{Url: "toto", Key: ".key"}, true},
		{&HttpJsonParams{Url: "toto", Key: ".key", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HttpJsonParams{Url: "toto", Key: ".key", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HttpJsonParams{Url: "toto", Key: ".key", Regex: "("}, false},
		{&HttpJsonParams{Url: "toto", Key: ".key", Regex: "(.*)"}, true},

		{&HttpYamlParams{}, false},
		{&HttpYamlParams{Url: "toto"}, false},
		{&HttpYamlParams{Url: "toto", Key: "."}, false},
		{&HttpYamlParams{Url: "toto", Key: ".key"}, true},
		{&HttpYamlParams{Url: "toto", Key: ".key", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HttpYamlParams{Url: "toto", Key: ".key", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HttpYamlParams{Url: "toto", Key: ".key", Regex: "("}, false},
		{&HttpYamlParams{Url: "toto", Key: ".key", Regex: "(.*)"}, true},
	} {
		assert.Equal(t, testcase.expected, testcase.params.IsValid())
	}
}

func TestHttpParams_GetRegex(t *testing.T) {
	for _, testcase := range []struct {
		params         RegexProvider
		expectedRegex  string
		expectedRegexp *regexp.Regexp
	}{
		{&HttpRawParams{}, "", nil},
		{&HttpRawParams{Regex: ""}, "", nil},
		{&HttpRawParams{Regex: "("}, "(", nil},
		{&HttpRawParams{Regex: "(.*)"}, "(.*)", regexp.MustCompile("(.*)")},

		{&HttpJsonParams{}, "", nil},
		{&HttpJsonParams{Regex: ""}, "", nil},
		{&HttpJsonParams{Regex: "("}, "(", nil},
		{&HttpJsonParams{Regex: "(.*)"}, "(.*)", regexp.MustCompile("(.*)")},

		{&HttpYamlParams{}, "", nil},
		{&HttpYamlParams{Regex: ""}, "", nil},
		{&HttpYamlParams{Regex: "("}, "(", nil},
		{&HttpYamlParams{Regex: "(.*)"}, "(.*)", regexp.MustCompile("(.*)")},
	} {
		assert.Equal(t, testcase.expectedRegex, testcase.params.GetRegex())
		if isValidRegex(testcase.params) {
			assert.Equal(t, testcase.expectedRegexp, testcase.params.GetRegexp())
		}
	}
}

func TestHttpJsonParams_FormatedDataProvider(t *testing.T) {
	for _, testcase := range []struct {
		params               FormatedDataProvider
		expectedKey          string
		expectedUnmarshaller func(data []byte, v interface{}) error
	}{
		{&HttpJsonParams{}, "", json.Unmarshal},
		{&HttpJsonParams{Key: ".key"}, ".key", json.Unmarshal},

		{&HttpYamlParams{}, "", yaml.Unmarshal},
		{&HttpYamlParams{Key: ".key"}, ".key", yaml.Unmarshal},
	} {
		assert.Equal(t, testcase.expectedKey, testcase.params.GetKey())

		// Tricks for testing 2 functions. See : https://github.com/stretchr/testify/issues/182#issuecomment-495359313
		funcName1 := runtime.FuncForPC(reflect.ValueOf(testcase.expectedUnmarshaller).Pointer()).Name()
		funcName2 := runtime.FuncForPC(reflect.ValueOf(testcase.params.GetUnmarshaller()).Pointer()).Name()
		assert.Equal(t, funcName1, funcName2)
	}
}
