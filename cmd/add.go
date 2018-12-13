package cmd

import (
	"fmt"
	"github.com/gofunct/scaffold/engine"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func init() {
	AddCmd.Flags().StringVarP(&packageName, "package", "t", "", "target package name")
	AddCmd.Flags().StringVarP(&parentName, "parent", "p", "RootCmd", "variable name of parent command for this command")
}

var packageName, parentName string

var AddCmd = &cobra.Command{
	Use:     "add [command name]",
	Aliases: []string{"command"},
	Short:   "Add a command to a GoRPC Application",
	Long: `Add (gorpc add) will create a new command, with a license and
the appropriate structure for a Cobra-based CLI application,
and register it to its parent (default RootCmd).

If you want your command to be public, pass in the command name
with an initial uppercase letter.

Example: gorpc add server -> resulting in a new cmd/server.go`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			scaffold.Er("add needs a name for the command")
		}

		var project *scaffold.Project
		if packageName != "" {
			project = scaffold.NewProject(packageName)
		} else {
			wd, err := os.Getwd()
			if err != nil {
				scaffold.Er(err)
			}
			project = scaffold.NewProjectFromPath(wd)
		}

		cmdName := scaffold.ValidateCmdName(args[0])
		cmdPath := filepath.Join(project.CmdPath(), cmdName+".go")
		scaffold.CreateCmdFile(cmdPath, cmdName)

		fmt.Fprintln(cmd.OutOrStdout(), cmdName, "created at", cmdPath)
	},
}
