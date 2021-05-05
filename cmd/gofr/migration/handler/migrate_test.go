package handler

import (
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/zopsmart/gofr/pkg/errors"
)

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
		{"mkdir error", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().Stat("build").Return(nil, &errors.Response{Reason: "test error"}).AnyTimes(),
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(true).AnyTimes(),
			mockFS.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(&errors.Response{Reason: "test error"}).Times(1),
		}, true},
		{"chdir error", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(nil).AnyTimes(),
			mockFS.EXPECT().Chdir(gomock.Any()).Return(&errors.Response{Reason: "test error"}).Times(1),
		}, true},
		{"openFile error", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, &errors.Response{Reason: "test error"}).Times(1),
		}, true},
		{"template execution error", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},
		{"success", args{"DOWN", "gorm", dir}, []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(f, nil).Times(1),
		}, false},
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
