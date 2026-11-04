package models

import (
	"testing"

	coreModels "github.com/monitoror/monitoror/models"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestIssuesCountParams_GetPriorityQuery(t *testing.T) {
	param := &IssuesCountParams{Query: "project: MO #Unresolved"}
	assert.Equal(t, "(project: MO #Unresolved) and Priority: {Show-stopper}", param.GetPriorityQuery())

	param.PriorityFieldLabel = pointer.ToString("Severity")
	param.PriorityFieldValue = pointer.ToString("Very High")
	assert.Equal(t, "(project: MO #Unresolved) and Severity: {Very High}", param.GetPriorityQuery())

	param.PriorityFieldLabel = pointer.ToString("Bug severity")
	assert.Equal(t, "(project: MO #Unresolved) and {Bug severity}: {Very High}", param.GetPriorityQuery())
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
