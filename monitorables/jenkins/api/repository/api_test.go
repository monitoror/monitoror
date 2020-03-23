package repository

import (
	"errors"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/config"
	pkgJenkins "github.com/monitoror/monitoror/pkg/gojenkins"
	"github.com/monitoror/monitoror/pkg/gojenkins/mocks"
	"github.com/monitoror/monitoror/pkg/gravatar"

	gojenkins "github.com/jsdidierlaurent/golang-jenkins"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, buildsAPI pkgJenkins.Jenkins) *jenkinsRepository {
	conf := &config.Jenkins{
		URL:             "http://jenkins.example.com",
		Login:           "test",
		Token:           "Test",
		Timeout:         config.Default.Timeout,
		SSLVerify:       config.Default.SSLVerify,
		InitialMaxDelay: config.Default.InitialMaxDelay,
	}

	repository := NewJenkinsRepository(conf)

	apiJenkinsRepository, ok := repository.(*jenkinsRepository)
	if assert.True(t, ok) {
		apiJenkinsRepository.jenkinsAPI = buildsAPI
		return apiJenkinsRepository
	}
	return nil
}

func TestRepository_GetJob_Error(t *testing.T) {
	jenkinsErr := errors.New("jenkins error")

	mocksJenkins := new(mocks.Jenkins)
	mocksJenkins.On("GetJob", AnythingOfType("string")).
		Return(gojenkins.Job{}, jenkinsErr)

	repository := initRepository(t, mocksJenkins)
	if repository != nil {
		_, err := repository.GetJob("master", "test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "jenkins error")
		mocksJenkins.AssertNumberOfCalls(t, "GetJob", 1)
		mocksJenkins.AssertExpectations(t)
	}
}

func TestRepository_GetJob_Success(t *testing.T) {
	jenkinsJob := gojenkins.Job{
		Buildable: true,
		InQueue:   false,
		QueueItem: gojenkins.QueueItem{
			InQueueSince: 123456789,
		},
	}

	mocksJenkins := new(mocks.Jenkins)
	mocksJenkins.On("GetJob", AnythingOfType("string")).
		Return(jenkinsJob, nil)

	// Expected
	expectedJob := &models.Job{
		ID:        "test/job/master",
		Buildable: true,
		InQueue:   false,
		QueuedAt:  nil,
		Branches:  []string{},
	}

	repository := initRepository(t, mocksJenkins)
	if repository != nil {
		job, err := repository.GetJob("test", "master")
		assert.NoError(t, err)
		assert.Equal(t, expectedJob, job)
		mocksJenkins.AssertNumberOfCalls(t, "GetJob", 1)
		mocksJenkins.AssertExpectations(t)
	}
}

func TestRepository_GetJob_SuccessWithQueue(t *testing.T) {
	jenkinsJob := gojenkins.Job{
		Buildable: true,
		InQueue:   true,
		QueueItem: gojenkins.QueueItem{
			InQueueSince: 123456789,
		},
	}

	mocksJenkins := new(mocks.Jenkins)
	mocksJenkins.On("GetJob", AnythingOfType("string")).
		Return(jenkinsJob, nil)

	// Expected
	date := parseDate(123456789)
	expectedJob := &models.Job{
		ID:        "test/job/master",
		Buildable: true,
		InQueue:   true,
		QueuedAt:  &date,
		Branches:  []string{},
	}

	repository := initRepository(t, mocksJenkins)
	if repository != nil {
		job, err := repository.GetJob("test", "master")
		assert.NoError(t, err)
		assert.Equal(t, expectedJob, job)
		mocksJenkins.AssertNumberOfCalls(t, "GetJob", 1)
		mocksJenkins.AssertExpectations(t)
	}
}

func TestRepository_GetJob_SuccessWithBranch(t *testing.T) {
	jenkinsJob := gojenkins.Job{
		Jobs: []gojenkins.SubJobDescription{{
			Name:  "master",
			Url:   "http://jenkins.example.com/job/test/job/master",
			Color: "blue",
		}},
	}

	mocksJenkins := new(mocks.Jenkins)
	mocksJenkins.On("GetJob", AnythingOfType("string")).
		Return(jenkinsJob, nil)

	// Expected
	expectedJob := &models.Job{
		ID:        "test/job/master",
		Buildable: false,
		InQueue:   false,
		Branches:  []string{"master"},
	}

	repository := initRepository(t, mocksJenkins)
	if repository != nil {
		job, err := repository.GetJob("test", "master")
		assert.NoError(t, err)
		assert.Equal(t, expectedJob, job)
		mocksJenkins.AssertNumberOfCalls(t, "GetJob", 1)
		mocksJenkins.AssertExpectations(t)
	}
}

func TestRepository_GetLastBuildStatus_Error(t *testing.T) {
	jenkinsErr := errors.New("jenkins error")

	mocksJenkins := new(mocks.Jenkins)
	mocksJenkins.On("GetLastBuildByJobId", AnythingOfType("string")).
		Return(gojenkins.Build{}, jenkinsErr)

	repository := initRepository(t, mocksJenkins)
	if repository != nil {
		build, err := repository.GetLastBuildStatus(&models.Job{ID: "test/job/master"})
		assert.Error(t, err)
		assert.Nil(t, build)
		mocksJenkins.AssertNumberOfCalls(t, "GetLastBuildByJobId", 1)
		mocksJenkins.AssertExpectations(t)
	}
}

func TestRepository_GetLastBuildStatus_Success(t *testing.T) {
	// Params
	jenkinsBuild := gojenkins.Build{
		Number:          1,
		FullDisplayName: "test/master",
		Building:        true,
		Result:          "SUCCESS",
		Timestamp:       123456789,
		Duration:        123,
		ChangeSets: []gojenkins.ScmChangeSet{{
			Items: []gojenkins.ChangeSetItem{{
				AuthorEmail: "test@test.test",
				Author: gojenkins.ScmAuthor{
					FullName: "test",
				},
			}},
		}},
	}

	mockJenkins := new(mocks.Jenkins)
	mockJenkins.On("GetLastBuildByJobId", AnythingOfType("string")).
		Return(jenkinsBuild, nil)

	// Expected
	expectedBuild := &models.Build{
		Number:   "1",
		FullName: jenkinsBuild.FullDisplayName,
		Author: &coreModels.Author{
			Name:      jenkinsBuild.ChangeSets[0].Items[0].Author.FullName,
			AvatarURL: gravatar.GetGravatarURL(jenkinsBuild.ChangeSets[0].Items[0].AuthorEmail),
		},

		Building:  jenkinsBuild.Building,
		Result:    jenkinsBuild.Result,
		StartedAt: parseDate(jenkinsBuild.Timestamp),
		Duration:  parseDuration(jenkinsBuild.Duration),
	}

	repository := initRepository(t, mockJenkins)
	if repository != nil {
		build, err := repository.GetLastBuildStatus(&models.Job{ID: "test/job/master"})
		assert.NoError(t, err)
		assert.Equal(t, expectedBuild, build)
		mockJenkins.AssertNumberOfCalls(t, "GetLastBuildByJobId", 1)
		mockJenkins.AssertExpectations(t)
	}
}
