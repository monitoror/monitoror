package repository

import (
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigRepository(t *testing.T) {
	assert.NotNil(t, NewConfigRepository())
}

func TestRepository_ReadConfig_Success(t *testing.T) {
	input := `
{
  "columns": 4,
  "apiBaseURL": "localhost:3000",
  "tiles": [
    { "type": "EMPTY" },
    { "type": "PING", "label": "...", "params": { "hostname": "server.com"}},
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" }},
      { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 }}
    ]}
  ]
}
`
	config, err := ReadConfig(strings.NewReader(input))

	assert.NoError(t, err)
	assert.Equal(t, 4, *config.Columns)
}

func TestRepository_ReadConfig_Error_WrongJson(t *testing.T) {
	input := `
{
  "columns": 4,
  "apiBaseURL": "localhost:3000",
  "tiles": [
    { "type": "EMPTY" },
    { "type": "PING", "label": "...", "params": { "hostname": "server.com"}},
    xxxx
  ]
}
`
	_, err := ReadConfig(strings.NewReader(input))

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid character 'x' looking for beginning of value")
}

func TestRepository_ReadConfig_Error_WrongJson2(t *testing.T) {
	input := `
{
  "columns": "4",
  "apiBaseURL": "localhost:3000",
  "tiles": [
    { "type": "EMPTY" },
    { "type": "PING", "label": "...", "params": { "hostname": "server.com"}}
  ]
}
`
	_, err := ReadConfig(strings.NewReader(input))

	assert.Error(t, err)
	assert.EqualError(t, err, "json: cannot unmarshal string into Go struct field Config.columns of type int")
}

func TestRepository_ReadConfig_Error_WrongReader(t *testing.T) {
	_, err := ReadConfig(iotest.TimeoutReader(strings.NewReader("blabla")))

	assert.Error(t, err)
	assert.EqualError(t, err, "timeout")
}
