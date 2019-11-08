//+build !faker

package usecase

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"

	"github.com/dustin/go-humanize"
	"github.com/jsdidierlaurent/echo-middleware/cache"
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

func NewHTTPUsecase(repository http.Repository, store cache.Store, upstreamCacheExpiration int) http.Usecase {
	return &httpUsecase{repository, store, upstreamCacheExpiration}
}

func (hu *httpUsecase) HTTPAny(params *httpModels.HTTPAnyParams) (*models.Tile, error) {
	return hu.httpAll(http.HTTPAnyTileType, params.URL, params)
}

func (hu *httpUsecase) HTTPRaw(params *httpModels.HTTPRawParams) (*models.Tile, error) {
	return hu.httpAll(http.HTTPRawTileType, params.URL, params)
}

func (hu *httpUsecase) HTTPJson(params *httpModels.HTTPJsonParams) (*models.Tile, error) {
	return hu.httpAll(http.HTTPJsonTileType, params.URL, params)
}

func (hu *httpUsecase) HTTPYaml(params *httpModels.HTTPYamlParams) (*models.Tile, error) {
	return hu.httpAll(http.HTTPYamlTileType, params.URL, params)
}

// httpAll handle all http usecase by checking if params match interfaces listed in models.params
func (hu *httpUsecase) httpAll(tileType models.TileType, url string, params interface{}) (*models.Tile, error) {
	tile := models.NewTile(tileType)
	tile.Label = url
	tile.Status = models.SuccessStatus

	// Download page
	response, err := hu.get(url)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: fmt.Sprintf("unable to get %s", url)}
	}

	// Check Status Code
	if statusCodeRangeProvider, ok := params.(httpModels.StatusCodesProvider); ok {
		if !checkStatusCode(statusCodeRangeProvider, response.StatusCode) {
			tile.Status = models.FailedStatus
			tile.Message = fmt.Sprintf("status code %d", response.StatusCode)
			return tile, nil
		}
	}

	// Unmarshal page
	var content string
	var match bool

	if formattedDataProvider, ok := params.(httpModels.FormatedDataProvider); ok {
		var jsonData interface{}
		err := formattedDataProvider.GetUnmarshaller()(response.Body, &jsonData)
		if err != nil {
			tile.Status = models.FailedStatus
			tile.Message = fmt.Sprintf("unable to unmarshal content")
			return tile, nil
		}

		// Lookup a key
		if match, content = lookupKey(formattedDataProvider, jsonData); !match {
			tile.Status = models.FailedStatus
			tile.Message = fmt.Sprintf(`unable to lookup for key "%s"`, formattedDataProvider.GetKey())
			return tile, nil
		}
	} else {
		content = string(response.Body)
	}

	// Match regex
	if regexProvider, ok := params.(httpModels.RegexProvider); ok {
		if match, content = matchRegex(regexProvider, content); !match {
			tile.Status = models.FailedStatus
			tile.Message = fmt.Sprintf(`pattern not found "%s"`, regexProvider.GetRegex())
			return tile, nil
		}
	}

	if s, err := strconv.ParseFloat(content, 64); err == nil {
		tile.Values = []float64{s}
	} else {
		tile.Message = content
	}

	return tile, nil
}

// Adding cache to Repository.Get
func (hu *httpUsecase) get(url string) (*httpModels.Response, error) {
	response := &httpModels.Response{}

	// Lookup in cache
	key := fmt.Sprintf("%s:%s", models.UpstreamStoreKeyPrefix, url)
	if err := hu.store.Get(key, response); err == nil {
		// Cache found, return
		return response, nil
	}

	// Download page
	response, err := hu.repository.Get(url)
	if err != nil {
		return nil, err
	}

	// Adding result in store
	_ = hu.store.Set(key, *response, time.Millisecond*time.Duration(hu.upstreamCacheExpiration))

	return response, nil
}

// checkStatusCode check if status code is between min / max
// if min/max are empty, use default value
func checkStatusCode(params httpModels.StatusCodesProvider, code int) bool {
	min, max := params.GetStatusCodes()
	return min <= code && code <= max
}

// matchRegex check if string match regex, if the regex match, try to extract first group
func matchRegex(params httpModels.RegexProvider, str string) (bool, string) {
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
func lookupKey(params httpModels.FormatedDataProvider, data interface{}) (bool, string) {
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

	switch data := data.(type) {
	case float64:
		return true, humanize.Ftoa(data)
	default:
		return true, fmt.Sprintf("%v", data)
	}
}
