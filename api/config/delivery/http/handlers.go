package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/models"
)

type ConfigDelivery struct {
	configUsecase config.Usecase
}

func NewConfigDelivery(cu config.Usecase) *ConfigDelivery {
	return &ConfigDelivery{cu}
}

func (h *ConfigDelivery) GetConfigList(c echo.Context) error {
	configList := h.configUsecase.GetConfigList()

	return c.JSON(http.StatusOK, configList)
}

func (h *ConfigDelivery) GetConfig(c echo.Context) error {
	// Bind / check Params
	params := &models.ConfigParams{}
	_ = c.Bind(params) // can't throw any error with this Params
	// Decode params
	params.Config, _ = url.QueryUnescape(params.Config)

	configBag := h.configUsecase.GetConfig(params)

	if len(configBag.Errors) == 0 {
		h.configUsecase.Verify(configBag)
	}
	if len(configBag.Errors) == 0 {
		h.configUsecase.Hydrate(configBag)
	}

	// By default, Marshall function escape <, > and & according https://golang.org/src/encoding/json/encode.go?s=6456:6499#L48
	// In Chromium on arm the UI code do not parse escaping character correctly
	encoded, _ := JSONMarshal(configBag) // Ignoring error, assuming there is no function or channel inside this struct

	return c.Blob(http.StatusOK, echo.MIMEApplicationJSONCharsetUTF8, encoded)
}

// JSONMarshal same as JSON.Marshall but with SetEscapeHTML(false)
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
