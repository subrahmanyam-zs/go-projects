package entity

import (
	"errors"
	"os"
	"testing"

	"developer.zopsmart.com/go/gofr/cmd/gofr/migration"

	"github.com/golang/mock/gomock"
)

func Test_populateEntityFile(t *testing.T) {
	currDir, _ := os.Getwd()

	dir := t.TempDir()
	projectDirectory := dir + "/testProject"

	_ = os.Mkdir(projectDirectory, os.ModePerm)
	_ = os.Chdir(projectDirectory)

	testFile, _ := os.OpenFile("test.go", os.O_CREATE|os.O_RDONLY, migration.RWMode)
	mainFile, _ := os.OpenFile("main.go", os.O_RDONLY, migration.RWMode)

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	c := NewMockfileSystem(ctrl)

	type args struct {
		entity string
		types  string
	}

	tests := []struct {
		name      string
		args      args
		mockCalls []*gomock.Call
		wantErr   bool
	}{
		{"error: Chdir", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(errors.New("test error")).Times(1),
		}, true},

		{"error: OpenFile", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).Times(1),
		}, true},

		{"error: OpenFile returns nil", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},

		{"error: OpenFile returns nil", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(testFile, nil).Times(1),
		}, true},

		{"error: OpenFile returns nil", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(mainFile, nil).Times(1),
		}, true},
	}

	for _, tt := range tests {
		if err := populateEntityFile(c, dir, projectDirectory, tt.args.entity, tt.args.types); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: populateEntityFile() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_createModel(t *testing.T) {
	currDir, _ := os.Getwd()

	dir := t.TempDir()
	projectDirectory := dir + "/testProject"
	_ = os.Mkdir(projectDirectory, os.ModePerm)
	_ = os.Chdir(projectDirectory)
	testFile, _ := os.OpenFile("testRead.go", os.O_CREATE|os.O_RDONLY, migration.RWMode)

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	c := NewMockfileSystem(ctrl)

	tests := []struct {
		name      string
		entity    string
		mockCalls []*gomock.Call
		wantErr   bool
	}{
		{"error case: chdir", "brand", []*gomock.Call{
			c.EXPECT().Stat(gomock.Any()).Return(nil, nil).AnyTimes(),
			c.EXPECT().IsNotExist(gomock.Any()).Return(false).AnyTimes(),
			c.EXPECT().Chdir(gomock.Any()).Return(errors.New("test error")).Times(1),
		}, true},

		{"error case: openfile", "brand", []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).Times(1),
		}, true},

		{"error case: openfile returns nil", "brand", []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},

		{"error case: openfile returns nil", "brand", []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(testFile, nil).Times(1),
		}, true},
	}
	for _, tt := range tests {
		if err := createModel(c, projectDirectory, tt.entity); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: createModel() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_populateInterfaceFiles(t *testing.T) {
	currDir, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(currDir)
	}()

	dir := t.TempDir()
	projectDirectory := dir + "/testProject"
	_ = os.Mkdir(projectDirectory, os.ModePerm)
	_ = os.Chdir(projectDirectory)

	testFile, _ := os.OpenFile("test.go", os.O_CREATE|os.O_RDONLY, migration.RWMode)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"error file is readOnly", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := populateInterfaceFiles("test", dir, "core", testFile); (err != nil) != tt.wantErr {
				t.Errorf("populateInterfaceFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
