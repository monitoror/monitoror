package usecase

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/monitorable/config/repository"

	"github.com/stretchr/testify/assert"
)

func TestUsecase_Hydrate(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "EMPTY" },
    { "type": "PING", "params": { "hostname": "aserver.com" } },
    { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } },
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" } },
      { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } }
    ]},
		{ "type": "JENKINS-BUILD", "params": { "job": "test" } },
		{ "type": "JENKINS-BUILD", "configVariant": "variant1", "params": { "job": "test" } }
  ]
}
`

	usecase := initConfigUsecase()
	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.GetConfig(reader)
	assert.NoError(t, err)

	err = usecase.Hydrate(config, "http://localhost:8080")
	assert.NoError(t, err)

	assert.Equal(t, "http://localhost:8080/ping?hostname=aserver.com", config.Tiles[1].Url)
	assert.Equal(t, "http://localhost:8080/port?hostname=bserver.com&port=22", config.Tiles[2].Url)

	group := config.Tiles[3].Tiles
	assert.Equal(t, "http://localhost:8080/ping?hostname=aserver.com", group[0].Url)
	assert.Equal(t, "http://localhost:8080/port?hostname=bserver.com&port=22", group[1].Url)

	assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[4].Url)
	assert.Equal(t, "http://localhost:8080/jenkins/variant1?job=test", config.Tiles[5].Url)
}
