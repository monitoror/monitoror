// Based on https://github.com/docker/cli/blob/master/templates/templates.go

package templates

import (
	"strings"
	"text/template"

	"github.com/labstack/gommon/color"
)

var colorer = color.New()

var basicFunctions = template.FuncMap{
	"split": strings.Split,
	"join":  strings.Join,
	"lower": strings.ToLower,
	"upper": strings.ToUpper,

	// For terminal only
	"blue":         colorer.Blue,
	"green":        colorer.Green,
	"red":          colorer.Red,
	"yellow":       colorer.Yellow,
	"grey":         colorer.Grey,
	"inverseColor": colorer.Inverse,
}

// New creates a new empty template with the provided tag and built-in
func New(tag string) *template.Template {
	return template.New(tag).Funcs(basicFunctions)
}
