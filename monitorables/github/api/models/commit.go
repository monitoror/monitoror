package models

import "github.com/monitoror/monitoror/models"

type (
	Commit struct {
		SHA    string
		Author *models.Author
	}
)
