package models

type MergeRequest struct {
	ID         int
	Repository string
	Ref        string
}
