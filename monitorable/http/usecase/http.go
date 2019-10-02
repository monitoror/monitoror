//+build !faker

package usecase

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jsdidierlaurent/echo-middleware/cache"
	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	"github.com/monitoror/monitoror/monitorable/http/models"
)

type (
	httpUsecase struct {
		repository http.Repository

		// Store used for caching request on same url
		store                   cache.Store
		upstreamCacheExpiration int
	}
)

var (
	KeySplitterRegex  = regexp.MustCompile(`"[^"]*"|[^.]+`)
	ArrayKeyPartRegex = regexp.MustCompile(`^\[(\d*)]$`)
)

const (
	HttpRequestStoreKeyPrefix = "monitoror.http.request"
)

func NewHttpUsecase(repository http.Repository, store cache.Store, upstreamCacheExpiration int) http.Usecase {
	return &httpUsecase{repository, store, upstreamCacheExpiration}
}

func (hu *httpUsecase) HttpAny(params *models.HttpAnyParams) (tile *Tile, err error) {
	return hu.httpAll(http.HttpAnyTileType, params.Url, params)
}

func (hu *httpUsecase) HttpRaw(params *models.HttpRawParams) (tile *Tile, err error) {
	return hu.httpAll(http.HttpRawTileType, params.Url, params)
}

func (hu *httpUsecase) HttpJson(params *models.HttpJsonParams) (tile *Tile, err error) {
	return hu.httpAll(http.HttpJsonTileType, params.Url, params)
}

func (hu *httpUsecase) HttpYaml(params *models.HttpYamlParams) (tile *Tile, err error) {
	return hu.httpAll(http.HttpYamlTileType, params.Url, params)
}

// httpAll handle all http usecase by checking if params match interfaces listed in models.params
func (hu *httpUsecase) httpAll(tileType TileType, url string, params interface{}) (tile *Tile, err error) {
	tile = NewTile(tileType)
	tile.Label = url
	tile.Status = SuccessStatus

	// Download page
	response, err := hu.get(url)
	if err != nil {
		return nil, &MonitororError{Err: err, Tile: tile, Message: fmt.Sprintf("unable to get %s", url)}
	}

	// Check Status Code
	if statusCodeRangeProvider, ok := params.(models.StatusCodesProvider); ok {
		if !checkStatusCode(statusCodeRangeProvider, response.StatusCode) {
			tile.Status = FailedStatus
			tile.Message = fmt.Sprintf("status code %d", response.StatusCode)
			return
		}
	}

	// Unmarshal page
	var content string
	var match bool

	if formatedDataProvider, ok := params.(models.FormatedDataProvider); ok {
		var jsonData interface{}
		err := formatedDataProvider.GetUnmarshaller()(response.Body, &jsonData)
		if err != nil {
			tile.Status = FailedStatus
			tile.Message = fmt.Sprintf("unable to unmarshal content")
			return tile, nil
		}

		// Lookup a key
		if match, content = lookupKey(formatedDataProvider, jsonData); !match {
			tile.Status = FailedStatus
			tile.Message = fmt.Sprintf(`unable to lookup for key "%s"`, formatedDataProvider.GetKey())
			return tile, nil
		}
	} else {
		content = string(response.Body)
	}

	// Match regex
	if regexProvider, ok := params.(models.RegexProvider); ok {
		if match, content = matchRegex(regexProvider, content); !match {
			tile.Status = FailedStatus
			tile.Message = fmt.Sprintf(`pattern not found "%s"`, regexProvider.GetRegex())
			return tile, nil
		}
	}

	if s, err := strconv.ParseFloat(content, 64); err == nil {
		tile.Values = []float64{s}
	} else {
		tile.Message = content
	}

	return
}

// Adding cache to Repository.Get
func (hu *httpUsecase) get(url string) (response *models.Response, err error) {
	response = &models.Response{}

	// Lookup in cache
	key := fmt.Sprintf("%s:%s", UpstreamStoreKeyPrefix, url)
	if err = hu.store.Get(key, response); err == nil {
		// Cache found, return
		return
	}

	// Download page
	response, err = hu.repository.Get(url)
	if err != nil {
		return
	}

	// Adding result in store
	err = hu.store.Set(key, *response, time.Millisecond*time.Duration(hu.upstreamCacheExpiration))

	return
}

// checkStatusCode check if status code is between min / max
// if min/max are empty, use default value
func checkStatusCode(params models.StatusCodesProvider, code int) bool {
	min, max := params.GetStatusCodes()
	return min <= code && code <= max
}

// matchRegex check if string match regex, if the regex match, try to extract first group
func matchRegex(params models.RegexProvider, str string) (bool, string) {
	regex := params.GetRegexp()
	if regex == nil {
		return true, str
	}

	if !regex.MatchString(str) {
		return false, ""
	}

	substrings := regex.FindAllStringSubmatch(str, -1)
	if len(substrings[0]) < 2 {
		return true, str
	}

	return true, substrings[0][1]
}

// extractValue extract value from interface{} (json/yaml/...)
// the key is in doted format like this ".bloc1."bloc.2".[2].value"
func lookupKey(params models.FormatedDataProvider, data interface{}) (bool, string) {
	// split key
	matchedString := KeySplitterRegex.FindAllStringSubmatch(params.GetKey(), -1)

	for _, part := range matchedString {
		keyPart := part[0]

		// Lookup for array
		r := ArrayKeyPartRegex.FindAllStringSubmatch(keyPart, 1)
		if len(r) == 1 {
			arrayIndex, _ := strconv.Atoi(r[0][1])
			// Look if data type is array and check if index wasn't out of bounds
			if array, ok := data.([]interface{}); ok && len(array) > arrayIndex && arrayIndex >= 0 {
				data = array[arrayIndex]
				continue
			}
			// If array didn't match, test with map
		}

		// Lookup for map
		keyPart = strings.ReplaceAll(keyPart, `"`, ``)
		// map[string]interface{} => JSON Style
		if m, ok := data.(map[string]interface{}); ok {
			// Check if keyPart is in map
			if value, ok := m[keyPart]; ok {
				data = value
				continue
			}
		}
		// map[interface{}]interface{} => YAML Style
		if m, ok := data.(map[interface{}]interface{}); ok {
			// Check if keyPart is in map
			if value, ok := m[keyPart]; ok {
				data = value
				continue
			}
		}

		return false, ""
	}

	switch data.(type) {
	case float64:
		return true, fmt.Sprintf("%s", humanize.Ftoa(data.(float64)))
	default:
		return true, fmt.Sprintf("%v", data)
	}
}
