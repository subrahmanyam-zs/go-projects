package file

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type gcp struct {
	client     storage.Client
	bucketName string
	object     string

	fileMode Mode
}

func newGCPFile(cfg *GCPConfig, object string, mode Mode) (*gcp, error) {
	c, err := storage.NewClient(context.Background(), option.WithCredentialsFile(cfg.GCPKey))
	if err != nil {
		return nil, err
	}

	gcpFile := &gcp{
		client:     *c,
		bucketName: cfg.BucketName,
		object:     object,
		fileMode:   mode,
	}

	return gcpFile, nil
}

func (g *gcp) fetch(fd *os.File) error {
	// download the gcp object from a bucket
	r, err := g.client.Bucket(g.bucketName).Object(g.object).NewReader(context.Background())
	if err != nil {
		return err
	}
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	_, err = fd.Write(data)

	return err
}

func (g *gcp) push(fd *os.File) error {
	w := g.client.Bucket(g.bucketName).Object(g.object).NewWriter(context.Background())
	if _, err := io.Copy(w, fd); err != nil {
		return err
	}

	err := w.Close()

	return err
}
