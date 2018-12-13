package hack

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/mattn/go-colorable"
	"github.com/xlab/treeprint"
	"go/format"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"text/template"
)

var cfgFile string

var (
	gopath       string
	templatePath string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scaffold.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".scaffold" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".scaffold")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func init() {
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		b, err := exec.Command("go", "env", "GOPATH").CombinedOutput()
		if err != nil {
			panic(string(b))
		}
		gopath = strings.TrimSpace(string(b))
	}

	if paths := filepath.SplitList(gopath); len(paths) > 0 {
		gopath = paths[0]
	}

	templatePath = filepath.Clean(filepath.Join(gopath, "/src/github.com/lileio/lile/template"))
	RootCmd.AddCommand(newCmd)
}

var out = colorable.NewColorableStdout()

func new(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Printf("You must supply a path for the service, e.g lile new lile/users\n")
		return
	}

	name := args[0]
	path := projectPath(name)
	fmt.Printf("Creating project in %s\n", path)

	if !askIsOK() {
		fmt.Println("Exiting..")
		return
	}

	p := newProject(path, name)

	err := p.write(templatePath)
	if err != nil {
		er(err)
	}

	p.Folder.print()
}

func askIsOK() bool {
	if os.Getenv("CI") != "" {
		return true
	}

	fmt.Fprintf(out, "Is this OK? %ses/%so\n",
		color.YellowString("[y]"),
		color.CyanString("[n]"),
	)
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	return strings.Contains(strings.ToLower(scan.Text()), "y")
}

func er(err error) {
	if err != nil {
		fmt.Fprintf(out, "%s: %s \n",
			color.RedString("[ERROR]"),
			err.Error(),
		)
		panic(err)
	}
}

type project struct {
	Name         string
	RelativeName string
	ProjectDir   string
	RelDir       string
	Folder       folder
}

func newProject(path, relativeName string) project {
	name := lastFromSplit(path, string(os.PathSeparator))
	relDir := projectBase(path)

	f := folder{Name: name, AbsPath: path}

	s := f.addFolder("server")
	s.addFile("server.go", "server.tmpl")
	s.addFile("server_test.go", "server_test.tmpl")

	subs := f.addFolder("subscribers")
	subs.addFile("subscribers.go", "subscribers.tmpl")

	cmd := f.addFolder(name)
	cmd.addFile("main.go", "cmd_main.tmpl")

	cmds := cmd.addFolder("cmd")
	cmds.addFile("root.go", "cmd_root.tmpl")
	cmds.addFile("up.go", "cmd_up.tmpl")

	f.addFile(name+".proto", "proto.tmpl")
	f.addFile("client.go", "client.tmpl")
	f.addFile("Makefile", "Makefile.tmpl")
	f.addFile("Dockerfile", "Dockerfile.tmpl")
	f.addFile(".gitignore", "gitignore.tmpl")

	return project{
		Name:         name,
		RelativeName: relativeName,
		RelDir:       relDir,
		ProjectDir:   path,
		Folder:       f,
	}
}

func (p project) write(templatePath string) error {
	err := os.MkdirAll(p.ProjectDir, os.ModePerm)
	if err != nil {
		return err
	}

	return p.Folder.render(templatePath, p)
}

// CamelCaseName returns a CamelCased name of the service
func (p project) CamelCaseName() string {
	return strcase.ToCamel(p.Name)
}

// SnakeCaseName returns a snake_cased_type name of the service
func (p project) SnakeCaseName() string {
	return strings.Replace(strcase.ToSnake(p.Name), "-", "_", -1)
}

// DNSName returns a snake-cased-type name that be used as a URL or packageName
func (p project) DNSName() string {
	return strings.Replace(strcase.ToSnake(p.Name), "_", "-", -1)
}

