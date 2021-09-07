package entity

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_addEntity(t *testing.T) {
	currDir, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(currDir)
	}()

	dir := t.TempDir()
	_ = os.Mkdir(dir+"/testEntity", os.ModePerm)

	var h Handler

	type args struct {
		entity     string
		entityType string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Invalid value for flag type ", args{"brand", "store"}, true},
		{"Success Case: core", args{"brand", "core"}, false},
		{"Success Case: composite", args{"brand", "composite"}, false},
		{"Success Case: consumer", args{"brand", "consumer"}, false},
	}

	for _, tt := range tests {
		_ = os.Chdir(dir + "/testEntity")

		if err := addEntity(h, tt.args.entityType, tt.args.entity); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: addEntity() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}

	{
		err := addEntity(h, "store", "brand")
		if err != nil && (err.Error() != invalidTypeError{}.Error()) {
			t.Errorf("invalid type error expected")
		}
	}
}

func TestErrors_addCore(t *testing.T) {
	currDir, _ := os.Getwd()

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	c := NewMockfileSystem(ctrl)
	dir := t.TempDir()
	path := dir + "/testEntity"
	_ = os.Mkdir(path, os.ModePerm)
	test, _ := os.Create(path + "/test.txt")
	testingFile, _ := os.Create(path + "/testingFile.txt")

	type args struct {
		name       string
		entityType string
	}

	tests := []struct {
		name        string
		args        args
		mockedCalls []*gomock.Call
		wantErr     bool
	}{
		{"error : Getwd()", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().Getwd().Return("", errors.New("test error")).Times(1),
		}, true},

		{"error : createChangeDir()", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().Getwd().Return(path, nil).AnyTimes(),
			c.EXPECT().Stat(gomock.Any()).Return(nil, errors.New("doesn't exist")).Times(1),
			c.EXPECT().IsNotExist(gomock.Any()).Return(true).Times(1),
			c.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(errors.New("test error")).Times(1),
		}, true},

		{"error composite: createChangeDir()", args{"brand", "composite"}, []*gomock.Call{
			c.EXPECT().Stat(gomock.Any()).Return(nil, errors.New("doesn't exist")).Times(1),
			c.EXPECT().IsNotExist(gomock.Any()).Return(true).Times(1),
			c.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(errors.New("test error")).Times(1),
		}, true},

		{"error consumer: createChangeDir()", args{"brand", "consumer"}, []*gomock.Call{
			c.EXPECT().Stat(gomock.Any()).Return(nil, errors.New("doesn't exist")).Times(1),
			c.EXPECT().IsNotExist(gomock.Any()).Return(true).Times(1),
			c.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(errors.New("test error")).Times(1),
		}, true},

		{"error: Chdir", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().IsNotExist(gomock.Any()).Return(false).AnyTimes(),
			c.EXPECT().Stat(gomock.Any()).Return(nil, nil).AnyTimes(),
			c.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().Chdir(path + "/core/brand").Return(errors.New("test error")).Times(1),
			c.EXPECT().Chdir(path + "/core").Return(nil).Times(1),
			c.EXPECT().OpenFile("interface.go", gomock.Any(), gomock.Any()).Return(test, nil).Times(1),
		}, true},

		{"error: OpenFile", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).Times(1),
		}, true},

		{"error: OpenFile, filePtr is nil", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},

		{"error: OpenFile, entity file open error", args{"brand", "core"}, []*gomock.Call{
			c.EXPECT().OpenFile("brand.go", gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).Times(1),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(testingFile, nil).AnyTimes(),
		}, true},
	}

	for _, tt := range tests {
		_ = os.Chdir(dir + "/testEntity")

		if err := addEntity(c, tt.args.entityType, tt.args.name); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: addEntity() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_addConsumer(t *testing.T) {
	currDir, _ := os.Getwd()
	dir := t.TempDir()
	projectDirectory := dir + "/testProject"
	_ = os.Mkdir(projectDirectory, os.ModePerm)

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	c := NewMockfileSystem(ctrl)

	type args struct {
		entity string
	}

	tests := []struct {
		name     string
		args     args
		mockCall []*gomock.Call
		wantErr  bool
	}{
		{"error http changeDir", args{"product"}, []*gomock.Call{
			c.EXPECT().Stat(gomock.Any()).Return(nil, nil).AnyTimes(),
			c.EXPECT().IsNotExist(gomock.Any()).Return(false).AnyTimes(),
			c.EXPECT().Chdir(projectDirectory + "/http").Return(errors.New("test error")).Times(1),
		}, true},

		{"error entity changeDir", args{"product"}, []*gomock.Call{
			c.EXPECT().Chdir("product").Return(errors.New("test error")).Times(1),
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
		}, true},

		{"error OpenFile", args{"product"}, []*gomock.Call{
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).Times(1),
		}, true},

		{"error OpenFile", args{"product"}, []*gomock.Call{
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},
	}

	for _, tt := range tests {
		if err := addConsumer(c, projectDirectory, tt.args.entity); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: addConsumer() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_addComposite(t *testing.T) {
	currDir, _ := os.Getwd()

	dir := t.TempDir()
	projectDirectory := dir + "/testProject"
	_ = os.Mkdir(projectDirectory, os.ModePerm)
	_ = os.Chdir(projectDirectory)
	testFile, _ := os.Create("test.go")
	compositePath := projectDirectory + "/composite"

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	c := NewMockfileSystem(ctrl)

	tests := []struct {
		name     string
		entity   string
		mockCall []*gomock.Call
		wantErr  bool
	}{
		{"error case: change dir", "brand", []*gomock.Call{
			c.EXPECT().Stat(gomock.Any()).Return(nil, nil).AnyTimes(),
			c.EXPECT().IsNotExist(gomock.Any()).Return(false).AnyTimes(),
			c.EXPECT().Chdir(gomock.Any()).Return(errors.New("test error")).Times(1),
		}, true},

		{"error case: openfile", "brand", []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).Times(2),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).Times(1),
		}, true},

		{"error case: openfile returns nil filePtr", "brand", []*gomock.Call{
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},

		{"error case: Chdir", "brand", []*gomock.Call{
			c.EXPECT().Chdir(compositePath + "/brand").Return(errors.New("test error")).Times(1),
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(testFile, nil).AnyTimes(),
		}, true},
	}
	for _, tt := range tests {
		if err := addComposite(c, projectDirectory, tt.entity); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: addComposite() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
