package models

import (
	"fmt"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestMergeRequestParams_Validate(t *testing.T) {
	param := &MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)}
	test.AssertParams(t, param, 0)

	param = &MergeRequestParams{ID: pointer.ToInt(10)}
	test.AssertParams(t, param, 1)

	param = &MergeRequestParams{ProjectID: pointer.ToInt(10)}
	test.AssertParams(t, param, 1)

	param = &MergeRequestParams{}
	test.AssertParams(t, param, 2)
}

func TestMergeRequestParams_String(t *testing.T) {
	param := &MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)}
	assert.Equal(t, "MERGEREQUEST-10-10", fmt.Sprint(param))
}
