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
		config *config.Config

		// Interfaces for Jenkins API
		jenkinsApi pkgJenkins.Jenkins
	}
)

func NewJenkinsRepository(conf *config.Config) jenkins.Repository {
	jenkinsConf := conf.Monitorable.Jenkins

	var auth *gojenkins.Auth
	if jenkinsConf.Login != "" {
		auth = &gojenkins.Auth{
			Username: jenkinsConf.Login,
			ApiToken: jenkinsConf.Token,
		}
	}
	jenkins := gojenkins.NewJenkins(auth, jenkinsConf.Url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !jenkinsConf.SSLVerify},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(jenkinsConf.Timeout) * time.Millisecond}
	jenkins.SetHTTPClient(client)

	return &jenkinsRepository{
		conf,
		jenkins,
	}
}

func (r *jenkinsRepository) GetJob(jobName string, jobParent string) (job *models.Job, err error) {
	jobId := jobName
	if jobParent != "" {
		jobId = fmt.Sprintf("%s/job/%s", jobParent, jobName)
	}

	jenkinsJob, err := r.jenkinsApi.GetJob(jobId)
	if err != nil {
		return nil, fmt.Errorf("unable to get job. %v", err)
	}

	job = &models.Job{}
	job.ID = jobId

	job.Buildable = jenkinsJob.Buildable
	job.InQueue = jenkinsJob.InQueue

	if job.InQueue {
		date := parseDate(jenkinsJob.QueueItem.InQueueSince)
		job.QueuedAt = &date
	}

	return
}

//GetBuildStatus fetch build information from travis-ci
func (r *jenkinsRepository) GetLastBuildStatus(job *models.Job) (build *models.Build, err error) {
	jenkinsBuild, err := r.jenkinsApi.GetLastBuildByJobId(job.ID)
	if err != nil {
		// No build found, return nil build but no error
		return nil, nil
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

func parseDate(date int) time.Time {
	return time.Unix(int64(date/int(time.Microsecond)), 0)
}

func parseDuration(duration int) time.Duration {
	return time.Duration(duration) * time.Millisecond
}
