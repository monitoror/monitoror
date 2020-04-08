package cli

import (
	"bytes"
	"errors"
	"testing"

	"github.com/monitoror/monitoror/cli/version"
	"github.com/stretchr/testify/assert"

	coreModels "github.com/monitoror/monitoror/models"
)

func TestPrintBanner(t *testing.T) {
	cli := New()
	output := &bytes.Buffer{}
	colorer.SetOutput(output)
	var actual string
	var expected string

	// Without BuildTags (default)
	output.Reset()
	version.Version = "0.0.0"
	cli.PrintBanner()
	actual = output.String()
	expected = `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  0.0.0

https://monitoror.com

`
	assert.Equal(t, expected, actual)

	// With BuildTags
	output.Reset()
	version.BuildTags = "test-tag"
	version.Version = "0.1.2"
	cli.PrintBanner()
	actual = output.String()
	expected = `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / /  test-tag ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  0.1.2

https://monitoror.com

`
	assert.Equal(t, expected, actual)
}

func TestPrintDevMode(t *testing.T) {
	cli := New()
	output := &bytes.Buffer{}
	colorer.SetOutput(output)
	cli.PrintDevMode()
	actual := output.String()
	expected := `
┌─ DEVELOPMENT MODE ──────────────────────────────┐
│ UI must be started via yarn serve from ./ui     │
│ For more details, check our development guide:  │
│ https://monitoror.com/guides/#development       │
└─────────────────────────────────────────────────┘

`
	assert.Equal(t, expected, actual)
}

func TestPrintMonitorableHeader(t *testing.T) {
	cli := New()
	output := &bytes.Buffer{}
	colorer.SetOutput(output)
	cli.PrintMonitorableHeader()
	actual := output.String()
	expected := `
ENABLED MONITORABLES

`
	assert.Equal(t, expected, actual)
}

func TestPrintMonitorable(t *testing.T) {
	cli := New()
	output := &bytes.Buffer{}
	colorer.SetOutput(output)
	var actual string
	var expected string

	// Not configured monitorable
	output.Reset()
	cli.PrintMonitorable("TEST1", nil, nil)
	assert.Equal(t, ``, output.String())

	// Default is configured
	output.Reset()
	cli.PrintMonitorable("TEST2", []coreModels.VariantName{coreModels.DefaultVariant}, nil)
	actual = output.String()
	expected = `  ✓ TEST2 ` + `
`
	assert.Equal(t, expected, actual)

	// A variant is configured
	output.Reset()
	cli.PrintMonitorable("TEST2bis", []coreModels.VariantName{"variant7"}, nil)
	actual = output.String()
	expected = `  ✓ TEST2bis [variants: [variant7]]
`
	assert.Equal(t, expected, actual)

	// Default is configured, a variant is errored
	output.Reset()
	cli.PrintMonitorable("TEST3", []coreModels.VariantName{coreModels.DefaultVariant}, []ErroredVariant{{"variant2", errors.New("config error details")}})
	actual = output.String()
	expected = `  ! TEST3 [default]
    /!\ Errored "variant2" variant configuration
        config error details
`
	assert.Equal(t, expected, actual)

	// Default and a variant are configured, a variant is errored
	output.Reset()
	cli.PrintMonitorable("TEST3bis", []coreModels.VariantName{coreModels.DefaultVariant, "variant1"}, []ErroredVariant{{"variant2", errors.New("config error details")}})
	actual = output.String()
	expected = `  ! TEST3bis [default, variants: [variant1]]
    /!\ Errored "variant2" variant configuration
        config error details
`
	assert.Equal(t, expected, actual)

	// Default not configured, a variant is errored
	output.Reset()
	cli.PrintMonitorable("TEST4", nil, []ErroredVariant{{"variant3", errors.New("config error details")}})
	actual = output.String()
	expected = `  ✕ TEST4 ` + `
    /!\ Errored "variant3" variant configuration
        config error details
`
	assert.Equal(t, expected, actual)

	// Default is errored, no other variants
	output.Reset()
	cli.PrintMonitorable("TEST5", nil, []ErroredVariant{{coreModels.DefaultVariant, errors.New("boom")}})
	actual = output.String()
	expected = `  ✕ TEST5 ` + `
    /!\ Errored default configuration
        boom
`
	assert.Equal(t, expected, actual)

	// Default is errored, a variant is configured
	output.Reset()
	cli.PrintMonitorable("TEST6 (faker)", []coreModels.VariantName{"variant1"}, []ErroredVariant{{coreModels.DefaultVariant, errors.New("boom")}})
	actual = output.String()
	expected = `  ! TEST6 (faker) [variants: [variant1]]
    /!\ Errored default configuration
        boom
`
	assert.Equal(t, expected, actual)

	// Multiple errored variants
	output.Reset()
	cli.PrintMonitorable("TEST7", []coreModels.VariantName{"variant1"}, []ErroredVariant{{coreModels.DefaultVariant, errors.New("boom")}, {"errored", errors.New("bim")}})
	actual = output.String()
	expected = `  ! TEST7 [variants: [variant1]]` + `
    /!\ Errored default configuration
        boom
    /!\ Errored "errored" variant configuration
        bim
`
	assert.Equal(t, expected, actual)
}

func TestPrintMonitorableFooter(t *testing.T) {
	cli := New()
	output := &bytes.Buffer{}
	colorer.SetOutput(output)

	// Production: documentation with version in URL
	output.Reset()
	version.Version = "1.2.3"
	cli.PrintMonitorableFooter(true, 42)
	actual := output.String()
	expected := `

42 more monitorables were ignored
Check the documentation to know how to enabled them:
https://monitoror.com/1.2/documentation/
`
	assert.Equal(t, expected, actual)

	// Development: latest documentation URL (without version)
	output.Reset()
	cli.PrintMonitorableFooter(false, 42)
	actual = output.String()
	expected = `

42 more monitorables were ignored
Check the documentation to know how to enabled them:
https://monitoror.com/documentation/
`
	assert.Equal(t, expected, actual)

	// No more non-enabled monitorables
	output.Reset()
	cli.PrintMonitorableFooter(true, 0)
	actual = output.String()
	expected = ``
	assert.Equal(t, expected, actual)

}

func TestPrintServerStartup(t *testing.T) {
	cli := New()
	output := &bytes.Buffer{}
	colorer.SetOutput(output)
	cli.PrintServerStartup("1.2.3.4", 9999)
	actual := output.String()
	expected := `

Monitoror is running at:
  http://localhost:9999
  http://1.2.3.4:9999

`
	assert.Equal(t, expected, actual)
}
