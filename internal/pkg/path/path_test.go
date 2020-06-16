package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToAbsolute(t *testing.T) {
	for _, testcase := range []struct {
		path string
		abs  string
	}{
		{path: "config.json", abs: "/root/config.json"},
		{path: "./config.json", abs: "/root/config.json"},
		{path: "/dir1/config.json", abs: "/dir1/config.json"},
		{path: "dir1/../config.json", abs: "/root/config.json"},
		{path: "/dir1/../config.json", abs: "/config.json"},
	} {
		assert.Equal(t, testcase.abs, ToAbsolute("/root", testcase.path))
	}
}
