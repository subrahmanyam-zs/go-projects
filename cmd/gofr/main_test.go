package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"developer.zopsmart.com/go/gofr/cmd/gofr/migration"
	"developer.zopsmart.com/go/gofr/pkg/gofr/assert"
)

func TestCLI(t *testing.T) {
	currDir, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(currDir)
	}()

	dir := t.TempDir()
	_ = os.Chdir(dir)

	flag.String("name", "", "")
	flag.String("methods", "", "")
	flag.String("path", "", "")
	flag.String("type", "", "")

	assert.CMDOutputContains(t, main, "gofr init -name=testGoProject", "Successfully created project: testGoProject")

	_ = os.Chdir(dir + "/testGoProject")

	assert.CMDOutputContains(t, main, "gofr add -methods=all -path=/foo", "Added route: /foo")

	_ = os.Chdir(dir + "/testGoProject")

	assert.CMDOutputContains(t, main, "gofr add -methods=all -path=/foo", "route foo is already present")

	_ = os.Chdir(dir + "/testGoProject")

	assert.CMDOutputContains(t, main, "gofr entity -type=core -name=person", "Successfully created entity: person")
}

func Test_Migrate(t *testing.T) {
	assert.CMDOutputContains(t, main, "gofr migrate -method=ABOVE -database=gorm", "invalid flag: method")
	assert.CMDOutputContains(t, main, "gofr migrate -method=UP -database=gorm", "migrations do not exist")
	assert.CMDOutputContains(t, main, "gofr migrate -method=UP -database=", "invalid flag: database")

	currDir := t.TempDir()
	_ = os.Chdir(currDir)

	path, _ := os.MkdirTemp(currDir, "migrateCreateTest")
	_ = os.Chdir(path)

	assert.CMDOutputContains(t, main, "gofr migrate create", "provide a name for migration")
	assert.CMDOutputContains(t, main, "gofr migrate create -name=testMigration", "Migration created")

	_ = os.Chdir(currDir + "/migrateCreateTest")

	assert.CMDOutputContains(t, main, "gofr migrate create -name=migrationTest", "Migration created")

	_ = os.Chdir(currDir + "/migrateCreateTest")

	assert.CMDOutputContains(t, main, "gofr migrate -method=UP -database=gorm", "migrations do not exists")

	_ = os.Chdir(currDir + "/migrateCreateTest")

	assert.CMDOutputContains(t, main, "gofr migrate -method=DOWN -database=gorm", "migrations do not exists")

	_ = os.Chdir(currDir + "/migrateCreateTest")

	assert.CMDOutputContains(t, main, "gofr migrate -method=DOWN -database=mongo", "migrations do not exists")

	_ = os.Chdir(currDir + "/migrateCreateTest")

	assert.CMDOutputContains(t, main, "gofr migrate -method=DOWN -database=cassandra", "migrations do not exists")

	_ = os.Chdir(currDir + "/migrateCreateTest")

	assert.CMDOutputContains(t, main, "gofr migrate -method=DOWN -database=ycql", "migrations do not exists")

	_ = os.Chdir(currDir + "/migrateCreateTest")

	assert.CMDOutputContains(t, main, "gofr migrate -method=DOWN -database=redis -tag=20200123143215", "migrations do not exists")
}

func Test_CreateMigration(t *testing.T) {
	currDir, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(currDir)
	}()

	_ = os.Chdir(t.TempDir())
	_ = os.Mkdir("migrationTest", migration.RWXMode)
	_ = os.Chdir("migrationTest")

	assert.CMDOutputContains(t, main, "gofr migrate create -name=removeColumn", "Migration created: removeColumn")
}

func Test_Integration(t *testing.T) {
	assert.CMDOutputContains(t, main, "gofr help", "Available Commands")
}

func Test_HelpGenerate(t *testing.T) {
	assert.CMDOutputContains(t, main, "gofr init -h", "creates a project structure inside the directory specified in the name flag")
	assert.CMDOutputContains(t, main, "gofr entity -h", "creates a template and interface for an entity")
	assert.CMDOutputContains(t, main, "gofr add -h", "add routes and creates a handler template")
	assert.CMDOutputContains(t, main, "gofr migrate -h", "usage: gofr migrate")
	assert.CMDOutputContains(t, main, "gofr migrate create -h", "usage: gofr migrate create")
}

