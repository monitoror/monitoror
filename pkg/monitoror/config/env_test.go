package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestEnv struct {
	Value string
}

func TestAnalyse_DefaultOverride(t *testing.T) {
	_ = os.Setenv("MO_VALUE", "test")

	variants := initEnvAndVariant("MO", "default", reflect.TypeOf(TestEnv{}))

	assert.Len(t, variants, 1)
	assert.True(t, variants["default"])
	assert.Equal(t, "test", os.Getenv("MO_DEFAULT_VALUE"))
	assert.Equal(t, "", os.Getenv("MO_VALUE"))
}

func TestAnalyse_SimpleLabel(t *testing.T) {
	_ = os.Setenv("MO_VARIANT1_VALUE", "test")

	variants := initEnvAndVariant("MO", "default", reflect.TypeOf(TestEnv{}))

	assert.Len(t, variants, 2)
	assert.True(t, variants["default"])
	assert.True(t, variants["variant1"])
	assert.Equal(t, "test", os.Getenv("MO_VARIANT1_VALUE"))
}

func TestAddDefaultLabel_WithConflict(t *testing.T) {
	_ = os.Setenv("MO_VALUE", "test")
	_ = os.Setenv("MO_DEFAULT_VALUE", "test")

	_ = initEnvAndVariant("MO", "default", reflect.TypeOf(TestEnv{}))

	assert.Equal(t, "test", os.Getenv("MO_VALUE"))
	assert.Equal(t, "test", os.Getenv("MO_DEFAULT_VALUE"))
}
