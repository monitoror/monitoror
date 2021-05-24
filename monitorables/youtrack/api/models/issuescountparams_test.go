package models

import (
	"testing"

	coreModels "github.com/monitoror/monitoror/models"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestIssuesCountParams_GetPriorityFieldLabelWithDefault(t *testing.T) {
	param := &IssuesCountParams{}
	assert.Equal(t, DefaultPriorityFieldLabel, param.GetPriorityFieldLabelWithDefault())
	param.PriorityFieldLabel = pointer.ToString("test")
	assert.Equal(t, "test", param.GetPriorityFieldLabelWithDefault())
}

func TestIssuesCountParams_GetPriorityFieldValueWithDefault(t *testing.T) {
	param := &IssuesCountParams{}
	assert.Equal(t, DefaultPriorityFieldValue, param.GetPriorityFieldValueWithDefault())
	param.PriorityFieldValue = pointer.ToString("test")
	assert.Equal(t, "test", param.GetPriorityFieldValueWithDefault())
}

func TestIssuesCountParams_Validate_Error(t *testing.T) {
	param := &IssuesCountParams{
		CountThreshold: map[coreModels.TileStatus]int{
			coreModels.TileStatus("TEST"):    1,
			coreModels.TileStatus("WARNING"): 1,
		},
		PriorityFieldThreshold: map[coreModels.TileStatus]int{
			coreModels.TileStatus("TEST"): 1,
		},
	}
	errs := param.Validate()
	assert.Len(t, errs, 2)
}

func TestIssuesCountParams_Validate_OK(t *testing.T) {
	param := &IssuesCountParams{
		CountThreshold: map[coreModels.TileStatus]int{
			coreModels.TileStatus("WARNING"): 1,
		},
		PriorityFieldThreshold: map[coreModels.TileStatus]int{
			coreModels.TileStatus("FAILURE"): 1,
		},
	}
	errs := param.Validate()
	assert.Len(t, errs, 0)
}
