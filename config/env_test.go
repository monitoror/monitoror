package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyse_DefaultOverride(t *testing.T) {
	_ = os.Setenv("MO_MONITORABLE_TRAVISCI_URL", "test")

	labels := initEnvAndVariant()

	assert.Len(t, labels["TravisCI"], 1)
	assert.True(t, labels["TravisCI"][DefaultVariant])
	assert.Equal(t, "test", os.Getenv("MO_MONITORABLE_TRAVISCI_DEFAULT_URL"))
	assert.Equal(t, "", os.Getenv("MO_MONITORABLE_TRAVISCI_URL"))
}

func TestAnalyse_SimpleLabel(t *testing.T) {
	_ = os.Setenv("MO_MONITORABLE_JENKINS_VARIANT1_URL", "test")

	labels := initEnvAndVariant()

	assert.Len(t, labels["Jenkins"], 2)
	assert.True(t, labels["Jenkins"]["variant1"])
}

func TestAddDefaultLabel_WithConflict(t *testing.T) {
	_ = os.Setenv("MO_TEST", "test")
	_ = os.Setenv("MO_DEFAULT_TEST", "test")

	addDefaultVariant("MO_TEST", "MO_DEFAULT_TEST", "test")

	assert.Equal(t, "test", os.Getenv("MO_TEST"))
	assert.Equal(t, "test", os.Getenv("MO_DEFAULT_TEST"))
}
