package repository

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/api"
	"github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/config"
	pkgJenkins "github.com/monitoror/monitoror/pkg/gojenkins"
	"github.com/monitoror/monitoror/pkg/gravatar"

	gojenkins "github.com/jsdidierlaurent/golang-jenkins"
)

type (
	jenkinsRepository struct {
		// Interfaces for Jenkins API
		jenkinsAPI pkgJenkins.Jenkins
	}
)

func NewJenkinsRepository(config *config.Jenkins) api.Repository {
	var auth *gojenkins.Auth
	if config.Login != "" {
		auth = &gojenkins.Auth{
			Username: config.Login,
			ApiToken: config.Token,
		}
	}

	// Remove last /
	if strings.HasSuffix(config.URL, "/") {
		config.URL = strings.TrimRight(config.URL, "/")
	}
	jenkins := gojenkins.NewJenkins(auth, config.URL)

	// Override transport
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
	jobID := jobName
	if branch != "" {
		jobID = fmt.Sprintf("%s/job/%s", jobName, branch)
	}

	jenkinsJob, err := r.jenkinsAPI.GetJob(jobID)
	if err != nil {
		return nil, err
	}

	job = &models.Job{}
	job.ID = jobID

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
func (r *jenkinsRepository) GetLastBuildStatus(job *models.Job) (*models.Build, error) {
	jenkinsBuild, err := r.jenkinsAPI.GetLastBuildByJobId(job.ID)
	if err != nil {
		return nil, err
	}

	build := &models.Build{}
	build.Number = fmt.Sprintf("%d", jenkinsBuild.Number)
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
		build.Author = &coreModels.Author{
			Name: item.Author.FullName,
		}

		if item.AuthorEmail != "" {
			build.Author.AvatarURL = gravatar.GetGravatarURL(item.AuthorEmail)
		}
	}

	return build, nil
}

func parseDate(date int64) time.Time {
	return time.Unix(date/int64(time.Microsecond), 0)
}

func parseDuration(duration int64) time.Duration {
	return time.Duration(duration) * time.Millisecond
}
