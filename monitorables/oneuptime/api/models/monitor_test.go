package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitor_MatchOneLabel(t *testing.T) {
	monitor := &Monitor{Labels: []string{"label1", "label2"}}
	assert.True(t, monitor.MatchOneLabel([]string{"label1"}))
	assert.True(t, monitor.MatchOneLabel([]string{"label1", "label2"}))
	assert.True(t, monitor.MatchOneLabel([]string{"label1", "label3"}))
	assert.False(t, monitor.MatchOneLabel([]string{"label3"}))
	assert.False(t, monitor.MatchOneLabel([]string{"label3", "label4"}))
}
