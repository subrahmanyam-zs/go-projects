package file

import "os"

type File interface {
	// Open should open the file in the provided mode. Implementation depends on the file storage to be used.
	Open() error
	// Read calls the internal file descriptor method to Read.
	Read([]byte) (int, error)
	// Write calls the internal file descriptor method to Write.
	Write([]byte) (int, error)
	// Seek calls the internal file descriptor method to Seek
	Seek(offset int64, whence int) (int64, error)
	// Close calls the internal file descriptor method to Close.
	Close() error
}

type cloudStore interface {
	fetch(fd *os.File) error
	push(fd *os.File) error
}
