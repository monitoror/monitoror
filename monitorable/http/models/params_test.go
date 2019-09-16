package models

import (
	"testing"

	"github.com/AlekSi/pointer"

	"github.com/stretchr/testify/assert"
)

func TestHttpParams_IsValid(t *testing.T) {
	httpAny := &HttpAnyParams{}
	assert.False(t, httpAny.IsValid())
	httpAny.Url = "toto"
	assert.True(t, httpAny.IsValid())
	httpAny.StatusCodeMin = pointer.ToInt(300)
	httpAny.StatusCodeMax = pointer.ToInt(299)
	assert.False(t, httpAny.IsValid())
	httpAny.StatusCodeMin = pointer.ToInt(200)
	httpAny.StatusCodeMax = pointer.ToInt(299)
	assert.True(t, httpAny.IsValid())

	httpRaw := &HttpRawParams{}
	assert.False(t, httpRaw.IsValid())
	httpRaw.Url = "toto"
	assert.True(t, httpRaw.IsValid())
	httpRaw.Regex = "("
	assert.False(t, httpRaw.IsValid())
	httpRaw.Regex = "(.*)"
	assert.True(t, httpRaw.IsValid())
	httpRaw.StatusCodeMin = pointer.ToInt(300)
	httpRaw.StatusCodeMax = pointer.ToInt(299)
	assert.False(t, httpRaw.IsValid())
	httpRaw.StatusCodeMin = pointer.ToInt(200)
	httpRaw.StatusCodeMax = pointer.ToInt(299)
	assert.True(t, httpRaw.IsValid())

	httpFormattedData := &HttpFormattedDataParams{}
	assert.False(t, httpFormattedData.IsValid())
	httpFormattedData.Key = ".bloc1"
	assert.False(t, httpFormattedData.IsValid())
	httpFormattedData.Url = "toto"
	assert.True(t, httpFormattedData.IsValid())
	httpFormattedData.Regex = "("
	assert.False(t, httpFormattedData.IsValid())
	httpFormattedData.Regex = "(.*)"
	assert.True(t, httpFormattedData.IsValid())
	httpFormattedData.StatusCodeMin = pointer.ToInt(300)
	httpFormattedData.StatusCodeMax = pointer.ToInt(299)
	assert.False(t, httpFormattedData.IsValid())
	httpFormattedData.StatusCodeMin = pointer.ToInt(200)
	httpFormattedData.StatusCodeMax = pointer.ToInt(299)
	assert.True(t, httpFormattedData.IsValid())
}

func TestHttpParams_GetRegex(t *testing.T) {
	httpRaw := &HttpRawParams{}
	assert.Nil(t, httpRaw.GetRegex())
	httpRaw.Regex = "(.*)"
	assert.NotNil(t, httpRaw.GetRegex())

	httpFormattedData := &HttpFormattedDataParams{}
	assert.Nil(t, httpFormattedData.GetRegex())
	httpFormattedData.Regex = "(.*)"
	assert.NotNil(t, httpFormattedData.GetRegex())
}
