package gogitlab

import "github.com/xanzy/go-gitlab"

type PipelinesService interface {
	ListProjectPipelines(
		pid interface{},
		opt *gitlab.ListProjectPipelinesOptions,
		options ...gitlab.OptionFunc,
	) ([]*gitlab.PipelineInfo, *gitlab.Response, error)

	GetPipeline(
		pid interface{},
		id int,
		options ...gitlab.OptionFunc,
	) (*gitlab.Pipeline, *gitlab.Response, error)
}
