package handler

import (
	"bufio"
	"bytes"
	"go/build"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"developer.zopsmart.com/go/gofr/cmd/gofr/helper"
	mg "developer.zopsmart.com/go/gofr/cmd/gofr/migration"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

const (
	DOWN = "DOWN"
	UP   = "UP"
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

func (h Handler) Help() interface{} {
	return helper.Generate(helper.Help{
		Example: `gofr migrate -method=UP -database=gorm`,
		Flag: `method: UP or DOWN
database: gorm  // gorm supports following dialects: mysql, mssql, postgres, sqlite`,
		Usage:       "gofr migrate -method=<method> -database=<database>",
		Description: "runs the migration for method UP or DOWN as provided and for the given database",
	})
}

func Migrate(c *gofr.Context) (interface{}, error) {
	h := Handler{}

	helpBool, _ := strconv.ParseBool(c.PathParam("h"))
	if helpBool {
		return h.Help(), nil
	}

	db := strings.ToUpper(c.PathParam("database"))
	if db == "" {
		return nil, &errors.Response{Reason: "invalid flag: database"}
	}

	method := strings.ToUpper(c.PathParam("method"))
	if method != UP && method != DOWN {
		return nil, &errors.Response{Reason: "invalid flag: method"}
	}

	var tagSlc []string

	if method == DOWN {
		tag := c.Param("tag")
		if tag != "" {
			tagSlc = strings.Split(tag, ",")
		}
	}

	return runMigration(h, method, db, tagSlc)
}

func runMigration(f FSMigrate, method, db string, tagSlc []string) (interface{}, error) {
	dir, err := f.Getwd()
	if err != nil {
		return nil, err
	}

	err = f.Chdir("migrations")
	if f.IsNotExist(err) {
		return nil, &errors.Response{Reason: "migrations do not exists! If you have created migrations " +
			"please run the command from the project's root directory"}
	}

	err = createMain(f, method, db, dir, tagSlc)
	if err != nil {
		return nil, err
	}

	output, err := exec.Command("go", "run", "main.go").Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func createMain(f FSMigrate, method, db, directory string, tagSlc []string) error {
	dbStr := ""

	switch strings.ToLower(db) {
	case "gorm":
		dbStr += "db := dbmigration.NewGorm(k.GORM())"
	case "mongo":
		dbStr += "db := dbmigration.NewMongo(k.MongoDB)"
	case "cassandra":
		dbStr += "db := dbmigration.NewCassandra(&k.Cassandra)"
	case "redis":
		dbStr += "db := dbmigration.NewRedis(k.Redis)"
	case "ycql":
		dbStr += "db := dbmigration.NewYCQL(&k.YCQL)"
	default:
		return &errors.Response{Reason: "database not supported"}
	}

	lastIndex := strings.LastIndex(directory, "/")
	projectName := directory[lastIndex+1:]

	moduleName, err := getModulePath(f, directory)
	if err != nil {
		return err
	}

	err = templateCreate(f, projectName, method, dbStr, moduleName, tagSlc)
	if err != nil {
		return err
	}

	return nil
}

func templateCreate(f FSMigrate, projectName, method, dbStr, moduleName string, tagSlc []string) error {
	migration := `migrations.All()` // if method is UP or DOWN method with no specific migrations to run, then `migrations.All() is used
	mainTemplate := template.Must(template.New("").Parse(`// This is auto-generated file using 'gofr migrate' tool. DO NOT EDIT.
package main

import (
	"{{.ModuleName}}/migrations"
	"developer.zopsmart.com/go/gofr/cmd/gofr/migration"
	"developer.zopsmart.com/go/gofr/cmd/gofr/migration/dbMigration"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()
	{{.Database}}	

	err := migration.Migrate("{{.ProjectName}}", db, {{.Migration}}, "{{.Method}}", k.Logger)
	if err != nil {
		k.Logger.Error(err)
	}
}
`))

	if method == DOWN && len(tagSlc) != 0 {
		migration = getDownString(tagSlc)
	}

	tData := struct {
		ProjectName string
		Method      string
		Database    string
		Migration   string
		ModuleName  string
	}{projectName, method, dbStr, migration, moduleName}

	if _, err := f.Stat("build"); f.IsNotExist(err) {
		if er := f.Mkdir("build", os.ModePerm); er != nil {
			return er
		}
	}

	if err := f.Chdir("build"); err != nil {
		return err
	}

	os.RemoveAll("main.go")

	mainFile, err := f.OpenFile("main.go", os.O_CREATE|os.O_WRONLY, mg.RWMode)
	if err != nil {
		return err
	}

	err = mainTemplate.Execute(mainFile, tData)
	if err != nil {
		return err
	}

	return nil
}

// getDownString adds all the migrations in the `tagSlc` in the required format
func getDownString(tagSlc []string) string {
	var allTemplate = template.Must(template.New("000_all").Parse(
		`map[string]dbmigration.Migrator{
{{range $key, $value := .}}"{{ $value }}": migrations.K{{ $value }}{},{{end}}
	}`))

	var buf bytes.Buffer

	err := allTemplate.Execute(&buf, tagSlc)
	if err != nil {
		return ""
	}

	return buf.String()
}

// Function to get modulePath that to be used in the import for Migration
func getModulePath(f FSMigrate, directory string) (string, error) {
	var modulePath string

	file, err := f.OpenFile("../go.mod", os.O_RDONLY, mg.RWMode)
	if err != nil {
		return checkGoPath(directory)
	}

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		modulePath = strings.Split(scanner.Text(), " ")[1]
	}

	defer file.Close()

	return modulePath, nil
}

func checkGoPath(directory string) (string, error) {
	var modulePath string

	goPath := build.Default.GOPATH

	if strings.Contains(directory, goPath+"/src/") {
		r := strings.SplitAfter(directory, goPath+"/src")
		if len(r) > 1 {
			modulePath = r[1]
		}
	} else {
		return "", &errors.Response{Reason: "Project is not in GOPATH and go.mod file not found in current directory"}
	}

	return modulePath, nil
}
