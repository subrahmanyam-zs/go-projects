package file

import (
	"os"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type fileAbstractor struct {
	fileName             string
	fileMode             int
	FD                   *os.File
	remoteFileAbstracter cloudStore
}

func newLocalFile(filename string, mode Mode) *fileAbstractor {
	return &fileAbstractor{
		fileName: filename,
		fileMode: fetchLocalFileMode(mode),
	}
}

func (l *fileAbstractor) Open() error {
	file, err := os.OpenFile(l.fileName, l.fileMode, os.ModePerm)
	if err != nil {
		return err
	}

	l.FD = file

	if l.remoteFileAbstracter == nil {
		return err
	}

	fileMode := l.fileMode

	tmpFileMode := fetchLocalFileMode(READWRITE) // tmp file should be opened in WRITE mode for downloading and READ mode for uploading
	if l.fileMode == fetchLocalFileMode(APPEND) {
		tmpFileMode |= os.O_APPEND
	}

	tmpFileName := l.fileName + randomString()

	l.fileName = "/tmp/" + tmpFileName
	l.fileMode = tmpFileMode

	if _, err = os.OpenFile(l.fileName, l.fileMode, os.ModePerm); err != nil {
		return err
	}

	err = l.remoteFileAbstracter.fetch(l.FD)
	if err != nil && fileMode == fetchLocalFileMode(READ) {
		return err
	}

	// if a file is requested in READ mode, then the temp file should have only READ access. (Note that Data has been downloaded
	// to it, That means it needs write access to do that.)
	if l.fileMode == fetchLocalFileMode(READ) {
		_ = l.Close()
		l.fileMode = fetchLocalFileMode(READ) // tmpFile should also be in READ mode if azure file is in READ mode
		_ = l.Open()
	} else {
		_, err = l.Seek(startOffset, defaultWhence)
	}

	return err
}

func (l *fileAbstractor) Read(b []byte) (int, error) {
	if l.FD == nil {
		return 0, errors.FileNotFound{}
	}

	return l.FD.Read(b)
}

func (l *fileAbstractor) Write(b []byte) (int, error) {
	if l.FD == nil {
		return 0, errors.FileNotFound{}
	}

	return l.FD.Write(b)
}

func (l *fileAbstractor) Seek(offset int64, whence int) (int64, error) {
	if l.FD == nil {
		return 0, errors.FileNotFound{}
	}

	return l.FD.Seek(offset, whence)
}

func (l *fileAbstractor) Close() error {
	if l.FD == nil {
		return errors.FileNotFound{}
	}

	if l.remoteFileAbstracter == nil {
		return l.FD.Close()
	}

	if _, err := l.Seek(startOffset, defaultWhence); err != nil { // offset is set to the start of the file
		return err
	}

	err := l.remoteFileAbstracter.push(l.FD)
	if err != nil {
		return err
	}

	err = l.FD.Close()
	if err != nil {
		return err
	}

	return os.Remove(l.fileName)
}

func (l *fileAbstractor) List(directory string) ([]string, error) {
	files := make([]string, 0)

	if l.remoteFileAbstracter != nil {
		return l.remoteFileAbstracter.list(directory)
	}

	fInfo, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for i := range fInfo {
		files = append(files, fInfo[i].Name())
	}

	return files, nil
}
