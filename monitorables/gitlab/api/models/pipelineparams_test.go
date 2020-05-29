package models

import (
	"fmt"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestPipeline_Validate(t *testing.T) {
	param := &PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"}
	test.AssertParams(t, param, 0)

	param = &PipelineParams{ProjectID: pointer.ToInt(10)}
	test.AssertParams(t, param, 1)

	param = &PipelineParams{Ref: "master"}
	test.AssertParams(t, param, 1)

	param = &PipelineParams{}
	test.AssertParams(t, param, 2)
}

func TestPipelineParams_String(t *testing.T) {
	param := &PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"}
	assert.Equal(t, "PIPELINE-10-master", fmt.Sprint(param))
}
