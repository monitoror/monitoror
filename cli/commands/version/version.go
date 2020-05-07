package version

import (
	"runtime"
	"text/template"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/cli/version"
	"github.com/monitoror/monitoror/pkg/templates"

	"github.com/spf13/cobra"
)

const versionTemplate = ` Version:    {{green .Version}}{{ with .BuildTags }}{{printf " (%s)" . | grey }}{{end}}
 Git commit: {{green .GitCommit}}
 Built:      {{green .BuildTime}}

 Go version: {{blue .GoVersion}}
 OS/Arch:    {{blue .Os}}/{{blue .Arch}}`

type versionInfo struct {
	Version   string
	GitCommit string
	GoVersion string
	Os        string
	Arch      string
	BuildTime string
	BuildTags string
}

var parsedTemplate *template.Template

func init() {
	// Print this error when you want to debug template
	parsedTemplate, _ = templates.New("version").Parse(versionTemplate)
}

func NewVersionCommand(monitororCli *cli.MonitororCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the Monitoror version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(monitororCli)
		},
	}
	return cmd
}

func runVersion(monitororCli *cli.MonitororCli) error {
	vi := &versionInfo{
		Version:   version.Version,
		GitCommit: version.GitCommit,
		BuildTime: version.BuildTime,
		BuildTags: version.BuildTags,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	return parsedTemplate.Execute(monitororCli.Output, vi)
}
