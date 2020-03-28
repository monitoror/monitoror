package cli

import (
	"fmt"
	"strings"

	"github.com/monitoror/monitoror/cli/version"
	coreModels "github.com/monitoror/monitoror/models"

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

`
	devModeTitle   = ` DEVELOPMENT MODE `
	uiStartCommand = `yarn serve`
	devMode        = `
┌─%s──────────────────────────────┐
│ UI must be started via %s from ./ui     │
│ For more details, check our development guide:  │
│ %s       │
└─────────────────────────────────────────────────┘

`

	monitorableHeader = `
ENABLED MONITORABLES
`

	monitorableFooterTitle = `%d more monitorables were ignored`
	monitorableFooter      = `

%s
Check the documentation to know how to enabled them:
%s
`
	echoStartup = `

Monitoror is running at:
`
)

var colorer = color.New()

func PrintBanner() {
	var tagFlag = ""
	if len(version.BuildTags) > 0 {
		tagFlag = colorer.Inverse(" " + version.BuildTags + " ")
	}

	colorer.Printf(banner, tagFlag, colorer.Green(version.Version), colorer.Blue(website))
}

func PrintDevMode() {
	colorer.Printf(devMode, colorer.Yellow(devModeTitle), colorer.Green(uiStartCommand), colorer.Blue(developmentGuides))
}

func PrintMonitorableHeader() {
	colorer.Println(colorer.Black(colorer.Green(monitorableHeader)))
}

func PrintMonitorable(displayName string, enabledVariants []coreModels.VariantName, erroredVariants map[coreModels.VariantName]error) {
	if len(enabledVariants) == 0 && len(erroredVariants) == 0 {
		return
	}

	// Stringify variants
	var strVariants string
	if len(enabledVariants) == 1 && enabledVariants[0] == coreModels.DefaultVariant {
		if len(erroredVariants) > 0 {
			strVariants = "[default]"
		}
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
		if strDefault != "" || len(variantsWithoutDefault) > 0 {
			strVariants = fmt.Sprintf("[%svariants: [%s]]", strDefault, strings.Join(variantsWithoutDefault, ", "))
		}
	}

	// Print Monitorable and variants
	prefixStatus := colorer.Green("✓")
	if len(erroredVariants) > 0 {
		if len(enabledVariants) > 0 {
			prefixStatus = colorer.Yellow("!")
		} else {
			prefixStatus = colorer.Red("✕")
		}
	}
	monitorableName := strings.Replace(displayName, "(faker)", colorer.Grey("(faker)"), 1)
	colorer.Printf("  %s %s %s\n", prefixStatus, monitorableName, colorer.Grey(strVariants))

	// Print Error
	for variant, err := range erroredVariants {
		if variant == coreModels.DefaultVariant {
			colorer.Printf(colorer.Red("    %s Errored %s configuration\n        %s\n"), errorSymbol, variant, err.Error())
		} else {
			colorer.Printf(colorer.Red("    %s Errored \"%s\" variant configuration\n        %s\n"), errorSymbol, variant, err.Error())
		}
	}
}

func PrintMonitorableFooter(isProduction bool, nonEnabledMonitorableCount int) {
	var documentationVersion string
	if isProduction {
		if splittedVersion := strings.Split(version.Version, "."); len(splittedVersion) == 3 {
			documentationVersion = fmt.Sprintf("%s.%s/", splittedVersion[0], splittedVersion[1])
		}
	}

	coloredMonitororFooterTitle := colorer.Yellow(fmt.Sprintf(monitorableFooterTitle, nonEnabledMonitorableCount))
	colorer.Printf(monitorableFooter, coloredMonitororFooterTitle, colorer.Blue(fmt.Sprintf(documentation, documentationVersion)))
}

func PrintServerStartup(ip string, port int) {
	colorer.Printf(echoStartup)
	colorer.Printf("  %s\n", colorer.Blue(fmt.Sprintf("http://localhost:%d", port)))
	colorer.Printf("  %s\n", colorer.Blue(fmt.Sprintf("http://%s:%d", ip, port)))
	colorer.Println()
}
