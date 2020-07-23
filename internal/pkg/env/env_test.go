package env

import (
	"os"
	"testing"

	"github.com/monitoror/monitoror/models"

	"github.com/stretchr/testify/assert"
)

type TestEnv struct {
	Value string
}

func TestAnalyse_DefaultOverride(t *testing.T) {
	_ = os.Setenv("VALUE1_TEST", "test")
	defaultVariant := string(models.DefaultVariantName)

	variantNames := InitEnvDefaultLabel("VALUE1", "TEST", defaultVariant)

	assert.Len(t, variantNames, 1)
	assert.True(t, variantNames[defaultVariant])
	assert.Equal(t, "test", os.Getenv("VALUE1_DEFAULT_TEST"))
	assert.Equal(t, "", os.Getenv("VALUE1_TEST"))
}

func TestAnalyse_WithoutSuffix(t *testing.T) {
	_ = os.Setenv("VALUE2", "test")
	defaultVariant := string(models.DefaultVariantName)

	variantNames := InitEnvDefaultLabel("VALUE2", "", defaultVariant)

	assert.Len(t, variantNames, 1)
	assert.True(t, variantNames[defaultVariant])
	assert.Equal(t, "test", os.Getenv("VALUE2_DEFAULT"))
	assert.Equal(t, "", os.Getenv("VALUE2"))
}

func TestAnalyse_SimpleLabel(t *testing.T) {
	_ = os.Setenv("VALUE3_VARIANT1_TEST", "test")
	defaultVariant := string(models.DefaultVariantName)

	variantNames := InitEnvDefaultLabel("VALUE3", "TEST", defaultVariant)

	assert.Len(t, variantNames, 2)
	assert.True(t, variantNames[defaultVariant])
	assert.True(t, variantNames["variant1"])
	assert.Equal(t, "test", os.Getenv("VALUE3_VARIANT1_TEST"))
}

func TestAddDefaultLabel_WithConflict(t *testing.T) {
	_ = os.Setenv("VALUE4", "test")
	_ = os.Setenv("VALUE4_DEFAULT", "test")

	defaultVariant := string(models.DefaultVariantName)

	_ = InitEnvDefaultLabel("VALUE4", "", defaultVariant)

	assert.Equal(t, "test", os.Getenv("VALUE4"))
	assert.Equal(t, "test", os.Getenv("VALUE4_DEFAULT"))
}

func TestAddDefaultLabel_WithEqual(t *testing.T) {
	_ = os.Setenv("VALUE5", "a=b")

	defaultVariant := string(models.DefaultVariantName)

	_ = InitEnvDefaultLabel("VALUE5", "", defaultVariant)

	assert.Equal(t, "a=b", os.Getenv("VALUE5_DEFAULT"))
}