// Copied and re-worked from
// https://github.com/spf13/cobra/blob/master/cobra/cmd/helpers.go
func projectPath(inputPath string) string {
	// if no path is provided... assume CWD.
	if inputPath == "" {
		x, err := os.Getwd()
		if err != nil {
			er(err)
		}

		return x
	}

	var projectPath string
	var projectBase string
	srcPath := srcPath()

	// if provided, inspect for logical locations
	if strings.ContainsRune(inputPath, os.PathSeparator) {
		if filepath.IsAbs(inputPath) || filepath.HasPrefix(inputPath, string(os.PathSeparator)) {
			// if Absolute, use it
			projectPath = filepath.Clean(inputPath)
			return projectPath
		}
		// If not absolute but contains slashes,
		// assuming it means create it from $GOPATH
		count := strings.Count(inputPath, string(os.PathSeparator))

		switch count {
		// If only one directory deep, assume "github.com"
		case 1:
			projectPath = filepath.Join(srcPath, "github.com", inputPath)
			return projectPath
		case 2:
			projectPath = filepath.Join(srcPath, inputPath)
			return projectPath
		default:
			er(errors.New("Unknown directory"))
		}
	}

	// hardest case.. just a word.
	if projectBase == "" {
		x, err := os.Getwd()
		if err == nil {
			projectPath = filepath.Join(x, inputPath)
			return projectPath
		}
		er(err)
	}

	projectPath = filepath.Join(srcPath, projectBase, inputPath)
	return projectPath
}

func projectBase(absPath string) string {
	rel, err := filepath.Rel(srcPath(), absPath)
	if err != nil {
		return filepath.ToSlash(absPath)
	}
	return filepath.ToSlash(rel)
}

func lastFromSplit(input, split string) string {
	rel := strings.Split(input, split)
	return rel[len(rel)-1]
}

func srcPath() string {
	return filepath.Join(gopath, "src") + string(os.PathSeparator)
}

type file struct {
	Name     string
	AbsPath  string
	Template string
}

type folder struct {
	Name    string
	AbsPath string

	// Unexported so you can't set them without methods
	files   []file
	folders []*folder
}

func (f *folder) addFolder(name string) *folder {
	newF := &folder{
		Name:    name,
		AbsPath: filepath.Join(f.AbsPath, name),
	}
	f.folders = append(f.folders, newF)
	return newF
}

func (f *folder) addFile(name, tmpl string) {
	f.files = append(f.files, file{
		Name:     name,
		Template: tmpl,
		AbsPath:  filepath.Join(f.AbsPath, name),
	})
}

func (f *folder) render(templatePath string, p project) error {
	for _, v := range f.files {
		t, err := template.ParseFiles(filepath.Join(templatePath, v.Template))
		if err != nil {
			return err
		}

		file, err := os.Create(v.AbsPath)
		if err != nil {
			return err
		}

		defer file.Close()

		if strings.Contains(v.AbsPath, ".go") {
			var out bytes.Buffer
			err = t.Execute(&out, p)
			if err != nil {
				log.Printf("Could not process template %s\n", v)
				return err
			}

			b, err := format.Source(out.Bytes())
			if err != nil {
				fmt.Print(string(out.Bytes()))
				log.Printf("\nCould not format Go file %s\n", v)
				return err
			}

			_, err = file.Write(b)
			if err != nil {
				return err
			}
		} else {
			err = t.Execute(file, p)
			if err != nil {
				return err
			}
		}
	}

	for _, v := range f.folders {
		err := os.Mkdir(v.AbsPath, os.ModePerm)
		if err != nil {
			return err
		}

		err = v.render(templatePath, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f folder) print() {
	t := f.tree(true, treeprint.New())
	fmt.Println(t.String())
}

func (f folder) tree(root bool, tree treeprint.Tree) treeprint.Tree {
	if !root {
		tree = tree.AddBranch(f.Name)
	}

	for _, v := range f.folders {
		v.tree(false, tree)
	}

	for _, v := range f.files {
		tree.AddNode(v.Name)
	}

	return tree
}
