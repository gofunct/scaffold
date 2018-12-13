package engine

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"github.com/gofunct/scaffold/source"
)

func InitializeProject(project *Project) {
	if !Exists(project.AbsPath()) { // If path doesn't yet exist, create it
		err := os.MkdirAll(project.AbsPath(), os.ModePerm)
		if err != nil {
			Er(err)
		}
	} else if !IsEmpty(project.AbsPath()) { // If path Exists and is not empty don't use it
		Er("GoRPC will not create a new project in a non empty directory: " + project.AbsPath())
	}

	Gen(source.DockerfileTemplate, "Dockerfile", project)
	Gen(source.MakefileTemplate, "Makefile", project)
	Gen(source.MainTemplate, "main.go", project)
	Gen(source.ConfigTemplate, "config.yaml", project)
	Gen(source.GitIgnoreTemplate, ".gitignore", project)
	GenCmd(source.RootTemplate, "root", project)
	GenCmd(source.GrpcServerTemplate, "grpc", project)
}

func InitFunc() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			Er(err)
		}
		var project *Project
		if len(args) == 0 {
			project = NewProjectFromPath(wd)
		} else if len(args) == 1 {
			arg := args[0]
			if arg[0] == '.' {
				arg = filepath.Join(wd, arg)
			}
			if filepath.IsAbs(arg) {
				project = NewProjectFromPath(arg)
			} else {
				project = NewProject(arg)
			}
		} else {
			Er("please provide only one argument")
		}

		InitializeProject(project)

		fmt.Fprintln(cmd.OutOrStdout(), `Your GoRPC application is ready at
`+project.AbsPath()+`

Give it a try by going there and running `+"`go run main.go`."+`
Add commands to it by running `+"`gorpc add [cmdname]`.")
	}
}