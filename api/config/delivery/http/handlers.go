package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/monitoror/monitoror/api/config"
	configModels "github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/models"

	"github.com/labstack/echo/v4"
)

type ConfigDelivery struct {
	configUsecase config.Usecase
}

func NewConfigDelivery(cu config.Usecase) *ConfigDelivery {
	return &ConfigDelivery{cu}
}

func (h *ConfigDelivery) GetConfig(c echo.Context) error {
	// Bind / check Params
	params := &configModels.ConfigParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

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
