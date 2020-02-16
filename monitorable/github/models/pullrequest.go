package models

type PullRequest struct {
	ID         int
	Owner      string
	Repository string
	Ref        string
}
