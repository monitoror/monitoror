package gogitlab

import "github.com/xanzy/go-gitlab"

type SearchService interface {
	Issues(
		query string,
		opt *gitlab.SearchOptions,
		options ...gitlab.OptionFunc,
	) ([]*gitlab.Issue, *gitlab.Response, error)
}
