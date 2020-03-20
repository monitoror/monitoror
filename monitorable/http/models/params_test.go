package models

import (
	"regexp"
	"testing"

	. "github.com/monitoror/monitoror/pkg/monitoror/utils"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestHTTPParams_IsValid(t *testing.T) {
	for _, testcase := range []struct {
		params   Validator
		expected bool
	}{
		{&HTTPStatusParams{}, false},
		{&HTTPStatusParams{URL: "toto"}, true},
		{&HTTPStatusParams{URL: "toto", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPStatusParams{URL: "toto", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},

		{&HTTPRawParams{}, false},
		{&HTTPRawParams{URL: "toto"}, true},
		{&HTTPRawParams{URL: "toto", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPRawParams{URL: "toto", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HTTPRawParams{URL: "toto", Regex: "("}, false},
		{&HTTPRawParams{URL: "toto", Regex: "(.*)"}, true},

		{&HTTPFormattedParams{}, false},
		{&HTTPFormattedParams{URL: "toto"}, false},
		{&HTTPFormattedParams{URL: "toto", Format: "unknown"}, false},
		{&HTTPFormattedParams{URL: "toto", Format: "JSON", Key: ""}, false},
		{&HTTPFormattedParams{URL: "toto", Format: "JSON", Key: "."}, false},
		{&HTTPFormattedParams{URL: "toto", Format: "JSON", Key: "key"}, true},
		{&HTTPFormattedParams{URL: "toto", Format: "JSON", Key: "key", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPFormattedParams{URL: "toto", Format: "JSON", Key: "key", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HTTPFormattedParams{URL: "toto", Format: "JSON", Key: "key", Regex: "("}, false},
		{&HTTPFormattedParams{URL: "toto", Format: "JSON", Key: "key", Regex: "(.*)"}, true},
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

		{&HTTPFormattedParams{}, "", nil},
		{&HTTPFormattedParams{Regex: ""}, "", nil},
		{&HTTPFormattedParams{Regex: "("}, "(", nil},
		{&HTTPFormattedParams{Regex: "(.*)"}, "(.*)", regexp.MustCompile("(.*)")},
	} {
		assert.Equal(t, testcase.expectedRegex, testcase.params.GetRegex())
		if isValidRegex(testcase.params) {
			assert.Equal(t, testcase.expectedRegexp, testcase.params.GetRegexp())
		}
	}
}

func TestHTTPFormattedParams_FormattedDataProvider(t *testing.T) {
	for _, testcase := range []struct {
		params         FormattedDataProvider
		expectedFormat string
		expectedKey    string
	}{
		{&HTTPFormattedParams{}, "", ""},
		{&HTTPFormattedParams{Format: JSONFormat}, JSONFormat, ""},
		{&HTTPFormattedParams{Format: YAMLFormat, Key: "key"}, YAMLFormat, "key"},
		{&HTTPFormattedParams{Format: XMLFormat, Key: "key"}, XMLFormat, "key"},
	} {
		assert.Equal(t, testcase.expectedFormat, testcase.params.GetFormat())
		assert.Equal(t, testcase.expectedKey, testcase.params.GetKey())
	}
}
