package models

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ConfigError struct {
	reasons []string
}

func NewConfigError() *ConfigError {
	return &ConfigError{}
}

func (ce *ConfigError) Add(reasons ...string) {
	ce.reasons = append(ce.reasons, reasons...)
}

func (ce *ConfigError) Count() int {
	return len(ce.reasons)
}

func (ce *ConfigError) Send(ctx echo.Context) {
	//TODO
	_ = ctx.NoContent(http.StatusBadRequest)
}

func (ce *ConfigError) Error() string {
	str := "invalid configuration: \n"
	for _, value := range ce.reasons {
		str += value + "\n"
	}

	return str
}
