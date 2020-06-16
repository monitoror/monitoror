package structs

import (
	"testing"

	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

func TestGetJSONFieldName(t *testing.T) {
	f := &struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2,omitempty"`
	}{}

	fields := structs.Fields(f)
	assert.Equal(t, "field1", GetJSONFieldName(fields[0]))
	assert.Equal(t, "field2", GetJSONFieldName(fields[1]))
}

func TestGetQueryFieldName(t *testing.T) {
	f := &struct {
		Field1 string `query:"field1"`
		Field2 string `query:"field2,omitempty"`
	}{}

	fields := structs.Fields(f)
	assert.Equal(t, "field1", GetQueryFieldName(fields[0]))
	assert.Equal(t, "field2", GetQueryFieldName(fields[1]))
}
