package engine

import (
	"path"
	"path/filepath"
)

func Gen(t string, f string, p *Project) {
		data := make(map[string]interface{})
		data["appName"] = path.Base(p.Name())
		data["importpath"] = path.Join(p.Name(), filepath.Base(p.CmdPath()))

		script, err := ExecuteTemplate(t, data)
		if err != nil {
			Er(err)
		}

	err = WriteStringToFile(filepath.Join(p.AbsPath(), f), script)
	if err != nil {
		Er(err)
	}
}


func GenCmd(t, name string, p *Project) {
	data := make(map[string]interface{})
	data["cmdPackage"] = filepath.Base(filepath.Dir("cmd/")) // last dir of path
	data["cmdName"] = name

	cmdScript, err := ExecuteTemplate(t, data)
	if err != nil {
		Er(err)
	}
	err = WriteStringToFile(filepath.Join(p.CmdPath(), name+".go"), cmdScript)
	if err != nil {
		Er(err)
	}
}

