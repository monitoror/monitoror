package models

import (
	"regexp"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestHTTPParams(t *testing.T) {
	for _, testcase := range []struct {
		params     params.Validator
		errorCount int
	}{
		{&HTTPStatusParams{}, 1},
		{&HTTPStatusParams{URL: "example.com"}, 1},
		{&HTTPStatusParams{URL: "http%sexample.com"}, 1},
		{&HTTPStatusParams{URL: "http://example.com"}, 0},
		{&HTTPStatusParams{URL: "http://example.com", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, 1},
		{&HTTPStatusParams{URL: "http://example.com", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, 0},

		{&HTTPRawParams{}, 1},
		{&HTTPRawParams{URL: "http://example.com"}, 0},
		{&HTTPRawParams{URL: "http://example.com", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, 1},
		{&HTTPRawParams{URL: "http://example.com", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, 0},
		{&HTTPRawParams{URL: "http://example.com", Regex: "("}, 1},
		{&HTTPRawParams{URL: "http://example.com", Regex: "(.*)"}, 0},

		{&HTTPFormattedParams{}, 3},
		{&HTTPFormattedParams{URL: "http://example.com"}, 2},
		{&HTTPFormattedParams{URL: "http://example.com", Format: "unknown"}, 2},
		{&HTTPFormattedParams{URL: "http://example.com", Format: "JSON", Key: ""}, 1},
		{&HTTPFormattedParams{URL: "http://example.com", Format: "JSON", Key: "."}, 1},
		{&HTTPFormattedParams{URL: "http://example.com", Format: "JSON", Key: "key"}, 0},
		{&HTTPFormattedParams{URL: "http://example.com", Format: "JSON", Key: "key", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, 1},
		{&HTTPFormattedParams{URL: "http://example.com", Format: "JSON", Key: "key", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, 0},
		{&HTTPFormattedParams{URL: "http://example.com", Format: "JSON", Key: "key", Regex: "("}, 1},
		{&HTTPFormattedParams{URL: "http://example.com", Format: "JSON", Key: "key", Regex: "(.*)"}, 0},
	} {
		test.AssertParams(t, testcase.params, testcase.errorCount)
		if testcase.errorCount == 0 {
			assert.NotEmpty(t, testcase.params.(GenericParamsProvider).GetURL())
		}
	}
}

func TestHTTPParams_GetRegex(t *testing.T) {
	for _, testcase := range []struct {
		params         RegexParamsProvider
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
		assert.Equal(t, testcase.expectedRegexp, testcase.params.GetRegexp())
	}
}

func TestHTTPFormattedParams_FormattedDataProvider(t *testing.T) {
	for _, testcase := range []struct {
		params         FormattedParamsProvider
		expectedFormat Format
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
