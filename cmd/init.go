package cmd

import (
	"fmt"
	"github.com/gofunct/scaffold/engine"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var InitCmd = &cobra.Command{
	Use:     "init [name]",
	Aliases: []string{"initialize", "initialise", "create"},
	Short:   "Initialize a Cobra Application",
	Long: `Initialize (cobra init) will create a new application, with a license
and the appropriate structure for a Cobra-based CLI application.

  * If a name is provided, it will be created in the current directory;
  * If no name is provided, the current directory will be assumed;
  * If a relative path is provided, it will be created inside $GOPATH
    (e.g. github.com/spf13/hugo);
  * If an absolute path is provided, it will be created;
  * If the directory already exists but is empty, it will be used.

Init will not use an existing directory with contents.`,

	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			engine.Er(err)
		}

		var project *engine.Project
		if len(args) == 0 {
			project = engine.NewProjectFromPath(wd)
		} else if len(args) == 1 {
			arg := args[0]
			if arg[0] == '.' {
				arg = filepath.Join(wd, arg)
			}
			if filepath.IsAbs(arg) {
				project = engine.NewProjectFromPath(arg)
			} else {
				project = engine.NewProject(arg)
			}
		} else {
			engine.Er("please provide only one argument")
		}

		engine.InitializeProject(project)

		fmt.Fprintln(cmd.OutOrStdout(), `Your GoRPC application is ready at
`+project.AbsPath()+`

Give it a try by going there and running `+"`go run main.go`."+`
Add commands to it by running `+"`gorpc add [cmdname]`.")
	},
}