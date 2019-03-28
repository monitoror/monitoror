package service

import (
	"github.com/jsdidierlaurent/monitowall/cli/version"
	"github.com/labstack/gommon/color"
)

const (
	website = "https://github.com/jsdidierlaurent/monitowall"
	banner  = `
    __  ___            _ __                      ____
   /  |/  /___  ____  (_) /_____ _      ______  / / /
  / /|_/ / __ \/ __ \/ / __/ __ \ | /| / / __ \/ / /
 / /  / / /_/ / / / / / /_/ /_/ / |/ |/ / /_/ / / /
/_/  /_/\____/_/ /_/_/\__/\____/|__/|__/\__,_/_/_/ %s

%s
___________________________________________________________
                                    
`
)

func PrintBanner() {
	colorer := color.New()
	colorer.Printf(banner, colorer.Red(version.Version), colorer.Blue(website))
}
