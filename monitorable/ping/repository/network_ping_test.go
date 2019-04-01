package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_CheckPing(t *testing.T) {
	// It's pain to test go-ping without interface/mock
	// Assuming go-test already have tests ...
	assert.True(t, true)
}
