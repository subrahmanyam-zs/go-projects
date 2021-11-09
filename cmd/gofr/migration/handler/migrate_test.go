package handler

import (
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

func Test_getModulePath(t *testing.T) {
	_ = os.Chdir("..")

	f, err := os.Create("go.mod")
	if err != nil {
		t.Errorf("error in creating mod file")
	}

	defer f.Close()

	defer os.Remove("go.mod")

	err = os.WriteFile("go.mod", []byte("module example.com/my-project\n\ngo 1.17\n"), os.ModeDevice)
	if err != nil {
		t.Errorf("error in writing to mod file")
	}

	ctrl := gomock.NewController(t)
	fs := NewMockFSMigrate(ctrl)

	fs.EXPECT().OpenFile("../go.mod", os.O_RDONLY, gomock.Any()).Return(f, nil)

	name, err := getModulePath(fs, "random-dir")

	assert.Nil(t, err)

	assert.Equal(t, "example.com/my-project", name)
}

func Test_createMain(t *testing.T) {
	currDir, _ := os.Getwd()

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	mockFS := NewMockFSMigrate(ctrl)
	dir := t.TempDir()
	_ = os.Chdir(dir)
	f, _ := os.Create("test.txt")
	f2, _ := os.Create("main.go")
	modFile, _ := os.Create("go.mod")
	_, _ = modFile.WriteString("module moduleName")

	defer modFile.Close()

	type args struct {
		method    string
		db        string
		directory string
	}

	tests := []struct {
		name      string
		args      args
		mockCalls []*gomock.Call
		wantErr   bool
	}{
		{"database not supported", args{"UP", "kafka", dir}, []*gomock.Call{}, true},
		{"Project Not in GOPATH error", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(f, &errors.Response{Reason: "test error"}).Times(1),
		}, true},
		{"success", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(modFile, nil).Times(1),
			mockFS.EXPECT().Stat("build").Return(nil, &errors.Response{Reason: "test error"}),
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(true),
			mockFS.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(nil).Times(1),
			mockFS.EXPECT().Chdir(gomock.Any()).Return(nil),
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(f2, nil).Times(1),
		}, false},
		{"Success case GOPATH", args{"DOWN", "gorm", "/home/appsadm/go/src/gofr"}, []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(f, &errors.Response{Reason: "test error"}).Times(1),
			mockFS.EXPECT().Stat("build").Return(nil, &errors.Response{Reason: "test error"}),
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(false),
			mockFS.EXPECT().Chdir(gomock.Any()).Return(nil).Times(1),
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(f2, nil).Times(1),
		}, false},
		{"mkdir error", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(modFile, nil).Times(1),
			mockFS.EXPECT().Stat("build").Return(nil, &errors.Response{Reason: "test error"}),
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(true),
			mockFS.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(&errors.Response{Reason: "test error"}).Times(1),
		}, true},
		{"chdir error", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(modFile, nil).Times(1),
			mockFS.EXPECT().Stat("build").Return(nil, nil),
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(false),
			mockFS.EXPECT().Chdir(gomock.Any()).Return(&errors.Response{Reason: "test error"}).Times(1),
		}, true},
		{"openFile error", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(modFile, nil).Times(1),
			mockFS.EXPECT().Stat("build").Return(nil, nil),
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(false),
			mockFS.EXPECT().Chdir(gomock.Any()).Return(nil).Times(1),
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, &errors.Response{Reason: "test error"}).Times(1),
		}, true},
		{"template execution error", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(modFile, nil).Times(1),
			mockFS.EXPECT().Stat("build").Return(nil, nil),
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(false),
			mockFS.EXPECT().Chdir(gomock.Any()).Return(nil).Times(1),
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},
	}

	for _, tt := range tests {
		tt := tt
		if err := createMain(mockFS, tt.args.method, tt.args.db, tt.args.directory, nil); (err != nil) != tt.wantErr {
			t.Errorf("TestCase %v: createMain() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_runMigration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFS := NewMockFSMigrate(ctrl)

	type args struct {
		method string
		db     string
	}

	tests := []struct {
		name      string
		args      args
		mockCalls []*gomock.Call
		want      interface{}
		wantErr   bool
	}{
		{"Getwd() error", args{}, []*gomock.Call{
			mockFS.EXPECT().Getwd().Return("", &errors.Response{Reason: "test error"}).Times(1),
		}, nil, true},

		{"Chdir and  dir not exists error", args{}, []*gomock.Call{
			mockFS.EXPECT().Getwd().Return("", nil).AnyTimes(),
			mockFS.EXPECT().Chdir("migrations").Return(nil).AnyTimes(),
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(true).Times(1),
		}, nil, true},

		{"createMain error", args{}, []*gomock.Call{
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(false).AnyTimes(),
			mockFS.EXPECT().Stat(gomock.Any()).Return(nil, nil).AnyTimes(),
			mockFS.EXPECT().Chdir("build").Return(&errors.Response{Reason: "test error"}).AnyTimes(),
		}, nil, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := runMigration(mockFS, tt.args.method, tt.args.db, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("runMigration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runMigration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
