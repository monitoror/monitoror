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

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"

	xml2json "github.com/basgys/goxml2json"
	"github.com/dustin/go-humanize"
	"github.com/ghodss/yaml"
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

func (hu *httpUsecase) HTTPStatus(params *httpModels.HTTPStatusParams) (*models.Tile, error) {
	return hu.httpAll(http.HTTPStatusTileType, params.URL, params)
}

func (hu *httpUsecase) HTTPRaw(params *httpModels.HTTPRawParams) (*models.Tile, error) {
	return hu.httpAll(http.HTTPRawTileType, params.URL, params)
}

func (hu *httpUsecase) HTTPFormatted(params *httpModels.HTTPFormattedParams) (*models.Tile, error) {
	return hu.httpAll(http.HTTPFormattedTileType, params.URL, params)
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
		// Convert XML to JSON if Format == XML
		if formattedDataProvider.GetFormat() == httpModels.XMLFormat {
			buffer, err := xml2json.Convert(bytes.NewReader(response.Body))
			if err != nil || strings.TrimSuffix(buffer.String(), "\n") == `""` {
				tile.Status = models.FailedStatus
				tile.Message = fmt.Sprintf("unable to convert xml to json")
				return tile, nil
			}
			response.Body = buffer.Bytes()
		}

		// Select Unmarshaller
		var unmarshaller func(data []byte, v interface{}) error
		if formattedDataProvider.GetFormat() == httpModels.JSONFormat ||
			formattedDataProvider.GetFormat() == httpModels.XMLFormat {
			unmarshaller = json.Unmarshal
		} else {
			unmarshaller = yaml.Unmarshal
		}

		var data interface{}
		err := unmarshaller(response.Body, &data)
		if err != nil {
			tile.Status = models.FailedStatus
			tile.Message = fmt.Sprintf("unable to unmarshal content")
			return tile, nil
		}

		// Lookup a key
		if match, content = lookupKey(formattedDataProvider, data); !match {
			tile.Status = models.FailedStatus
			tile.Message = fmt.Sprintf(`unable to lookup for key "%s"`, formattedDataProvider.GetKey())
			return tile, nil
		}
	} else {
		content = string(response.Body)
	}

	// Match regex
	if regexProvider, ok := params.(httpModels.RegexProvider); ok {
		match, matchedContent := matchRegex(regexProvider, content)
		if match {
			content = matchedContent
		} else {
			tile.Status = models.FailedStatus
		}
	}

	if content != "" {
		if _, err := strconv.ParseFloat(content, 64); err == nil {
			tile.WithValue(models.NumberUnit)
		} else {
			tile.WithValue(models.RawUnit)
		}
		tile.Value.Values = []string{content}
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

	substrings := regex.FindStringSubmatch(str)
	if len(substrings) == 1 {
		return true, str
	}

	return true, substrings[1]
}

// extractValue extract value from interface{} (json/yaml/...)
// the key is in doted format like this ".bloc1."bloc.2".[2].value"
func lookupKey(params httpModels.FormatedDataProvider, data interface{}) (bool, string) {
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

	switch data := data.(type) {
	case float64:
		return true, humanize.Ftoa(data)
	default:
		return true, fmt.Sprintf("%v", data)
	}
}
