package file

import (
	"errors"
	"io"
	"os"
	"testing"

	pkgSftp "github.com/pkg/sftp"
	"github.com/stretchr/testify/assert"

	pkgErr "developer.zopsmart.com/go/gofr/pkg/errors"
)

func Test_NewSFTPFile(t *testing.T) {
	filename := "test.txt"
	mode := READWRITE
	c1 := &SFTPConfig{Host: "localhost", User: "", Password: "", Port: 22}
	expErr := errors.New("")
	_, err := newSFTPFile(c1, filename, mode)
	assert.IsTypef(t, expErr, err, "Test failed, Expected:%v, got:%v ", expErr, err)
}

type mockSftpClient struct {
	t *testing.T
	*pkgSftp.Client
}

func (s mockSftpClient) Open(fileName string) (io.ReadWriteCloser, error) {
	if fileName == "Open error.txt" {
		return nil, pkgErr.FileNotFound{}
	}
	// Creating temporary directory for tests
	d := s.t.TempDir()
	_ = os.Chdir(d)

	// Creating file in the temp directory
	fd, _ := os.Create(fileName)

	return fd, nil
}

func (s mockSftpClient) Create(fileName string) (io.ReadWriteCloser, error) {
	if fileName == "Create error.txt" {
		return nil, errors.New("error in creating the file")
	}
	// Creating temporary directory for tests
	d := s.t.TempDir()
	_ = os.Chdir(d)

	// Creating file in the temp directory
	fd, _ := os.Create(fileName)

	return fd, nil
}

type mockFileInfo struct {
	name string
	os.FileInfo
}

func (m mockFileInfo) Name() string {
	return m.name
}

func (s mockSftpClient) ReadDir(dirName string) ([]os.FileInfo, error) {
	if dirName == "ErrorDirectory" {
		return nil, errors.New("error while reading directory")
	}

	files := make([]os.FileInfo, 0)
	m1 := mockFileInfo{name: "test1.txt"}
	m2 := mockFileInfo{name: "test2.txt"}
	files = append(files, m1, m2)

	return files, nil
}

func Test_Fetch(t *testing.T) {
	s1 := &sftp{
		filename: "Open error.txt",
		fileMode: "r",
		client:   mockSftpClient{},
	}
	s2 := &sftp{
		filename: "Copy error.txt",
		fileMode: "r",
		client:   mockSftpClient{t: t},
	}
	s3 := &sftp{
		filename: "test2.txt",
		fileMode: "rw",
		client:   mockSftpClient{t: t},
	}

	openErr := pkgErr.FileNotFound{}
	copyErr := errors.New("invalid argument")

	testcases := []struct {
		desc   string
		s      *sftp
		expErr error
	}{
		{"Open Error", s1, openErr},
		{"Copy Error", s2, copyErr},
		{"Success", s3, nil},
	}
	for i, tc := range testcases {
		l := newLocalFile(tc.s.filename, tc.s.fileMode)
		_ = l.Open()
		err := tc.s.fetch(l.FD)
		assert.Equal(t, tc.expErr, err, "Test [%v] failed. Expected: %v, got: %v,", i, tc.expErr, err)
	}
}

func Test_Push(t *testing.T) {
	s1 := &sftp{
		filename: "Create error.txt",
		fileMode: "r",
		client:   mockSftpClient{},
	}
	s2 := &sftp{
		filename: "Copy error.txt",
		fileMode: "r",
		client:   mockSftpClient{t: t},
	}
	s3 := &sftp{
		filename: "test1.txt",
		fileMode: "rw",
		client:   mockSftpClient{t: t},
	}
	createErr := errors.New("error in creating the file")
	copyErr := errors.New("invalid argument")
	testcases := []struct {
		desc   string
		s      *sftp
		expErr error
	}{
		{"Create Error", s1, createErr},
		{"Copy Error", s2, copyErr},
		{"Success", s3, nil},
	}

	for i, tc := range testcases {
		l := newLocalFile(tc.s.filename, tc.s.fileMode)
		_ = l.Open()
		err := tc.s.push(l.FD)
		assert.Equal(t, tc.expErr, err, "Test [%v] failed. Expected: %v, got: %v,", i, tc.expErr, err)
	}
}

func Test_SftpList(t *testing.T) {
	s := &sftp{
		filename: "",
		fileMode: "",
		client:   mockSftpClient{},
	}
	// Creating temporary directory for tests
	d := t.TempDir()
	_ = os.Chdir(d)

	// Creating two files in the temp directory
	_, _ = os.Create("test1.txt")
	_, _ = os.Create("test2.txt")

	dirErr := errors.New("error while reading directory")

	testcases := []struct {
		desc    string
		dirName string
		expErr  error
	}{
		{"Read Error", "ErrorDirectory", dirErr},
		{"Success", d, nil},
	}
	for i, tc := range testcases {
		_, err := s.list(tc.dirName)
		assert.Equal(t, tc.expErr, err, "Test [%v] failed. Expected: %v, got: %v,", i, tc.expErr, err)
	}
}
