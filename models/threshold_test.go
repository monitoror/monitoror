package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreshold_GetTileStatus_Empty(t *testing.T) {
	tr := IntThreshold{WarningStatus: 80, FailedStatus: 100}
	assert.Equal(t, SuccessStatus, tr.GetTileStatus(10, SuccessStatus))
	assert.Equal(t, SuccessStatus, tr.GetTileStatus(50, SuccessStatus))
	assert.Equal(t, SuccessStatus, tr.GetTileStatus(-10, SuccessStatus))
	assert.Equal(t, WarningStatus, tr.GetTileStatus(80, SuccessStatus))
	assert.Equal(t, WarningStatus, tr.GetTileStatus(99, SuccessStatus))
	assert.Equal(t, FailedStatus, tr.GetTileStatus(100, SuccessStatus))
	assert.Equal(t, FailedStatus, tr.GetTileStatus(1000, SuccessStatus))
}

func TestTreshold_UnmarshalParam_ErrorInvalidFormat(t *testing.T) {
	var test = IntThreshold{}
	err := test.UnmarshalParam("map[WARNING:5,ERROR]")
	if assert.Error(t, err) {
		assert.Equal(t, "invalid format: ERROR", err.Error())
	}
}

func TestTreshold_UnmarshalParam_ErrorNonNumericValue(t *testing.T) {
	var test = IntThreshold{}
	err := test.UnmarshalParam("map[WARNING:test]")
	if assert.Error(t, err) {
		assert.Equal(t, "non-numeric value for WARNING", err.Error())
	}
}

func TestTreshold_UnmarshalParam_ErrorInvalidTileStatus(t *testing.T) {
	var test = IntThreshold{}
	err := test.UnmarshalParam("map[TEST:1]")
	if assert.Error(t, err) {
		assert.Equal(t, "invalid TileStatus: TEST", err.Error())
	}
}

func TestTreshold_UnmarshalParam_Success(t *testing.T) {
	var test = IntThreshold{}
	err := test.UnmarshalParam("map[WARNING:1]")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, test[WarningStatus])
	}
}

func TestTreshold_UnmarshalParam_Empty(t *testing.T) {
	var test = IntThreshold{}
	err := test.UnmarshalParam("map[]")
	if assert.NoError(t, err) {
		assert.Len(t, test, 0)
	}
}
