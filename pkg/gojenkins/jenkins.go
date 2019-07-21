package gojenkins

import (
	gojenkins "github.com/jsdidierlaurent/golang-jenkins"
)

type Jenkins interface {
	GetJob(jobName string) (job gojenkins.Job, err error)
	GetBuildByJobId(jobId string, number int) (build gojenkins.Build, err error)
	GetLastBuildByJobId(jobId string) (build gojenkins.Build, err error)
}
