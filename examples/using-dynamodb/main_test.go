package main

import (
	"bytes"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestMain(m *testing.M) {
	k := gofr.New()

	table := "person"
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: aws.String("S")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: aws.String("HASH")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ReadCapacityUnits: aws.Int64(10), WriteCapacityUnits: aws.Int64(5)},
		TableName:             aws.String(table),
	}

	_, err := k.DynamoDB.CreateTable(input)
	if err != nil {
		k.Logger.Errorf("Failed creation of table %v, %v", table, err)
	}

	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	go main()
	time.Sleep(2 * time.Second)

	tcs := []struct {
		method        string
		endpoint      string
		expStatusCode int
		body          []byte
	}{
		{http.MethodPost, "person", http.StatusCreated, []byte(`{"id":"1", "name":  "gofr", "email": "gofr@zopsmart.com"}`)},
		{http.MethodGet, "person/1", http.StatusOK, nil},
		{http.MethodPut, "person/1", http.StatusOK, []byte(`{"id":"1", "name":  "gofr1", "email": "gofrone@zopsmart.com"}`)},
		{http.MethodDelete, "person/1", http.StatusNoContent, nil},
	}

	for i, tc := range tcs {
		req, _ := request.NewMock(tc.method, "http://localhost:9091/"+tc.endpoint, bytes.NewBuffer(tc.body))

		cl := http.Client{}

		resp, err := cl.Do(req)
		if err != nil {
			t.Error(err)
			continue
		}

		if resp.StatusCode != tc.expStatusCode {
			t.Errorf("Testcase[%v] Failed.\tExpected %v\tGot %v\n", i, tc.expStatusCode, resp.StatusCode)
		}

		_ = resp.Body.Close()
	}
}
