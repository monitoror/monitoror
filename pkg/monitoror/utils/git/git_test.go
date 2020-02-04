package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanizeBranch(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%s%s", branchPrefix, "master"), HumanizeBranch("refs/heads/master"))
	assert.Equal(t, fmt.Sprintf("%s%s", branchPrefix, "develop"), HumanizeBranch("refs/heads/develop"))
	assert.Equal(t, fmt.Sprintf("%s%s", branchPrefix, "feat/toto"), HumanizeBranch(fmt.Sprintf("%s%s", branchPrefix, "refs/heads/feat/toto")))
	assert.Equal(t, fmt.Sprintf("%s%s", branchPrefix, "feat/toto"), HumanizeBranch("feat/toto"))
	assert.Equal(t, fmt.Sprintf("%s%s", branchPrefix, "feat/toto"), HumanizeBranch(fmt.Sprintf("%s%s", branchPrefix, "feat/toto")))
}
