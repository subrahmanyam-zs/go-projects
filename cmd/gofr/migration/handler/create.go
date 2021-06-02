package handler

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"developer.zopsmart.com/go/gofr/cmd/gofr/helper"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Create struct {
}

func (c Create) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (c Create) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (c Create) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (c Create) ReadDir(dir string) ([]os.DirEntry, error) {
	return os.ReadDir(dir)
}

func (c Create) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (c Create) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (c Create) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (c Create) Help() interface{} {
	return helper.Generate(helper.Help{
		Example:     `gofr migrate create -name=AddForeignKey`,
		Flag:        `name: name of the migration`,
		Usage:       "gofr migrate create -name=<migration_name>",
		Description: "creates a migration template inside migrations folder and the name of the file is the name provided in the `name` flag",
	})
}

func CreateMigration(c *gofr.Context) (interface{}, error) {
	h := Create{}

	helpBool, _ := strconv.ParseBool(c.PathParam("h"))
	if helpBool {
		return h.Help(), nil
	}

	migrationName := c.PathParam("name")
	if migrationName == "" {
		return nil, &errors.Response{Reason: "provide a name for migration"}
	}

	err := create(h, migrationName)
	if err != nil {
		return nil, err
	}

	return "Migration created: " + migrationName, nil
}

func create(f FSCreate, name string) error {
	err := createMigrationFile(f, name)
	if err != nil {
		return err
	}

	prefixes, err := getPrefixes(f)
	if err != nil {
		return err
	}

	sort.Strings(prefixes)

	err = createAllFile(f, prefixes)
	if err != nil {
		return err
	}

	return nil
}

func getPrefixes(f FSCreate) ([]string, error) {
	var prefixes []string

	files, err := f.ReadDir("./")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileParts := strings.Split(file.Name(), "_")
		if len(fileParts) < 2 || file.Name() == "000_all.go" {
			continue
		} else {
			prefixes = append(prefixes, fileParts[0])
		}
	}

	return prefixes, nil
}

// createAllFile creates the file which stores all the migrations of the project in the form of a map
func createAllFile(f FSCreate, prefixes []string) error {
	// Write 000_all.go
	fAll, err := f.Create("000_all.go")
	if err != nil {
		return err
	}

	defer fAll.Close()

	var allTemplate = template.Must(template.New("000_all").Parse(
		`// This is auto-generated file using 'gofr migrate' tool. DO NOT EDIT.
package migrations

import (
	dbmigration "developer.zopsmart.com/go/gofr/cmd/gofr/migration/dbMigration"
)

func All() map[string]dbmigration.Migrator{
	return map[string]dbmigration.Migrator{
{{range $key, $value := .}}	
		"{{ $value }}": K{{ $value }}{},{{end}}
	}
}
`))

	err = allTemplate.Execute(fAll, prefixes)
	if err != nil {
		return err
	}

	return nil
}

// createMigrationFile creates a .go file which contains the template for writing up and down migration
func createMigrationFile(f FSCreate, migrationName string) error {
	if _, err := f.Stat("migrations"); f.IsNotExist(err) {
		if er := f.Mkdir("migrations", 0777); er != nil {
			return er
		}
	}

	if err := f.Chdir("migrations"); err != nil {
		return err
	}

	currTimeStamp := time.Now().Format("20060102150405")

	migrationName = currTimeStamp + "_" + migrationName

	migrationTemplate := template.Must(template.New("migration").Parse(`package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K{{.Timestamp}} struct {
}

func (k K{{.Timestamp}}) Up(d *datastore.DataStore, logger log.Logger) error {
	return nil
}

func (k K{{.Timestamp}}) Down(d *datastore.DataStore, logger log.Logger) error {
	return nil
}
`))

	file, err := f.OpenFile(migrationName+".go", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	tData := struct {
		Timestamp         string
		MigrationFileName string
	}{currTimeStamp, migrationName}

	err = migrationTemplate.Execute(file, tData)
	if err != nil {
		return err
	}

	return nil
}
