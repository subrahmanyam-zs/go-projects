package file

import (
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
)

const testFile = "/tmp/testData.txt"

func TestLocalFileOpen(t *testing.T) {
	testcases := []struct {
		filename string
		mode     Mode
		expErr   error
	}{
		{"newTest.txt", READ, &os.PathError{
			Op:   "open",
			Path: "/tmp/newTest.txt",
			Err:  syscall.ENOENT,
		}}, // opening a new file in read mode does not make sense!
		{"test.txt", WRITE, nil},
		{"test.txt", READ, nil},
		{"test1.txt", READWRITE, nil},
		{"test1.txt", APPEND, nil},
		{"test1.txt", "unknown", nil},
	}

	c := &config.MockConfig{Data: map[string]string{
		"FILE_STORE": "LOCAL",
	}}

	for _, tc := range testcases {
		f, err := NewWithConfig(c, "/tmp/"+tc.filename, tc.mode)
		if err != nil {
			t.Error(err)
		}

		err = f.Open()
		assert.Equal(t, tc.expErr, err)
	}
}

func TestLocal_WriteInReadMode(t *testing.T) {
	c := &config.MockConfig{Data: map[string]string{
		"FILE_STORE": "LOCAL",
	}}

	err := createTestFile(testFile, []byte("The quick brown fox jumps over the lazy dog"))
	if err != nil {
		t.Error(err)
	}

	f, err := NewWithConfig(c, testFile, READ)
	if err != nil {
		t.Error(err)
	}

	defer f.Close()

	err = f.Open()
	if err != nil {
		t.Error(err)
	}

	dataToWrite := []byte("The quick brown fox jumps over the lazy dog")

	_, err = f.Write(dataToWrite)
	if err == nil {
		t.Error("Expected error while writing to a Read only file!")
	}
}

func TestLocal_ReadInWriteMode(t *testing.T) {
	c := &config.MockConfig{Data: map[string]string{
		"FILE_STORE": "LOCAL",
	}}

	err := createTestFile(testFile, []byte("The quick brown fox jumps over the lazy dog"))
	if err != nil {
		t.Error(err)
	}

	f, err := NewWithConfig(c, testFile, WRITE)
	if err != nil {
		t.Error(err)
	}

	defer f.Close()

	err = f.Open()
	if err != nil {
		t.Error(err)
	}

	b := make([]byte, 50)
	if _, err = f.Read(b); err == nil {
		t.Error("Expected error while reading from a Write only file!")
	}
}

func TestNilFileDescriptor(t *testing.T) {
	file := &fileAbstractor{FD: nil}
	b := make([]byte, 50)
	offset := int64(2)
	whence := 0

	_, err := file.Read(b)
	if err == nil {
		t.Error("Expected error while Reading from nil file descriptor")
	}

	_, err = file.Write(b)
	if err == nil {
		t.Error("Expected error while Writing from nil file descriptor")
	}

	err = file.Close()
	if err == nil {
		t.Error("Expected error while Closing nil file descriptor")
	}

	_, err = file.Seek(offset, whence)
	if err == nil {
		t.Error("Expected error while seeking nil file descriptor")
	}
}

func TestNotNilFileDescriptor(t *testing.T) {
	mode := fetchLocalFileMode(READWRITE)

	tests := []struct {
		fileMode         int
		str              string
		appendOrOverride string
		output           string
	}{
		{mode, "Test read write ", "Override the existing string", "Override the existing string"},
		{mode | os.O_APPEND, "Test read write ", "Append in the existing string", "Test read write Append in the existing string"},
	}

	for i, tc := range tests {
		fileName := "/tmp/testFile.txt"
		b := performFileOps(t, tc.fileMode, fileName, tc.str, tc.appendOrOverride)
		// if the file has been opened in READWRITE mode then tc.str content should get overwritten by tc.appendOrOverride
		// if it doesn't happen then through an error
		if tc.fileMode == mode {
			if strings.Contains(string(b), tc.str) {
				t.Errorf("Unexpected string: %v", tc.str)
			}
		}

		if !strings.Contains(string(b), tc.output) {
			t.Errorf("Failed[%v]Expect %v got %v", i, tc.output, string(b))
		}
	}
}

func performFileOps(t *testing.T, fileMode int, fileName, str, appendOrOverride string) []byte {
	b := make([]byte, 60)
	offset := int64(0)
	whence := 0
	l := fileAbstractor{fileName: fileName, fileMode: fileMode}

	if err := l.Open(); err != nil {
		t.Error(err)
	}

	defer os.Remove(fileName)

	if _, err := l.Write([]byte(str)); err != nil {
		t.Error(err)
	}
	// offset is set to the start of the file
	if _, err := l.Seek(offset, whence); err != nil {
		t.Error(err)
	}

	if _, err := l.Write([]byte(appendOrOverride)); err != nil {
		t.Error(err)
	}
	// offset is set to the start of the file
	if _, err := l.Seek(offset, whence); err != nil {
		t.Error(err)
	}

	if _, err := l.Read(b); err != nil {
		t.Error(err)
	}

	return b
}

func TestLocal_Seek(t *testing.T) {
	err := createTestFile(testFile, []byte("The quick brown fox jumps over the lazy dog"))
	if err != nil {
		t.Error(err)
	}

	defer os.Remove(testFile)

	tests := []struct {
		mode   Mode
		offset int64
		whence int
	}{
		{READWRITE, 0, 0},
		{WRITE, 2, 0},
		{READ, 1, 2},
		{APPEND, 0, 0},
	}

	for i, tc := range tests {
		l := fileAbstractor{fileName: "testFile.txt", fileMode: fetchLocalFileMode(tc.mode)}
		if err := l.Open(); err != nil {
			t.Error(err)
		}

		offset, err := l.Seek(tc.offset, tc.whence)

		assert.Equal(t, tc.offset, offset, i)

		if err != nil {
			t.Errorf("expect nil got err %v", err)
		}

		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}
}

func createTestFile(filePath string, dataToWrite []byte) error {
	file, err := os.OpenFile(filePath, fetchLocalFileMode(WRITE), os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(dataToWrite)

	return err
}
