package handler

import (
	"os"
	"testing"

	"developer.zopsmart.com/go/gofr/cmd/gofr/migration"
	"developer.zopsmart.com/go/gofr/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_create(t *testing.T) {
	var (
		rwxMode = os.FileMode(migration.RWXMode)
		rwMode  = os.FileMode(migration.RWMode)
	)

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()
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
			mockFS.EXPECT().Stat("migrations").Return(nil, &errors.Response{Reason: "test error"}).MaxTimes(6),
			mockFS.EXPECT().IsNotExist(&errors.Response{Reason: "test error"}).Return(true).MaxTimes(6),
			mockFS.EXPECT().Mkdir("migrations", rwxMode).Return(&errors.Response{Reason: "test error"}).Times(1),
		}, true},

		{"chdir error", "testing", []*gomock.Call{
			mockFS.EXPECT().Mkdir("migrations", rwxMode).Return(nil).MaxTimes(5),
			mockFS.EXPECT().Chdir("migrations").Return(&errors.Response{Reason: "test error"}).Times(1),
		}, true},

		{"unable to openfile", "testing", []*gomock.Call{
			mockFS.EXPECT().Chdir("migrations").Return(nil).MaxTimes(4),
			mockFS.EXPECT().OpenFile(gomock.Any(), os.O_CREATE|os.O_WRONLY, rwMode).Return(nil, &errors.Response{Reason: "test error"}).Times(1),
		}, true},

		{"template execution error", "testing", []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), os.O_CREATE|os.O_WRONLY, rwMode).Return(nil, nil).Times(1),
		}, true},

		{"read dir error", "testing", []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), os.O_CREATE|os.O_WRONLY, rwMode).Return(file, nil).Times(1),
			mockFS.EXPECT().ReadDir(gomock.Any()).Return(nil, &errors.Response{Reason: "test error"}).Times(1),
		}, true},

		{"create file error", "testing", []*gomock.Call{
			mockFS.EXPECT().OpenFile(gomock.Any(), os.O_CREATE|os.O_WRONLY, rwMode).Return(file1, nil),
			mockFS.EXPECT().ReadDir("./").Return(allFiles, nil),
			mockFS.EXPECT().Create("000_all.go").Return(nil, &errors.Response{Reason: "test error"}).Times(1),
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

func Test_getPrefixes(t *testing.T) {
	dir := t.TempDir()

	_ = os.Chdir(dir)
	// these files will be ignored
	_, _ = os.Create("20190320095356_test.go")                       // files with len < 2 will be ignored
	_, _ = os.Create("000_all.go")                                   // 000_all.go will be ignored.
	_, _ = os.Create("20220320095352_table_employee_create_test.go") // ignores the files that have the suffix test

	// files whose prefixes will be added to the slice prefixes
	_, _ = os.Create("20220410095352_table_employee_create.go")
	_, _ = os.Create("20210520095352_table_employee_create.go")
	_, _ = os.Create("20190320095352_table_employee_create.go")

	allFiles, _ := os.ReadDir(dir)
	ctrl := gomock.NewController(t)
	mockFS := NewMockFSCreate(ctrl)

	tests := []struct {
		desc      string
		output    []string
		err       error
		mockEntry []os.DirEntry
	}{
		{"Success", []string{"20190320095352", "20210520095352", "20220410095352"}, nil, allFiles},
		{"Error in Reading files", nil, errors.Error("Error while reading file"), nil},
	}
	for i, tc := range tests {
		mockFS.EXPECT().ReadDir("./").Return(tc.mockEntry, tc.err)

		result, err := getPrefixes(mockFS)

		assert.Equal(t, tc.output, result, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}
