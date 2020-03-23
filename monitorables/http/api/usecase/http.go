//+build !faker

package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/models"
	"github.com/monitoror/monitoror/pkg/humanize"

	xml2json "github.com/basgys/goxml2json"
	"github.com/ghodss/yaml"
	"github.com/jsdidierlaurent/echo-middleware/cache"
)

type (
	httpUsecase struct {
		repository api.Repository

		// store used for caching request on same url
		store           cache.Store
		cacheExpiration int
	}
)

var (
	KeySplitterRegex  = regexp.MustCompile(`"[^"]*"|[^.]+`)
	ArrayKeyPartRegex = regexp.MustCompile(`^\[(\d*)]$`)
)

func NewHTTPUsecase(repository api.Repository, store cache.Store, cacheExpiration int) api.Usecase {
	return &httpUsecase{repository, store, cacheExpiration}
}

func (hu *httpUsecase) HTTPStatus(params *models.HTTPStatusParams) (*coreModels.Tile, error) {
	return hu.httpAll(api.HTTPStatusTileType, params.URL, params)
}

func (hu *httpUsecase) HTTPRaw(params *models.HTTPRawParams) (*coreModels.Tile, error) {
	return hu.httpAll(api.HTTPRawTileType, params.URL, params)
}

func (hu *httpUsecase) HTTPFormatted(params *models.HTTPFormattedParams) (*coreModels.Tile, error) {
	return hu.httpAll(api.HTTPFormattedTileType, params.URL, params)
}

// httpAll handle all http usecase by checking if params match interfaces listed in coreModels.params
func (hu *httpUsecase) httpAll(tileType coreModels.TileType, url string, params interface{}) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(tileType)
	tile.Label = url
	tile.Status = coreModels.SuccessStatus

	// Download page
	response, err := hu.get(url)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: fmt.Sprintf("unable to get %s", url)}
	}

	// Check Status Code
	if statusCodeRangeProvider, ok := params.(models.StatusCodesProvider); ok {
		if !checkStatusCode(statusCodeRangeProvider, response.StatusCode) {
			tile.Status = coreModels.FailedStatus
			tile.Message = fmt.Sprintf("status code %d", response.StatusCode)
			return tile, nil
		}
	}

	if tileType == api.HTTPStatusTileType {
		return tile, nil
	}

	// Unmarshal page
	var content string
	var match bool

	if formattedDataProvider, ok := params.(models.FormattedDataProvider); ok {
		// Convert XML to JSON if Format == XML
		if formattedDataProvider.GetFormat() == models.XMLFormat {
			buffer, err := xml2json.Convert(bytes.NewReader(response.Body))
			if err != nil || strings.TrimSuffix(buffer.String(), "\n") == `""` {
				tile.Status = coreModels.FailedStatus
				tile.Message = fmt.Sprintf("unable to convert xml to json")
				return tile, nil
			}
			response.Body = buffer.Bytes()
		}

		// Select Unmarshaller
		var unmarshaller func(data []byte, v interface{}) error
		if formattedDataProvider.GetFormat() == models.JSONFormat ||
			formattedDataProvider.GetFormat() == models.XMLFormat {
			unmarshaller = json.Unmarshal
		} else {
			unmarshaller = yaml.Unmarshal
		}

		var data interface{}
		err := unmarshaller(response.Body, &data)
		if err != nil {
			tile.Status = coreModels.FailedStatus
			tile.Message = fmt.Sprintf("unable to unmarshal content")
			return tile, nil
		}

		// Lookup a key
		if match, content = lookupKey(formattedDataProvider, data); !match {
			tile.Status = coreModels.FailedStatus
			tile.Message = fmt.Sprintf(`unable to lookup for key %q`, formattedDataProvider.GetKey())
			return tile, nil
		}
	} else {
		content = string(response.Body)
	}

	// Match regex
	if regexProvider, ok := params.(models.RegexProvider); ok {
		match, matchedContent := matchRegex(regexProvider, content)
		if match {
			content = matchedContent
		} else {
			tile.Status = coreModels.FailedStatus
		}
	}

	if content != "" {
		if _, err := strconv.ParseFloat(content, 64); err == nil {
			tile.WithValue(coreModels.NumberUnit)
		} else {
			tile.WithValue(coreModels.RawUnit)
		}
		tile.Value.Values = []string{content}
	}

	return tile, nil
}

// Adding cache to Repository.Get
func (hu *httpUsecase) get(url string) (*models.Response, error) {
	response := &models.Response{}

	// Lookup in cache
	key := fmt.Sprintf("%s:%s", coreModels.UpstreamStoreKeyPrefix, url)
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
	_ = hu.store.Set(key, *response, time.Millisecond*time.Duration(hu.cacheExpiration))

	return response, nil
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

	substrings := regex.FindStringSubmatch(str)
	if len(substrings) == 1 {
		return true, str
	}

	return true, substrings[1]
}

// extractValue extract value from interface{} (json/yaml/...)
// the key is in doted format like this ".bloc1."bloc.2".[2].value"
func lookupKey(params models.FormattedDataProvider, data interface{}) (bool, string) {
	// split key
	matchedString := KeySplitterRegex.FindAllStringSubmatch(params.GetKey(), -1)

	for _, part := range matchedString {
		keyPart := part[0]

		// Lookup for array
		r := ArrayKeyPartRegex.FindStringSubmatch(keyPart)
		if len(r) == 2 {
			arrayIndex, _ := strconv.Atoi(r[1])
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

	return true, humanize.Interface(data)
}
