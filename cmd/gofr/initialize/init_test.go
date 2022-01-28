package initialize

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_createProjectErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fs := NewMockFileSystem(ctrl)

	type args struct {
		f    fileSystem
		name string
	}

	tests := []struct {
		name      string
		mockCalls []*gomock.Call
		args      args
		wantErr   bool
	}{
		{"Error Mkdir", []*gomock.Call{
			fs.EXPECT().Mkdir("testProject", gomock.Any()).Return(errors.New("test error")).Times(1),
		}, args{fs, "testProject"}, true},

		{"Error Chdir", []*gomock.Call{
			fs.EXPECT().Mkdir("testProject", gomock.Any()).Return(nil).Times(1),
			fs.EXPECT().Chdir(gomock.Any()).Return(errors.New("test error")).Times(1),
		}, args{fs, "testProject"}, true},

		{"Error Mkdir - Standard Directories", []*gomock.Call{
			fs.EXPECT().Chdir("testProject").Return(nil).AnyTimes(),
			fs.EXPECT().Mkdir("configs", gomock.Any()).Return(errors.New("test error")).Times(1),
			fs.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(nil).AnyTimes(),
		}, args{fs, "testProject"}, true},

		{"Error Create", []*gomock.Call{
			fs.EXPECT().Create(gomock.Any()).Return(nil, errors.New("test error")).Times(1),
		}, args{fs, "testProject"}, true},

		{"Error WriteString", []*gomock.Call{
			fs.EXPECT().Create(gomock.Any()).Return(nil, nil).Times(1),
		}, args{fs, "testProject"}, true},
	}

	for _, tt := range tests {
		if err := createProject(tt.args.f, tt.args.name); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_createProject(t *testing.T) {
	var h Handler

	currDir, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(currDir)
	}()

	type args struct {
		f    fileSystem
		name string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success Case", args{h, "testProject"}, false},
		{"Project with same name already exists error", args{h, "testProject"}, true},
	}

	dir := t.TempDir()

	for _, tt := range tests {
		_ = os.Chdir(dir)

		if err := createProject(h, tt.args.name); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: createProject() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
