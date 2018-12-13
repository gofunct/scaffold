package scaffold

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

var srcPaths []string

func init() {
	// Initialize srcPaths.
	envGoPath := os.Getenv("GOPATH")
	goPaths := filepath.SplitList(envGoPath)
	if len(goPaths) == 0 {


		goExecutable := os.Getenv("COBRA_GO_EXECUTABLE")
		if len(goExecutable) <= 0 {
			goExecutable = "go"
		}

		out, Err := exec.Command(goExecutable, "env", "GOPATH").Output()
		if Err != nil {
			Er(Err)
		}

		toolchainGoPath := strings.TrimSpace(string(out))
		goPaths = filepath.SplitList(toolchainGoPath)
		if len(goPaths) == 0 {
			Er("$GOPATH is not set")
		}
	}
	srcPaths = make([]string, 0, len(goPaths))
	for _, goPath := range goPaths {
		srcPaths = append(srcPaths, filepath.Join(goPath, "src"))
	}
}

func Er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

// IsEmpty checks if a given path is empty.
// Hidden files in path are ignored.
func IsEmpty(path string) bool {
	fi, Err := os.Stat(path)
	if Err != nil {
		Er(Err)
	}

	if !fi.IsDir() {
		return fi.Size() == 0
	}

	f, Err := os.Open(path)
	if Err != nil {
		Er(Err)
	}
	defer f.Close()

	names, Err := f.Readdirnames(-1)
	if Err != nil && Err != io.EOF {
		Er(Err)
	}

	for _, name := range names {
		if len(name) > 0 && name[0] != '.' {
			return false
		}
	}
	return true
}

// Exists checks if a file or directory Exists.
func Exists(path string) bool {
	if path == "" {
		return false
	}
	_, Err := os.Stat(path)
	if Err == nil {
		return true
	}
	if !os.IsNotExist(Err) {
		Er(Err)
	}
	return false
}

func ExecuteTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, Err := template.New("").Funcs(template.FuncMap{"comment": CommentifyString}).Parse(tmplStr)
	if Err != nil {
		return "", Err
	}

	buf := new(bytes.Buffer)
	Err = tmpl.Execute(buf, data)
	return buf.String(), Err
}

func WriteStringToFile(path string, s string) error {
	return WriteToFile(path, strings.NewReader(s))
}

// WriteToFile writes r to file with path only
// if file/directory on given path doesn't exist.
func WriteToFile(path string, r io.Reader) error {
	if Exists(path) {
		return fmt.Errorf("%v already Exists", path)
	}

	dir := filepath.Dir(path)
	if dir != "" {
		if Err := os.MkdirAll(dir, 0777); Err != nil {
			return Err
		}
	}

	file, Err := os.Create(path)
	if Err != nil {
		return Err
	}
	defer file.Close()

	_, Err = io.Copy(file, r)
	return Err
}

// commentfyString comments evEry line of in.
func CommentifyString(in string) string {
	var newlines []string
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			newlines = append(newlines, line)
		} else {
			if line == "" {
				newlines = append(newlines, "//")
			} else {
				newlines = append(newlines, "// "+line)
			}
		}
	}
	return strings.Join(newlines, "\n")
}

func ValidateCmdName(source string) string {
	i := 0
	l := len(source)
	// The output is initialized on demand, then first dash or underscore
	// occurs.
	var output string

	for i < l {
		if source[i] == '-' || source[i] == '_' {
			if output == "" {
				output = source[:i]
			}

			// If it's last rune and it's dash or underscore,
			// don't add it output and break the loop.
			if i == l-1 {
				break
			}

			// If next character is dash or underscore,
			// just skip the current character.
			if source[i+1] == '-' || source[i+1] == '_' {
				i++
				continue
			}

			// If the current character is dash or underscore,
			// upper next letter and add to output.
			output += string(unicode.ToUpper(rune(source[i+1])))
			// We know, what source[i] is dash or underscore and source[i+1] is
			// uppered character, so make i = i+2.
			i += 2
			continue
		}

		// If the current character isn't dash or underscore,
		// just add it.
		if output != "" {
			output += string(source[i])
		}
		i++
	}

	if output == "" {
		return source // source is initially valid name.
	}
	return output
}

