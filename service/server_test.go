package service

import (
	"testing"

	"github.com/monitoror/monitoror/config"

	"github.com/GeertJohan/go.rice/embedded"
	"github.com/stretchr/testify/assert"
)

// /!\ this is an integration test /!\
// Note : It may be necessary to separate them from unit tests

func TestInit_Dev(t *testing.T) {
	assert.NotPanics(t, func() {
		Init(&config.Config{Env: "develop"})
	})
}

func TestInit_Prod_WithoutRicebox(t *testing.T) {
	delete(embedded.EmbeddedBoxes, "../ui/dist")
	assert.Panics(t, func() {
		Init(&config.Config{Env: "production"})
	})
}

func TestInit_Prod_WithRicebox(t *testing.T) {
	delete(embedded.EmbeddedBoxes, "../ui/dist")
	embedded.RegisterEmbeddedBox("../ui/dist", &embedded.EmbeddedBox{
		Name: "../ui/dist",
	})
	defer delete(embedded.EmbeddedBoxes, "../ui/dist")
	assert.NotPanics(t, func() {
		Init(&config.Config{Env: "production"})
	})
}
