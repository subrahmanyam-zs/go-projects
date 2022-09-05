package file

import (
	"fmt"
	"io"
	"os"

	pkgSftp "github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type sftp struct {
	filename string
	fileMode Mode
	client   fileOp
}

type sftpClient struct {
	*pkgSftp.Client
}

type fileOp interface {
	Open(fileName string) (io.ReadWriteCloser, error)
	Create(fileName string) (io.ReadWriteCloser, error)
	ReadDir(dirName string) ([]os.FileInfo, error)
}

func (s sftpClient) Open(fileName string) (io.ReadWriteCloser, error) {
	return s.Client.Open(fileName)
}

func (s sftpClient) Create(fileName string) (io.ReadWriteCloser, error) {
	return s.Client.Create(fileName)
}

func (s sftpClient) ReadDir(dirName string) ([]os.FileInfo, error) {
	return s.Client.ReadDir(dirName)
}

func newSFTPFile(c *SFTPConfig, filename string, mode Mode) (*sftp, error) {
	sftpFile := &sftp{filename: filename, fileMode: mode}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec // using InsecureIgnoreHostKey to accept any host key
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	client, err := pkgSftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	sftpFile.client = sftpClient{Client: client}

	return sftpFile, nil
}

func (s *sftp) fetch(fd *os.File) error {
	srcFile, err := s.client.Open(s.filename)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	_, err = io.Copy(fd, srcFile)
	if err != nil {
		return err
	}

	return nil
}
func (s *sftp) push(fd *os.File) error {
	destFile, err := s.client.Create(s.filename)
	if err != nil {
		return err
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, fd)
	if err != nil {
		return err
	}

	return nil
}

func (s *sftp) list(folderName string) ([]string, error) {
	files := make([]string, 0)

	fInfo, err := s.client.ReadDir(folderName)
	if err != nil {
		return nil, err
	}

	for i := range fInfo {
		files = append(files, fInfo[i].Name())
	}

	return files, nil
}
