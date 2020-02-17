package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanizeBranch(t *testing.T) {
	assert.Equal(t, "master", HumanizeBranch("refs/heads/master"))
	assert.Equal(t, "develop", HumanizeBranch("refs/heads/develop"))
	assert.Equal(t, "feat/toto", HumanizeBranch("refs/heads/feat/toto"))
	assert.Equal(t, "feat/toto", HumanizeBranch("feat/toto"))
	assert.Equal(t, "feat/toto", HumanizeBranch("feat/toto"))
}
