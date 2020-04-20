package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestPullRequestGeneratorParams_Validate(t *testing.T) {
	param := &PullRequestGeneratorParams{Owner: "test", Repository: "test"}
	test.AssertParams(t, param, 0)

	param = &PullRequestGeneratorParams{Owner: "test"}
	test.AssertParams(t, param, 1)

	param = &PullRequestGeneratorParams{}
	test.AssertParams(t, param, 2)
}
