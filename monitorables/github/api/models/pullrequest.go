package models

import (
	coreModels "github.com/monitoror/monitoror/models"
)

type PullRequest struct {
	ID     int
	Title  string
	Author coreModels.Author

	Owner      string
	Repository string
	Branch     string
	CommitSHA  string
}
