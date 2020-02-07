package github

import "github.com/monitoror/monitoror/monitorable/github/models"

type (
	Repository interface {
		GetIssuesCount(query string) (int, error)
		GetChecks(owner, repository, ref string) (*models.Checks, error)
		GetPullRequests(owner, repository string) ([]models.PullRequest, error)
		GetCommit(owner, repository, sha string) (*models.Commit, error)
	}
)
