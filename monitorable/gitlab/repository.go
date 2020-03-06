package gitlab

import "github.com/monitoror/monitoror/monitorable/gitlab/models"

type (
	Repository interface {
		GetCount(query string) (int, error)
		GetPipelines(repository, ref string) (*models.Pipelines, error)
		GetMergeRequests(repository string) ([]models.MergeRequest, error)
		GetCommit(repository, sha string) (*models.Commit, error)
	}
)
