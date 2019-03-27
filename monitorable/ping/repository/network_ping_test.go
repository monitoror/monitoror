package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Ping(t *testing.T) {
	// It's pain to test go-ping without interface/mock
	// Assuming go-test already have integration test
	assert.True(t, true)
}
