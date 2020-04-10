package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/monitoror/monitoror/monitorables/github/config"
	"github.com/monitoror/monitoror/pkg/gogithub/mocks"
	"github.com/monitoror/monitoror/pkg/gravatar"

	. "github.com/AlekSi/pointer"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T) *githubRepository {
	conf := &config.Github{
		URL:                  "https://github.example.com",
		Token:                "xxx",
		Timeout:              1000,
		CountCacheExpiration: 10000,
	}

	repository := NewGithubRepository(conf)

	assert.Equal(t, "https://github.example.com/", conf.URL)

	apiGithubRepository, ok := repository.(*githubRepository)
	if assert.True(t, ok) {
		return apiGithubRepository
	}
	return nil
}

func TestRepository_GetSearchCount_Error(t *testing.T) {
	githubErr := errors.New("github error")

	mocksSearchService := new(mocks.SearchService)
	mocksSearchService.On("Issues", Anything, AnythingOfType("string"), Anything).
		Return(nil, nil, githubErr)

	repository := initRepository(t)
	if repository != nil {
		repository.searchService = mocksSearchService

		_, err := repository.GetCount("test")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "github error")
			mocksSearchService.AssertNumberOfCalls(t, "Issues", 1)
			mocksSearchService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetSearchCount_Success(t *testing.T) {
	mocksSearchService := new(mocks.SearchService)
	mocksSearchService.On("Issues", Anything, AnythingOfType("string"), Anything).
		Return(&github.IssuesSearchResult{Total: ToInt(42)}, nil, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.searchService = mocksSearchService

		value, err := repository.GetCount("test")
		if assert.NoError(t, err) {
			assert.Equal(t, 42, value)
			mocksSearchService.AssertNumberOfCalls(t, "Issues", 1)
			mocksSearchService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetChecks_CheckServiceError(t *testing.T) {
	githubErr := errors.New("github error")

	mocksChecksService := new(mocks.ChecksService)
	mocksChecksService.
		On("ListCheckRunsForRef", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string"), Anything).
		Return(nil, nil, githubErr)

	repository := initRepository(t)
	if repository != nil {
		repository.checksService = mocksChecksService

		_, err := repository.GetChecks("test", "test", "master")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "github error")
			mocksChecksService.AssertNumberOfCalls(t, "ListCheckRunsForRef", 1)
			mocksChecksService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetChecks_RepositoriesServiceError(t *testing.T) {
	githubErr := errors.New("github error")

	mocksChecksService := new(mocks.ChecksService)
	mocksChecksService.
		On("ListCheckRunsForRef", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string"), Anything).
		Return(&github.ListCheckRunsResults{}, nil, nil)

	mocksRepositoriesService := new(mocks.RepositoriesService)
	mocksRepositoriesService.
		On("ListStatuses", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string"), Anything).
		Return(nil, nil, githubErr)

	repository := initRepository(t)
	if repository != nil {
		repository.checksService = mocksChecksService
		repository.repositoriesService = mocksRepositoriesService

		_, err := repository.GetChecks("test", "test", "master")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "github error")
			mocksChecksService.AssertNumberOfCalls(t, "ListCheckRunsForRef", 1)
			mocksChecksService.AssertExpectations(t)
			mocksRepositoriesService.AssertNumberOfCalls(t, "ListStatuses", 1)
			mocksRepositoriesService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetChecks_Success(t *testing.T) {
	checkRunsResults := &github.ListCheckRunsResults{
		CheckRuns: []*github.CheckRun{
			{
				Name:        ToString("build 1"),
				Status:      ToString("completed"),
				Conclusion:  ToString("success"),
				StartedAt:   &github.Timestamp{Time: time.Now()},
				CompletedAt: &github.Timestamp{Time: time.Now()},
				HeadSHA:     ToString("sha"),
			},
		},
	}
	statuses := []*github.RepoStatus{
		{
			Context:   ToString("app 1"),
			State:     ToString("success"),
			CreatedAt: ToTime(time.Now()),
			UpdatedAt: ToTime(time.Now()),
			URL:       ToString("/sha"),
		},
	}

	mocksChecksService := new(mocks.ChecksService)
	mocksChecksService.
		On("ListCheckRunsForRef", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string"), Anything).
		Return(checkRunsResults, nil, nil)

	mocksRepositoriesService := new(mocks.RepositoriesService)
	mocksRepositoriesService.
		On("ListStatuses", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string"), Anything).
		Return(statuses, nil, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.checksService = mocksChecksService
		repository.repositoriesService = mocksRepositoriesService

		checks, err := repository.GetChecks("test", "test", "test")
		if assert.NoError(t, err) {
			assert.Len(t, checks.Runs, 1)
			assert.Len(t, checks.Statuses, 1)

			assert.Equal(t, *checkRunsResults.CheckRuns[0].Name, checks.Runs[0].Title)
			assert.Equal(t, *checkRunsResults.CheckRuns[0].Status, checks.Runs[0].Status)
			assert.Equal(t, *checkRunsResults.CheckRuns[0].Conclusion, checks.Runs[0].Conclusion)
			assert.Equal(t, &checkRunsResults.CheckRuns[0].StartedAt.Time, checks.Runs[0].StartedAt)
			assert.Equal(t, &checkRunsResults.CheckRuns[0].CompletedAt.Time, checks.Runs[0].CompletedAt)

			assert.Equal(t, *statuses[0].Context, checks.Statuses[0].Title)
			assert.Equal(t, *statuses[0].State, checks.Statuses[0].State)
			assert.Equal(t, *statuses[0].CreatedAt, checks.Statuses[0].CreatedAt)
			assert.Equal(t, *statuses[0].UpdatedAt, checks.Statuses[0].UpdatedAt)

			assert.Equal(t, *checkRunsResults.CheckRuns[0].HeadSHA, *checks.HeadCommit)

			mocksChecksService.AssertNumberOfCalls(t, "ListCheckRunsForRef", 1)
			mocksChecksService.AssertExpectations(t)
			mocksRepositoriesService.AssertNumberOfCalls(t, "ListStatuses", 1)
			mocksRepositoriesService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetPullRequest_Error(t *testing.T) {
	githubErr := errors.New("github error")

	mocksPullRequestService := new(mocks.PullRequestService)
	mocksPullRequestService.On("List", Anything, AnythingOfType("string"), AnythingOfType("string"), Anything).
		Return(nil, nil, githubErr)

	repository := initRepository(t)
	if repository != nil {
		repository.pullRequestService = mocksPullRequestService

		_, err := repository.GetPullRequests("test", "test")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "github error")
			mocksPullRequestService.AssertNumberOfCalls(t, "List", 1)
			mocksPullRequestService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetPullRequest_Success(t *testing.T) {
	mocksPullRequestService := new(mocks.PullRequestService)
	mocksPullRequestService.On("List", Anything, AnythingOfType("string"), AnythingOfType("string"), Anything).
		Return([]*github.PullRequest{
			{
				Number: ToInt(10),
				Title:  ToString("Test"),
				Head:   &github.PullRequestBranch{Ref: ToString("master")},
			},
		}, nil, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.pullRequestService = mocksPullRequestService

		pullRequests, err := repository.GetPullRequests("test", "test")
		if assert.NoError(t, err) {
			assert.Len(t, pullRequests, 1)
			assert.Equal(t, 10, pullRequests[0].ID)
			assert.Equal(t, "test", pullRequests[0].Owner)
			assert.Equal(t, "test", pullRequests[0].Repository)
			assert.Equal(t, "master", pullRequests[0].Ref)

			mocksPullRequestService.AssertNumberOfCalls(t, "List", 1)
			mocksPullRequestService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetCommit_Error(t *testing.T) {
	githubErr := errors.New("github error")

	mocksGitService := new(mocks.GitService)
	mocksGitService.On("GetCommit", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, nil, githubErr)

	repository := initRepository(t)
	if repository != nil {
		repository.gitService = mocksGitService

		_, err := repository.GetCommit("test", "test", "sha")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "github error")
		mocksGitService.AssertNumberOfCalls(t, "GetCommit", 1)
		mocksGitService.AssertExpectations(t)
	}
}

func TestRepository_GetCommit_Success(t *testing.T) {
	mocksGitService := new(mocks.GitService)
	mocksGitService.On("GetCommit", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&github.Commit{
			Author: &github.CommitAuthor{
				Name:  ToString("Test"),
				Email: ToString("test@example.com"),
			},
		}, nil, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.gitService = mocksGitService

		commit, err := repository.GetCommit("test", "test", "sha")
		if assert.NoError(t, err) {
			assert.Equal(t, "sha", commit.SHA)
			assert.Equal(t, "Test", commit.Author.Name)
			assert.Equal(t, gravatar.GetGravatarURL("test@example.com"), commit.Author.AvatarURL)
			mocksGitService.AssertNumberOfCalls(t, "GetCommit", 1)
			mocksGitService.AssertExpectations(t)
		}
	}
}
