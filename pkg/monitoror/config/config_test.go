package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	ValueA string
	ValueB bool
	ValueC int
}

func TestLoadConfig(t *testing.T) {
	// SETUP
	prefix := "MO_TEST"
	defaultVariant := "default"
	conf := make(map[string]*TestStruct)
	defaultValue := &TestStruct{ValueA: "default", ValueB: true, ValueC: 1337}

	// ENV
	_ = os.Setenv("MO_TEST_TESTSTRUCT_VALUEA", "test value A")
	_ = os.Setenv("MO_TEST_TESTSTRUCT_VARIANT1_VALUEC", "1000")

	// TEST
	LoadConfigWithVariant(prefix, defaultVariant, &conf, defaultValue)

	assert.Len(t, conf, 2)

	assert.Equal(t, "test value A", conf[defaultVariant].ValueA)
	assert.Equal(t, defaultValue.ValueB, conf[defaultVariant].ValueB)
	assert.Equal(t, defaultValue.ValueC, conf[defaultVariant].ValueC)

	assert.Equal(t, defaultValue.ValueA, conf["variant1"].ValueA)
	assert.Equal(t, defaultValue.ValueB, conf["variant1"].ValueB)
	assert.Equal(t, 1000, conf["variant1"].ValueC)
}

func TestLoadConfig_Panic(t *testing.T) {
	conf := make(map[string]*TestStruct)
	defaultValue := TestStruct{ValueA: "default", ValueB: true, ValueC: 1337}

	assert.Panics(t, func() { LoadConfigWithVariant("", "", nil, nil) })
	assert.Panics(t, func() { LoadConfigWithVariant("", "", conf, &defaultValue) }) // Need pointer
	assert.Panics(t, func() { LoadConfigWithVariant("", "", &conf, defaultValue) }) // Need pointer

	conf2 := make(map[int]TestStruct)
	assert.Panics(t, func() { LoadConfigWithVariant("", "", &conf2, defaultValue) }) // Need *map[string] ...
}
