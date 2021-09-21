package handler

import (
	"os"
	"testing"

	"developer.zopsmart.com/go/gofr/cmd/gofr/migration"
	"developer.zopsmart.com/go/gofr/pkg/errors"

	"github.com/golang/mock/gomock"
)

func Test_create(t *testing.T) {
	currDir, _ := os.Getwd()

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	dir := t.TempDir()
	_ = os.Chdir(dir)
	allFiles, _ := os.ReadDir(dir)
	file, _ := os.OpenFile("test.txt", os.O_CREATE|os.O_WRONLY, migration.RWMode)
	file1, _ := os.OpenFile("test1.txt", os.O_CREATE|os.O_WRONLY, migration.RWMode)

	mockFS := NewMockFSCreate(ctrl)

	tests := []struct {
		name      string
		fileName  string
		mockCalls []*gomock.Call
		wantErr   bool
	}{
		{"mkdir error", "testing", []*gomock.Call{
			mockFS.EXPECT().Stat("migrations").Return(nil, &errors.Response{Reason: "test error"}).AnyTimes(),
			mockFS.EXPECT().IsNotExist(gomock.Any()).Return(true).AnyTimes(),
			mockFS.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(&errors.Response{Reason: "test error"}).Times(1),
		}, true},

		{"chdir error", "testing", []*gomock.Call{
			mockFS.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(nil).AnyTimes(),
			mockFS.EXPECT().Chdir(gomock.Any()).Return(&errors.Response{Reason: "test error"}).Times(1),
		}, true},

		{"unable to openfile", "testing", []*gomock.Call{
			mockFS.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, &errors.Response{Reason: "test error"}).Times(1),
		}, true},

		{"template execution error", "testing", []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},

		{"read dir error", "testing", []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(file, nil).Times(1),
			mockFS.EXPECT().ReadDir(gomock.Any()).Return(nil, &errors.Response{Reason: "test error"}).Times(1),
		}, true},

		{"create file error", "testing", []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(file1, nil).AnyTimes(),
			mockFS.EXPECT().ReadDir(gomock.Any()).Return(allFiles, nil).AnyTimes(),
			mockFS.EXPECT().Create(gomock.Any()).Return(nil, &errors.Response{Reason: "test error"}).Times(1),
		}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := create(mockFS, tt.fileName); (err != nil) != tt.wantErr {
				t.Errorf("create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
