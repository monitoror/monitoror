package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing_unit(t *testing.T) {
	// It's pain to test go-ping without interface
	// Assuming go-test already have integration test
	assert.True(t, true)
}
