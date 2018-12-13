package engine

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Project contains name, license and paths to projects.
type Project struct {
	absPath string
	cmdPath string
	srcPath string
	name    string
}

// NewProject returns Project with specified project name.
func NewProject(projectName string) *Project {
	if projectName == "" {
		Er("can't create project with blank name")
	}

	p := new(Project)
	p.name = projectName

	// 1. Find already created protect.
	p.absPath = FindPackage(projectName)

	// 2. If there are no created project with this path, and user is in GOPATH,
	// then use GOPATH/src/projectName.
	if p.absPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			Er(err)
		}
		for _, srcPath := range srcPaths {
			goPath := filepath.Dir(srcPath)
			if FilePathHasPrefix(wd, goPath) {
				p.absPath = filepath.Join(srcPath, projectName)
				break
			}
		}
	}

	// 3. If user is not in GOPATH, then use (first GOPATH)/src/projectName.
	if p.absPath == "" {
		p.absPath = filepath.Join(srcPaths[0], projectName)
	}

	return p
}

// FindPackage returns full path to existing go package in GOPATHs.
func FindPackage(packageName string) string {
	if packageName == "" {
		return ""
	}

	for _, srcPath := range srcPaths {
		packagePath := filepath.Join(srcPath, packageName)
		if Exists(packagePath) {
			return packagePath
		}
	}

	return ""
}

// NewProjectFromPath returns Project with specified absolute path to
// package.
func NewProjectFromPath(absPath string) *Project {
	if absPath == "" {
		Er("can't create project: absPath can't be blank")
	}
	if !filepath.IsAbs(absPath) {
		Er("can't create project: absPath is not absolute")
	}

	// If absPath is symlink, use its destination.
	fi, err := os.Lstat(absPath)
	if err != nil {
		Er("can't read path info: " + err.Error())
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		path, err := os.Readlink(absPath)
		if err != nil {
			Er("can't read the destination of symlink: " + err.Error())
		}
		absPath = path
	}

	p := new(Project)
	p.absPath = strings.TrimSuffix(absPath, FindCmdDir(absPath))
	p.name = filepath.ToSlash(TrimSrcPath(p.absPath, p.SrcPath()))
	return p
}

func TrimSrcPath(absPath, srcPath string) string {
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil {
		Er(err)
	}
	return relPath
}

func (p Project) Name() string {
	return p.name
}

func (p *Project) CmdPath() string {
	if p.absPath == "" {
		return ""
	}
	if p.cmdPath == "" {
		p.cmdPath = filepath.Join(p.absPath, FindCmdDir(p.absPath))
	}
	return p.cmdPath
}

func FindCmdDir(absPath string) string {
	if !Exists(absPath) || IsEmpty(absPath) {
		return "cmd"
	}

	if IsCmdDir(absPath) {
		return filepath.Base(absPath)
	}

	files, _ := filepath.Glob(filepath.Join(absPath, "c*"))
	for _, file := range files {
		if IsCmdDir(file) {
			return filepath.Base(file)
		}
	}

	return "cmd"
}

// IsCmdDir checks if base of name is one of cmdDir.
func IsCmdDir(name string) bool {
	name = filepath.Base(name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if name == cmdDir {
			return true
		}
	}
	return false
}

// AbsPath returns absolute path of project.
func (p Project) AbsPath() string {
	return p.absPath
}

// SrcPath returns absolute path to $GOPATH/src where project is located.
func (p *Project) SrcPath() string {
	if p.srcPath != "" {
		return p.srcPath
	}
	if p.absPath == "" {
		p.srcPath = srcPaths[0]
		return p.srcPath
	}

	for _, srcPath := range srcPaths {
		if FilePathHasPrefix(p.absPath, srcPath) {
			p.srcPath = srcPath
			break
		}
	}

	return p.srcPath
}

func FilePathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}
	return path[0:len(prefix)] == prefix

}
