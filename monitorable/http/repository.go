package http

import (
	"github.com/monitoror/monitoror/monitorable/http/models"
)

type (
	Repository interface {
		Get(url string) (*models.Response, error)
	}
)
