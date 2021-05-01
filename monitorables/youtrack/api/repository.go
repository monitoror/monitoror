package api

import (
	"github.com/monitoror/monitoror/monitorables/youtrack/api/models"
)

type (
	Repository interface {
		GetIssues(query string) (*models.Issues, error)
	}
)
