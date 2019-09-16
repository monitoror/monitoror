package http

import (
	"github.com/monitoror/monitoror/monitorable/http/models"
)

// Repository represent the http's repository contract
type (
	Repository interface {
		Get(url string) (*models.Response, error)
	}
)
