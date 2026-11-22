package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestCountParams_Validate(t *testing.T) {
	param := &CountParams{Owner: "test", Repository: "test"}
	test.AssertParams(t, param, 0)

	param = &CountParams{Owner: "test", Repository: "test", Type: "pulls", State: "all", Labels: "bug,ui", Query: "crash"}
	test.AssertParams(t, param, 0)

	param = &CountParams{Owner: "test", Repository: "test", Type: "commits"}
	test.AssertParams(t, param, 1)

	param = &CountParams{Owner: "test", Repository: "test", State: "merged"}
	test.AssertParams(t, param, 1)

	param = &CountParams{Owner: "test"}
	test.AssertParams(t, param, 1)

	param = &CountParams{}
	test.AssertParams(t, param, 2)
}

func TestCountParams_String(t *testing.T) {
	param := &CountParams{Owner: "test", Repository: "repo", Type: "pulls", State: "open", Labels: "bug", Query: "crash"}
	assert.Equal(t, "COUNT-test-repo-pulls-open-bug-crash", fmt.Sprint(param))
}
