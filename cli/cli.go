package cli

import (
	"fmt"
	"strings"

	"github.com/monitoror/monitoror/cli/version"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/pkg/system"

	"github.com/labstack/gommon/color"
)

const (
	website           = "https://monitoror.com"
	developmentGuides = "https://monitoror.com/guides/#development"
	documentation     = "https://monitoror.com/" + "%s" + "documentation/"

	errorSymbol = `/!\`

	banner = `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / %s
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  %s

%s
_____________________________________________________
`
	devMode = `
You are in dev mode

To see UI, run ` + "`yarn serve`" + ` in the ` + "`ui`" + ` folder.
For more details, read the development guide:
%s
_____________________________________________________
`

	monitorableHeader = `
Enabled modules
===============
`

	monitorableFooter = `
The module you need is not listed here?
Please read the module documentation to discover all available modules and see how to enable them:
%s
_____________________________________________________
`
	echoStartup = `
Monitoror is running at:
`
)

var colorer = color.New()

func PrintBanner() {
	colorer.Printf(banner, colorer.Green(version.BuildTags), colorer.Green(version.Version), colorer.Blue(website))
}

func PrintDevMode() {
	colorer.Printf(devMode, colorer.Blue(developmentGuides))
}

func PrintMonitorableHeader() {
	colorer.Printf(monitorableHeader)
}

func PrintMonitorable(displayName string, enabledVariants []coreModels.Variant, erroredVariants map[coreModels.Variant]error) {
	if len(enabledVariants) == 0 && len(erroredVariants) == 0 {
		return
	}

	// Stringify variants
	var strVariants string
	if len(enabledVariants) == 1 && enabledVariants[0] == coreModels.DefaultVariant {
		// Only Default variant, skip
		strVariants = ""
	} else {
		var strDefault string
		var variantsWithoutDefault []string

		for _, variant := range enabledVariants {
			if variant == coreModels.DefaultVariant {
				strDefault = fmt.Sprintf("%s, ", variant)
			} else {
				variantsWithoutDefault = append(variantsWithoutDefault, string(variant))
			}
		}
		strVariants = fmt.Sprintf("[%svariants: [%s]]", strDefault, strings.Join(variantsWithoutDefault, ", "))
	}

	// Print Minitorable and variants
	colorer.Printf("- %s %s\n", colorer.Green(displayName), strVariants)

	// Print Error
	for variant, err := range erroredVariants {
		if variant == coreModels.DefaultVariant {
			colorer.Printf(" %[1]s Errored %[2]s configuration %[1]s\n     %[3]s\n", colorer.Red(errorSymbol), variant, err.Error())

		} else {
			colorer.Printf(" %[1]s Errored %[2]s configuration variant %[1]s\n     %[3]s\n", colorer.Red(errorSymbol), variant, err.Error())
		}
	}
}

func PrintMonitorableFooter(isProduction bool) {
	var documentationVersion string
	if isProduction {
		if splittedVersion := strings.Split(version.Version, "."); len(splittedVersion) == 3 {
			documentationVersion = fmt.Sprintf("%s.%s/", splittedVersion[0], splittedVersion[1])
		}
	}

	colorer.Printf(monitorableFooter, colorer.Blue(fmt.Sprintf(documentation, documentationVersion)))
}

func PrintServerStartup(port int) {
	color.Print(echoStartup)

	ips, _ := system.ListLocalhostIpv4()

	// in case of empty ips
	if len(ips) == 0 {
		ips = append(ips, "127.0.0.1")
	}

	for _, ip := range ips {
		if ip == "127.0.0.1" {
			ip = "localhost"
		}

		color.Printf("   %s\n", colorer.Blue(fmt.Sprintf("http://%s:%d", ip, port)))
	}

	color.Println()
}
