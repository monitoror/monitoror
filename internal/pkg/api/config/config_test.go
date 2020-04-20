package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	keys := Keys(map[string]string{"test": "test", "test2": "test"})
	assert.Contains(t, keys, ", ")
	assert.Contains(t, keys, "test")
	assert.Contains(t, keys, "test2")
}

func TestStringify(t *testing.T) {
	test := struct {
		Test  string `json:"test"`
		Test2 int    `json:"test2"`
		Test3 *int   `json:"test3"`
	}{
		Test:  "test",
		Test2: 1000,
	}

	assert.Equal(t, `{"test":"test","test2":1000}`, Stringify(test))
}
