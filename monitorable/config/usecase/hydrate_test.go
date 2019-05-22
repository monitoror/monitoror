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
    { "type": "empty" },
    { "type": "ping", "params": { "hostname": "aserver.com" } },
    { "type": "port", "params": { "hostname": "bserver.com", "port": 22 } },
    { "type": "group", "label": "...", "tiles": [
      { "type": "ping", "params": { "hostname": "aserver.com" } },
      { "type": "port", "params": { "hostname": "bserver.com", "port": 22 } }
    ]}
  ]
}
`

	usecase := initConfigUsecase()
	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.GetConfig(reader)
	assert.NoError(t, err)

	err = usecase.Hydrate(config, "http://localhost:8080")
	assert.NoError(t, err)

	assert.Equal(t, "http://localhost:8080/ping?hostname=aserver.com", config.Tiles[1][UrlKey])
	assert.Equal(t, "http://localhost:8080/port?hostname=bserver.com&port=22", config.Tiles[2][UrlKey])

	group := config.Tiles[3][TilesKey].([]interface{})
	gtile1 := group[0].(map[string]interface{})
	assert.Equal(t, "http://localhost:8080/ping?hostname=aserver.com", gtile1[UrlKey])
	gtile2 := group[1].(map[string]interface{})
	assert.Equal(t, "http://localhost:8080/port?hostname=bserver.com&port=22", gtile2[UrlKey])
}
