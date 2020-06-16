package templates

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStringFunctions(t *testing.T) {
	tm, err := New("").Parse(`{{join (split . ":") "/"}}`)
	assert.NoError(t, err)

	var b bytes.Buffer
	assert.NoError(t, tm.Execute(&b, "text:with:colon"))
	want := "text/with/colon"
	assert.Equal(t, want, b.String())
}

func TestNew(t *testing.T) {
	tm, err := New("foo").Parse("this is a {{ . }}")
	assert.NoError(t, err)

	var b bytes.Buffer
	assert.NoError(t, tm.Execute(&b, "string"))
	want := "this is a string"
	assert.Equal(t, want, b.String())
}
