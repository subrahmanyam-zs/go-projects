package main

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	grpc2 "developer.zopsmart.com/go/gofr/examples/sample-grpc/handler/grpc"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"

	"google.golang.org/grpc"
)

func TestIntegration(t *testing.T) {
	go main()
	time.Sleep(time.Second * 5)

	tcs := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{"GET", "/example?id=1", 200, nil},
		{"GET", "/example?id=2", 404, nil},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, "http://localhost:9093/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("error on making request %v", err)
		}

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}

		_ = resp.Body.Close()
	}

	testClient(t)
}

func testClient(t *testing.T) {
	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		t.Errorf("did not connect: %s", err)
		return
	}

	defer conn.Close()

	c := grpc2.NewExampleServiceClient(conn)

	_, err = c.Get(context.TODO(), &grpc2.ID{Id: "1"})
	if err != nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}

	_, err = c.Get(context.TODO(), &grpc2.ID{Id: "2"})
	if err == nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}
}
