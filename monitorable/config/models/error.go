package models

import "errors"

var (
	ErrConfigFileNotFound = errors.New("config file not found")
	ErrInvalidVersionFormat = errors.New("invalid version format")
)
