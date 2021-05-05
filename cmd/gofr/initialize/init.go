package initialize

import (
	"os"
	"strconv"

	"github.com/zopsmart/gofr/cmd/gofr/helper"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Handler struct {
}

func (h Handler) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (h Handler) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (h Handler) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (h Handler) Help() string {
	return helper.Generate(helper.Help{
		Example:     "gofr init -name=testProject",
		Flag:        "name provide the name of the project",
		Usage:       "init -name=<project_name>",
		Description: "creates a project structure inside the directory specified in the name flag",
	})
}

func Init(c *gofr.Context) (interface{}, error) {
	var h Handler

	helpBool, _ := strconv.ParseBool(c.PathParam("h"))
	if helpBool {
		return h.Help(), nil
	}

	projectName := c.PathParam("name")

	err := createProject(h, projectName)
	if err != nil {
		return nil, err
	}

	return "Successfully created project: " + projectName, nil
}

func createProject(f fileSystem, projectName string) error {
	standardDirectories := []string{
		"cmd",
		"configs",
		"internal",
	}

	standardEnvFiles := []string{
		".env",
		".test.env",
	}

	err := f.Mkdir(projectName, os.ModePerm)
	if err != nil {
		return err
	}

	err = f.Chdir(projectName)
	if err != nil {
		return err
	}

	for _, name := range standardDirectories {
		if er := f.Mkdir(name, 0777); er != nil {
			return er
		}
	}

	mainFile, err := f.Create("main.go")
	if err != nil {
		return err
	}

	mainString := `package main

import (
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	// Sample Route
	k.GET("/hello", func(c *gofr.Context) (interface{}, error) {
		return "Hello World!!!", nil
	})

	// Add the routes here

	k.Start()
}
`

	_, err = mainFile.WriteString(mainString)
	if err != nil {
		_ = os.Remove(projectName)
		return err
	}

	err = createEnvFiles(f, standardEnvFiles)
	if err != nil {
		return err
	}

	return nil
}

func createEnvFiles(f fileSystem, envFiles []string) error {
	err := f.Chdir("configs")
	if err != nil {
		return err
	}

	for _, fileName := range envFiles {
		_, err = f.Create(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}
