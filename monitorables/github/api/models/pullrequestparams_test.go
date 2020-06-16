package models

import (
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestPullRequest_Validate(t *testing.T) {
	param := &PullRequestParams{Owner: "test", Repository: "test", ID: pointer.ToInt(10)}
	test.AssertParams(t, param, 0)

	param = &PullRequestParams{Owner: "test", Repository: "test"}
	test.AssertParams(t, param, 1)

	param = &PullRequestParams{Owner: "test"}
	test.AssertParams(t, param, 2)

	param = &PullRequestParams{}
	test.AssertParams(t, param, 3)
}

func TestPullRequestParams_String(t *testing.T) {
	param := &PullRequestParams{Owner: "test", Repository: "test", ID: pointer.ToInt(10)}
	assert.Equal(t, "PULLREQUEST-test-test-10", fmt.Sprint(param))
}
