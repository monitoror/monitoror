package cli

import (
	"github.com/monitoror/monitoror/cli/version"

	"github.com/labstack/gommon/color"
)

const (
	website = "https://github.com/monitoror/monitoror"
	banner  = `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / %s
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  %s

%s
_____________________________________________________

`
)

func PrintBanner() {
	colorer := color.New()
	colorer.Printf(banner, colorer.Green(version.BuildTags), colorer.Green(version.Version), colorer.Blue(website))
}