// nolint:funlen // reducing the function length reduces readability
func Test_test_Success(t *testing.T) {
	const ymlStr = `openapi: 3.0.1
info:
  title: LogisticsAPI
  version: '0.1'
servers:
  - url: 'http://api.staging.zopsmart.com'
paths:
  /hello-world:
    get:
      tags:
        - Hello
      description: Sample API Hello
      responses:
        '200':
          description: Sample API Hello
  /hello:
    get:
      tags:
        - Hello
      description: Sample API Hello with name
      parameters:
        - name: x-zopsmart-tenent
          in: header
          schema:
            type: string
            format: uuid
          example: 'good4more'
        - name: x-correlation-ID
          in: header
          schema:
            type: string
            format: uuid
          example: 
        - name: custom-header
          in: header
          schema:
            type: string
            format: uuid
          example: 'abc,xyz,ijk'
        - name: name
          in: query
          schema:
            type: string
          example: 'Roy'
        - name: age
          in: body
          schema:
            type: float
          example: 32189.5
        - name: hasAcc
          in: query
          schema:
            type: bool
          example: true
        - name: nick_names
          in: query
          schema:
            type: array
          example: [abc, def, ghi]
      responses:
        '200':
          description: Sample API Hello with name
    post:
      tags:
        - Hello
      description: Sample API Hello with name
      parameters:
        - name: x-zopsmart-tenent
          in: header
          schema:
            type: string
            format: uuid
          example: 'good4more'
        - name: x-correlation-ID
          in: header
          schema:
            type: string
            format: uuid
          example: 
        - name: custom-header
          in: header
          schema:
            type: string
            format: uuid
          example: 'abc,xyz,ijk'
        - name: id
          in: path
          schema:
            type: int
          example: 5
        - name: catalog_item
          in: body
          schema:
            type: object
            properties:
              id:
                type: integer
              name:
                type: string
            required:
              - id
              - name
          example:
            id: 38
            name: T-shirt
            salary: 452.05
      responses:
        '200':
          description: Sample API Hello with name`

	currDir, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(currDir)
	}()

	d1 := []byte(ymlStr)

	tempFile, err := os.CreateTemp(t.TempDir(), "dat1.yml")
	if err != nil {
		t.Error(err)
	}

	_, err = tempFile.Write(d1)
	if err != nil {
		t.Error(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode("{}")
	}))
	hostPort := strings.Replace(server.URL, "http://", "", 1)
	assert.CMDOutputContains(t, main, "gofr test -host="+hostPort+" -source="+tempFile.Name(), "Test Passed!")
}

func Test_test_Error(t *testing.T) {
	const ymlStr = `openapi: 3.0.1
info:
  title: LogisticsAPI
  version: '0.1'
servers:
  - url: 'http://api.staging.zopsmart.com'
paths:
  /hello/{id}:
    put:
      tags:
        - Hello
      description: Sample API Hello with name
      parameters:
        - name: id
          in: path
          schema:
            type: int
          example: 5
        - name: catalog_item
          in: body
          schema:
            type: object
            properties:
              id:
                type: integer
              name:
                type: string
            required:
              - id
              - name
          example:
            id: 38
            name: T-shirt
            salary: 452.05
      responses:
        '403':
          description: Sample API Hello with name
    delete:
      tags:
        - Post Hello
      description: Sample API Hello with name
      parameters:
        - name: id
          in: path
          schema:
            type: int
          example: 5
      responses:
        '400':
          description: Sample API Hello`

	currDir, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(currDir)
	}()

	d1 := []byte(ymlStr)

	tempFile, err := os.CreateTemp(t.TempDir(), "dat1.yml")
	if err != nil {
		t.Error(err)
	}

	_, err = tempFile.Write(d1)
	if err != nil {
		t.Error(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode("{}")
	}))
	hostPort := strings.Replace(server.URL, "http://", "", 1)
	assert.CMDOutputContains(t, main, "gofr test -host="+hostPort+" -source="+tempFile.Name(), "failed")

	// case to check test help
	assert.CMDOutputContains(t, main, "gofr test -h", "runs integration test for a given configuration")

	// case when source not specified
	assert.CMDOutputContains(t, main, "gofr test -host="+hostPort, "source not specified")

	// case when host not specified
	assert.CMDOutputContains(t, main, "gofr test -source="+tempFile.Name(), "please provide host")

	// case when source is incorrect
	assert.CMDOutputContains(t, main, "gofr test -host="+hostPort+" -source=/some/fake/path/data.yml", "no such file or directory")
}
