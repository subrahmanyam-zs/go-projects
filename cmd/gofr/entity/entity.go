package entity

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"developer.zopsmart.com/go/gofr/cmd/gofr/helper"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Handler struct {
}

func (h Handler) Getwd() (string, error) {
	return os.Getwd()
}

func (h Handler) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (h Handler) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (h Handler) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (h Handler) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (h Handler) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (h Handler) Help() string {
	return helper.Generate(helper.Help{
		Example: "gofr entity -type=core -name=persons",
		Flag: `type specify the layer: core, composite or consumer
name entity name`,
		Usage:       "entity -type=<layer> -name=<entity_name>",
		Description: "creates a template and interface for an entity",
	})
}

func AddEntity(c *gofr.Context) (interface{}, error) {
	var h Handler

	helpBool, _ := strconv.ParseBool(c.PathParam("h"))
	if helpBool {
		return h.Help(), nil
	}

	layer := c.PathParam("type")
	name := c.PathParam("name")

	err := addEntity(h, layer, name)
	if err != nil {
		return nil, err
	}

	return "Successfully created entity: " + name, nil
}

type invalidTypeError struct{}

func (i invalidTypeError) Error() string {
	return "invalid type"
}

func addEntity(f fileSystem, entityType, entity string) error {
	projectDirectory, err := f.Getwd()
	if err != nil {
		return err
	}

	switch entityType {
	case "core":
		err := addCore(f, projectDirectory, entity)
		if err != nil {
			return err
		}
	case "composite":
		err := addComposite(f, projectDirectory, entity)
		if err != nil {
			return err
		}
	case "consumer":
		err := addConsumer(f, projectDirectory, entity)
		if err != nil {
			return err
		}
	default:
		return invalidTypeError{}
	}

	return nil
}

func addCore(f fileSystem, projectDirectory, entity string) error {
	path := projectDirectory + "/core"

	err := createChangeDir(f, path)
	if err != nil {
		return err
	}
	// create the interfaceFile , interface.go,  for core layer
	interfaceFile, err := f.OpenFile("interface.go", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer interfaceFile.Close()

	err = populateInterfaceFiles(strings.Title(entity), projectDirectory, "cores", interfaceFile)
	if err != nil {
		return err
	}

	entityPath := path + "/" + entity

	err = createChangeDir(f, entityPath)
	if err != nil {
		return err
	}

	err = populateEntityFile(f, projectDirectory, entityPath, entity, "core")
	if err != nil {
		return err
	}

	err = createModel(f, projectDirectory, entity)
	if err != nil {
		return err
	}

	return nil
}

func addComposite(f fileSystem, projectDirectory, entity string) error {
	compositePath := projectDirectory + "/composite"

	err := createChangeDir(f, compositePath)
	if err != nil {
		return err
	}

	interfaceFile, err := f.OpenFile("interface.go", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer interfaceFile.Close()

	err = populateInterfaceFiles(strings.Title(entity), projectDirectory, "composites", interfaceFile)
	if err != nil {
		return err
	}

	err = createChangeDir(f, compositePath+"/"+entity)
	if err != nil {
		return err
	}

	err = populateEntityFile(f, projectDirectory, compositePath+"/"+entity, entity, "composite")
	if err != nil {
		return err
	}

	return nil
}

func addConsumer(f fileSystem, projectDirectory, entity string) error {
	path := projectDirectory + "/http"

	err := createChangeDir(f, path)
	if err != nil {
		return err
	}

	err = createChangeDir(f, entity)
	if err != nil {
		return err
	}

	filePtr, err := f.OpenFile(entity+".go", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer filePtr.Close()

	_, err = filePtr.WriteString(fmt.Sprintf("package %s", entity))
	if err != nil {
		return err
	}

	return nil
}
