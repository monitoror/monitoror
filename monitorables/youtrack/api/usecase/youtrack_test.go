package usecase

import (
	"errors"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/api"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/mocks"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/models"

	"github.com/stretchr/testify/assert"
)

const query = "project: MO #Unresolved"
const priorityQuery = "(project: MO #Unresolved) and Priority: {Show-stopper}"

func TestYoutrackUsecase_CountIssues_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("CountIssues", query).Return(12, nil)

	tile, err := NewYoutrackUsecase(mockRepository).CountIssues(&models.IssuesCountParams{Query: query})
	if assert.NoError(t, err) {
		assert.Equal(t, api.YoutrackCountIssuesTileType, tile.Type)
		assert.Equal(t, coreModels.SuccessStatus, tile.Status)
		assert.Equal(t, []string{"12"}, tile.Metrics.Values)
		mockRepository.AssertNumberOfCalls(t, "CountIssues", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestYoutrackUsecase_CountIssues_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("CountIssues", query).Return(0, errors.New("boom"))

	tile, err := NewYoutrackUsecase(mockRepository).CountIssues(&models.IssuesCountParams{Query: query})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to count issues", err.Error())
	}
}

func TestYoutrackUsecase_CountIssues_PriorityError(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("CountIssues", query).Return(12, nil)
	mockRepository.On("CountIssues", priorityQuery).Return(0, errors.New("boom"))

	params := &models.IssuesCountParams{
		Query:                  query,
		PriorityFieldThreshold: coreModels.IntThreshold{coreModels.FailedStatus: 1},
	}

	tile, err := NewYoutrackUsecase(mockRepository).CountIssues(params)
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, "unable to count priority issues", err.Error())
	}
}

func TestYoutrackUsecase_CountIssues_Thresholds(t *testing.T) {
	for _, testcase := range []struct {
		count, priorityCount int
		expected             coreModels.TileStatus
	}{
		{count: 2, priorityCount: 0, expected: coreModels.SuccessStatus},
		{count: 5, priorityCount: 0, expected: coreModels.WarningStatus},
		{count: 10, priorityCount: 0, expected: coreModels.FailedStatus},
		{count: 2, priorityCount: 1, expected: coreModels.WarningStatus},
		{count: 2, priorityCount: 3, expected: coreModels.FailedStatus},
		// A less severe priority status must not hide the count status
		{count: 10, priorityCount: 1, expected: coreModels.FailedStatus},
	} {
		mockRepository := new(mocks.Repository)
		mockRepository.On("CountIssues", query).Return(testcase.count, nil)
		mockRepository.On("CountIssues", priorityQuery).Return(testcase.priorityCount, nil)

		params := &models.IssuesCountParams{
			Query:                  query,
			CountThreshold:         coreModels.IntThreshold{coreModels.WarningStatus: 5, coreModels.FailedStatus: 10},
			PriorityFieldThreshold: coreModels.IntThreshold{coreModels.WarningStatus: 1, coreModels.FailedStatus: 3},
		}

		tile, err := NewYoutrackUsecase(mockRepository).CountIssues(params)
		if assert.NoError(t, err) {
			assert.Equal(t, testcase.expected, tile.Status, "count=%d priorityCount=%d", testcase.count, testcase.priorityCount)
			mockRepository.AssertNumberOfCalls(t, "CountIssues", 2)
		}
	}
}
