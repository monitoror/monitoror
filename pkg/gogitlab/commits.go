package gogitlab

import "github.com/xanzy/go-gitlab"

type CommitsService interface {
	GetCommit(
		pid interface{},
		sha string,
		options ...gitlab.OptionFunc,
	) (*gitlab.Commit, *gitlab.Response, error)
}
