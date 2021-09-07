package file

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type mockClient struct{}

func (mc *mockClient) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	switch *params.Bucket {
	case "test-bucket-zs":
		return &s3.GetObjectOutput{Body: ioutil.NopCloser(strings.NewReader("Successful fetch"))}, nil
	default:
		return nil, errors.InvalidParam{Param: []string{"bucket"}}
	}
}

func (mc *mockClient) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	switch *params.Bucket {
	case "test-bucket-zs":
		return nil, nil
	default:
		return nil, errors.InvalidParam{Param: []string{"bucket"}}
	}
}

func Test_NewAWSFile(t *testing.T) {
	cfg := AWSConfig{
		AccessKey: "random-access-key",
		SecretKey: "random-secret-key",
		Bucket:    "test_bucket",
		Region:    "us-east-2",
	}
	filename := "testfile.txt"
	mode := READWRITE
	f := newAWSS3File(&cfg, filename, mode)

	if f.client == nil {
		t.Error("error expected not nil client")
	}
}

func TestAws_fetch(t *testing.T) {
	m := &mockClient{}
	tests := []struct {
		cfg *aws
		err error
	}{
		{&aws{fileName: "aws.txt", fileMode: APPEND, client: m, bucketName: "test-bucket-zs"}, nil},
		{&aws{fileName: "aws.txt", fileMode: READ, client: m, bucketName: "random-bucket"}, errors.InvalidParam{Param: []string{"bucket"}}},
	}

	for i, tc := range tests {
		l := newLocalFile(tc.cfg.fileName, tc.cfg.fileMode)
		_ = l.Open()
		err := tc.cfg.fetch(l.FD)
		assert.Equal(t, tc.err, err, i)

		_ = l.Close()
	}
}

func TestAws_push(t *testing.T) {
	m := &mockClient{}
	tests := []struct {
		cfg *aws
		err error
	}{
		{&aws{fileName: "aws.txt", fileMode: READWRITE, client: m, bucketName: "random-bucket"}, errors.InvalidParam{Param: []string{"bucket"}}},
		{&aws{fileName: "awstest.txt", fileMode: READ, client: m, bucketName: "test-bucket-zs"}, nil},
	}

	for i, tc := range tests {
		l := newLocalFile(tc.cfg.fileName, tc.cfg.fileMode)
		_ = l.Open()
		err := tc.cfg.push(l.FD)
		assert.Equal(t, tc.err, err, i)

		_ = l.Close()
	}
}
