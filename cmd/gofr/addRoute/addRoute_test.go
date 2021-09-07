package addroute

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_addRoute(t *testing.T) {
	currDir, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(currDir)
	}()

	dir := t.TempDir()
	_ = os.Mkdir(dir+"/testEntity", os.ModePerm)
	_ = os.Chdir(dir + "/testEntity")
	_, _ = os.Create("main.go")

	var h Handler

	type args struct {
		methods string
		path    string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"success : no methods", args{"", "/hello-world"}, nil},
		{"success : one method", args{"GET", "/hello"}, nil},
		{"success : multiple methods", args{"GET,POST", "/test"}, nil},
	}

	for _, tt := range tests {
		_ = os.Chdir(dir + "/testEntity")

		err := addRoute(h, tt.args.methods, tt.args.path)
		if err != nil && (err.Error() != tt.wantErr.Error()) {
			t.Errorf("Test %v: addRoute() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestErrors(t *testing.T) {
	currDir, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(currDir)
	}()

	dir := t.TempDir()
	_ = os.Chdir(dir)
	_, _ = os.Create("main.go")

	var h Handler

	type args struct {
		path   string
		method string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"invalid path", args{"/$/{id}", http.MethodGet}, invalidPathError{"$/{id}"}},
		{"invalid method", args{"/abcd/{id}", http.MethodPatch}, invalidMethodError{"PATCH"}},
	}

	for _, tt := range tests {
		_ = os.Chdir(dir)

		err := addRoute(h, tt.args.method, tt.args.path)
		if err != nil && err.Error() != tt.wantErr.Error() {
			t.Errorf("Test %v: addRoute() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestErrors_FileSystem(t *testing.T) {
	currDir, _ := os.Getwd()
	dir := t.TempDir()

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	c := NewMockfileSystem(ctrl)
	_ = os.Chdir(dir)
	test, _ := os.Create("test.go")

	type args struct {
		path   string
		method string
	}

	tests := []struct {
		name      string
		args      args
		mockCalls []*gomock.Call
		wantErr   bool
	}{
		{"error: Match error", args{"/brand", http.MethodGet}, []*gomock.Call{
			c.EXPECT().Match(gomock.Any(), gomock.Any()).Return(false, errors.New("test error")).Times(1),
		}, true},

		{"error: OpenFile", args{"/brand", http.MethodGet}, []*gomock.Call{
			c.EXPECT().Match(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).Times(1),
		}, true},

		{"error: Getwd", args{"/brand", http.MethodGet}, []*gomock.Call{
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(test, nil).AnyTimes(),
			c.EXPECT().Getwd().Return("", errors.New("test error")).Times(1),
		}, true},
	}

	for _, tt := range tests {
		if err := addRoute(c, tt.args.method, tt.args.path); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: addRoute() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_populateMain(t *testing.T) {
	currDir, _ := os.Getwd()
	dir := t.TempDir()

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	c := NewMockfileSystem(ctrl)
	_ = os.Chdir(dir)
	testFile, _ := os.OpenFile("testing.go", os.O_CREATE|os.O_RDONLY, 0666)

	type args struct {
		mainString string
	}

	tests := []struct {
		name      string
		args      args
		mockCalls []*gomock.Call
		wantErr   bool
	}{
		{"Error chdir", args{"package main"}, []*gomock.Call{
			c.EXPECT().Getwd().Return(dir+"/testEntity", nil).AnyTimes(),
			c.EXPECT().Chdir(gomock.Any()).Return(errors.New("test error")).Times(1),
		}, true},

		{"Error OpenFile", args{"package main"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).Times(1),
		}, true},

		{"Error OpenFile", args{"package main"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},

		{"Error OpenFile", args{"package main"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(testFile, nil).Times(1),
		}, true},
	}

	for _, tt := range tests {
		if err := populateMain(c, tt.args.mainString, ""); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: populateMain() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_populateHandler(t *testing.T) {
	currDir, _ := os.Getwd()
	dir := t.TempDir()

	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()

		_ = os.Chdir(currDir)
	}()

	c := NewMockfileSystem(ctrl)
	_ = os.Chdir(dir)
	testFile, _ := os.OpenFile("testing.go", os.O_CREATE|os.O_RDONLY, 0666)

	type args struct {
		path          string
		handlerString string
	}

	tests := []struct {
		name     string
		args     args
		mockCall []*gomock.Call
		wantErr  bool
	}{
		{"chdir error", args{"brand", "package brand"}, []*gomock.Call{
			c.EXPECT().Stat(gomock.Any()).Return(nil, nil).AnyTimes(),
			c.EXPECT().IsNotExist(gomock.Any()).Return(false).AnyTimes(),
			c.EXPECT().Chdir("http").Return(errors.New("test error")).Times(1),
		}, true},

		{"chdir error", args{"brand", "package brand"}, []*gomock.Call{
			c.EXPECT().Chdir("http").Return(nil).Times(1),
			c.EXPECT().Chdir("brand").Return(errors.New("test error")).Times(1),
		}, true},

		{"openfile error", args{"brand", "package brand"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test error")).Times(1),
		}, true},

		{"openfile error: nil returned", args{"brand", "package brand"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1),
		}, true},

		{"openfile error: read only file given to write the content", args{"brand", "package brand"}, []*gomock.Call{
			c.EXPECT().Chdir(gomock.Any()).Return(nil).AnyTimes(),
			c.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(testFile, nil).Times(1),
		}, true},
	}
	for _, tt := range tests {
		if err := populateHandler(c, tt.args.path, tt.args.handlerString); (err != nil) != tt.wantErr {
			t.Errorf("Test %v: populateHandler() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_createChangeDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockfileSystem(ctrl)
	c.EXPECT().Stat(gomock.Any()).Return(nil, errors.New("test error")).Times(1)
	c.EXPECT().IsNotExist(gomock.Any()).Return(true).Times(1)
	c.EXPECT().Mkdir(gomock.Any(), gomock.Any()).Return(errors.New("test error")).Times(1)

	tests := []struct {
		name      string
		directory string
		wantErr   bool
	}{
		{"mkdir error", "test", true},
	}

	for _, tt := range tests {
		if err := createChangeDir(c, tt.directory); (err != nil) != tt.wantErr {
			t.Errorf("Test %v : createChangeDir() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_existCheck(t *testing.T) {
	type args struct {
		file io.ReadSeeker
		elem string
	}

	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{"error case", args{nil, ""}, 0, false},
	}

	for _, tt := range tests {
		got, got1 := existCheck(tt.args.file, tt.args.elem)

		if got != tt.want {
			t.Errorf("existCheck() got = %v, want %v", got, tt.want)
		}

		if got1 != tt.want1 {
			t.Errorf("existCheck() got1 = %v, want %v", got1, tt.want1)
		}
	}
}

func Test_importSortCheck(t *testing.T) {
	lineString := importSortCheck("abc", "developer.zopsmart.com/go/gofr/pkg/gofr")
	if strings.Contains(lineString, `"abc"
		developer.zopsmart.com/go/gofr/pkg/gofr`) {
		t.Errorf("import sort failed. Got: %v", lineString)
	}
}
