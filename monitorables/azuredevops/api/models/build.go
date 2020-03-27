package models

import (
	"time"

	"github.com/monitoror/monitoror/models"
)

type (
	Build struct {
		BuildNumber    string
		DefinitionName string
		Branch         string
		Author         *models.Author

		Status string
		Result string

		FinishedAt *time.Time
		StartedAt  *time.Time
		QueuedAt   *time.Time
	}
)
