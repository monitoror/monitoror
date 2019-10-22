package models

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitororError_Error(t *testing.T) {
	me := &MonitororError{Err: errors.New("boom")}
	assert.Equal(t, "boom", me.Error())

	me.Message = "big boom"
	assert.Equal(t, "big boom", me.Error())

	me.Err = nil
	assert.Equal(t, "big boom", me.Error())

	me.Message = ""
	assert.Equal(t, "", me.Error())
}

func TestMonitororError_Unwrap(t *testing.T) {
	err := context.DeadlineExceeded
	me := &MonitororError{Err: err}
	assert.Equal(t, err, me.Unwrap())
}

func TestMonitororError_Timeout(t *testing.T) {
	me := &MonitororError{}
	assert.False(t, me.Timeout())

	me = &MonitororError{Err: &net.DNSError{IsTimeout: true}}
	assert.True(t, me.Timeout())

	me = &MonitororError{Err: fmt.Errorf("boom, %w", &net.DNSError{IsNotFound: true})}
	assert.True(t, me.Timeout())

	me = &MonitororError{Err: errors.New("net/http: request canceled while waiting for connection")}
	assert.True(t, me.Timeout())

	me = &MonitororError{Err: errors.New("boom")}
	assert.False(t, me.Timeout())
}
