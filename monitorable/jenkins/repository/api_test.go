package repository

import (
	"errors"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/gravatar"

	"github.com/monitoror/monitoror/monitorable/jenkins/models"

	gojenkins "github.com/jsdidierlaurent/golang-jenkins"

	. "github.com/monitoror/monitoror/config"
	pkgJenkins "github.com/monitoror/monitoror/pkg/gojenkins"
	"github.com/monitoror/monitoror/pkg/gojenkins/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, buildsApi pkgJenkins.Jenkins) *jenkinsRepository {
	conf := InitConfig()
	conf.Monitorable.Jenkins[DefaultVariant].Url = "http://jenkins.test.com"
	conf.Monitorable.Jenkins[DefaultVariant].Login = "test"
	conf.Monitorable.Jenkins[DefaultVariant].Token = "test"

	repository := NewJenkinsRepository(conf.Monitorable.Jenkins[DefaultVariant])

	apiJenkinsRepository, ok := repository.(*jenkinsRepository)
	if assert.True(t, ok) {
		apiJenkinsRepository.jenkinsApi = buildsApi
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
		assert.True(t, strings.Contains(err.Error(), "unable to get job. jenkins error"))
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
	}

	repository := initRepository(t, mocksJenkins)
	if repository != nil {
		job, err := repository.GetJob("master", "test")
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
	}

	repository := initRepository(t, mocksJenkins)
	if repository != nil {
		job, err := repository.GetJob("master", "test")
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
		Number:   string(jenkinsBuild.Number),
		FullName: jenkinsBuild.FullDisplayName,
		Author: &models.Author{
			Name:      jenkinsBuild.ChangeSets[0].Items[0].Author.FullName,
			AvatarUrl: gravatar.GetGravatarUrl(jenkinsBuild.ChangeSets[0].Items[0].AuthorEmail),
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
