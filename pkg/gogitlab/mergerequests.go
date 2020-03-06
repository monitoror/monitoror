package gogitlab

import "github.com/xanzy/go-gitlab"

type MergeRequestsService interface {
	ListProjectMergeRequests(
		pid interface{},
		opt *gitlab.ListProjectMergeRequestsOptions,
		options ...gitlab.OptionFunc,
	) ([]*gitlab.MergeRequest, *gitlab.Response, error)
}
