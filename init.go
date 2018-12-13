package scaffold

import (
	"os"
	"path"
	"path/filepath"
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

	CreateMainFile(project)
	CreateRootCmdFile(project)
	CreateGrpcServerCmdFile(project)
}

func CreateMainFile(project *Project) {
	mainTemplate := MainTemplate
	data := make(map[string]interface{})
	data["importpath"] = path.Join(project.Name(), filepath.Base(project.CmdPath()))

	mainScript, err := ExecuteTemplate(mainTemplate, data)
	if err != nil {
		Er(err)
	}

	err = WriteStringToFile(filepath.Join(project.AbsPath(), "main.go"), mainScript)
	if err != nil {
		Er(err)
	}
}

func CreateRootCmdFile(project *Project) {
	template := RootTemplate

	data := make(map[string]interface{})
	data["appName"] = path.Base(project.Name())

	rootCmdScript, err := ExecuteTemplate(template, data)
	if err != nil {
		Er(err)
	}

	err = WriteStringToFile(filepath.Join(project.CmdPath(), "root.go"), rootCmdScript)
	if err != nil {
		Er(err)
	}

}

func CreateGrpcServerCmdFile(project *Project) {
	template := GrpcServerTemplate
	data := make(map[string]interface{})

	rootCmdScript, err := ExecuteTemplate(template, data)
	if err != nil {
		Er(err)
	}

	err = WriteStringToFile(filepath.Join(project.CmdPath(), "grpc.go"), rootCmdScript)
	if err != nil {
		Er(err)
	}
}

func CreateDockerFile(project *Project) {
	template := DockerfileTemplate

	data := make(map[string]interface{})
	data["appName"] = path.Base(project.Name())

	dockerScript, err := ExecuteTemplate(template, data)
	if err != nil {
		Er(err)
	}
	err = WriteStringToFile(filepath.Join(project.AbsPath(), "Dockerfile"), dockerScript)
	if err != nil {
		Er(err)
	}
}

func CreateProtoFile(project *Project) {
	template := ProtoTemplate

	data := make(map[string]interface{})
	data["appName"] = path.Base(project.Name())

	genScript, err := ExecuteTemplate(template, data)
	if err != nil {
		Er(err)
	}
	err = WriteStringToFile(filepath.Join(project.AbsPath(), "proto.proto"), genScript)
	if err != nil {
		Er(err)
	}
}

func CreateMakeFile(project *Project) {
	template := MakefileTemplate

	data := make(map[string]interface{})
	data["appName"] = path.Base(project.Name())

	makeScript, err := ExecuteTemplate(template, data)
	if err != nil {
		Er(err)
	}
	err = WriteStringToFile(filepath.Join(project.AbsPath(), "Makefile"), makeScript)
	if err != nil {
		Er(err)
	}
}

func CreateConfigFile(project *Project) {
	template := CmdTemplate

	data := make(map[string]interface{})
	data["appName"] = path.Base(project.Name())

	cfgScript, err := ExecuteTemplate(template, data)
	if err != nil {
		Er(err)
	}
	err = WriteStringToFile(filepath.Join(project.AbsPath(), "config.yaml"), cfgScript)
	if err != nil {
		Er(err)
	}
}

func CreateCmdFile(path, cmdName string) {
	template := CmdTemplate

	data := make(map[string]interface{})
	data["cmdPackage"] = filepath.Base(filepath.Dir(path)) // last dir of path
	data["cmdName"] = cmdName

	cmdScript, err := ExecuteTemplate(template, data)
	if err != nil {
		Er(err)
	}
	err = WriteStringToFile(path, cmdScript)
	if err != nil {
		Er(err)
	}
}