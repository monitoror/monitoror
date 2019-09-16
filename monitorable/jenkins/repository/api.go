package repository

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	"github.com/monitoror/monitoror/monitorable/jenkins/models"
	pkgJenkins "github.com/monitoror/monitoror/pkg/gojenkins"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/gravatar"

	gojenkins "github.com/jsdidierlaurent/golang-jenkins"
)

type (
	jenkinsRepository struct {
		// Interfaces for Jenkins API
		jenkinsApi pkgJenkins.Jenkins
	}
)

func NewJenkinsRepository(config *config.Jenkins) jenkins.Repository {
	var auth *gojenkins.Auth
	if config.Login != "" {
		auth = &gojenkins.Auth{
			Username: config.Login,
			ApiToken: config.Token,
		}
	}
	jenkins := gojenkins.NewJenkins(auth, config.Url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.SSLVerify},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(config.Timeout) * time.Millisecond}
	jenkins.SetHTTPClient(client)

	return &jenkinsRepository{
		jenkins,
	}
}

func (r *jenkinsRepository) GetJob(jobName string, branch string) (job *models.Job, err error) {
	jobId := jobName
	if branch != "" {
		jobId = fmt.Sprintf("%s/job/%s", jobName, branch)
	}

	jenkinsJob, err := r.jenkinsApi.GetJob(jobId)
	if err != nil {
		return nil, err
	}

	job = &models.Job{}
	job.ID = jobId

	job.Buildable = jenkinsJob.Buildable
	job.InQueue = jenkinsJob.InQueue

	if job.InQueue {
		date := parseDate(jenkinsJob.QueueItem.InQueueSince)
		job.QueuedAt = &date
	}

	job.Branches = []string{}
	for _, jenkinsSubJob := range jenkinsJob.Jobs {
		if jenkinsSubJob.Color != "disabled" { // Filter old merged branch
			job.Branches = append(job.Branches, jenkinsSubJob.Name)
		}
	}

	return
}

//GetBuildStatus fetch build information from travis-ci
func (r *jenkinsRepository) GetLastBuildStatus(job *models.Job) (build *models.Build, err error) {
	jenkinsBuild, err := r.jenkinsApi.GetLastBuildByJobId(job.ID)
	if err != nil {
		return nil, err
	}

	build = &models.Build{}
	build.Number = string(jenkinsBuild.Number)
	build.FullName = jenkinsBuild.FullDisplayName

	build.Building = jenkinsBuild.Building

	build.Result = jenkinsBuild.Result
	build.StartedAt = parseDate(jenkinsBuild.Timestamp)
	build.Duration = parseDuration(jenkinsBuild.Duration)

	// ChangeSet or ChangeSets in case of pipeline
	changeSet := jenkinsBuild.ChangeSet
	if len(jenkinsBuild.ChangeSets) > 0 {
		changeSet = jenkinsBuild.ChangeSets[0]
	}

	if len(changeSet.Items) > 0 {
		item := changeSet.Items[0]
		build.Author = &models.Author{
			Name: item.Author.FullName,
		}

		if item.AuthorEmail != "" {
			build.Author.AvatarUrl = gravatar.GetGravatarUrl(item.AuthorEmail)
		}
	}

	return
}

func parseDate(date int64) time.Time {
	return time.Unix(date/int64(time.Microsecond), 0)
}

func parseDuration(duration int64) time.Duration {
	return time.Duration(duration) * time.Millisecond
}
