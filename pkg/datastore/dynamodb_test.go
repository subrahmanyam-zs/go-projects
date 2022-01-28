package datastore

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func newDynamoDB(t *testing.T) DynamoDB {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	cfg := DynamoDBConfig{
		Region:            c.Get("DYNAMODB_REGION"),
		Endpoint:          c.Get("DYNAMODB_ENDPOINT_URL"),
		AccessKeyID:       c.Get("DYNAMODB_ACCESS_KEY_ID"),
		SecretAccessKey:   c.Get("DYNAMODB_SECRET_ACCESS_KEY"),
		ConnRetryDuration: 0,
	}

	db, err := NewDynamoDB(logger, cfg)
	if err != nil {
		t.Errorf("error in making connection to DynamoDB, %v", err)
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: aws.String("S")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: aws.String("HASH")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ReadCapacityUnits: aws.Int64(10), WriteCapacityUnits: aws.Int64(5)},
		TableName:             aws.String("test"),
	}

	_, err = db.CreateTable(input)
	if err != nil {
		a, ok := err.(awserr.Error)
		if ok && a.Code() != dynamodb.ErrCodeResourceInUseException {
			t.Errorf("Failed creation of table, %v", err)
		}
	}

	return db
}

func TestGetNewDynamoDB(t *testing.T) {
	tcs := []struct {
		region string
		err    error
	}{
		{"", awserr.New("MissingRegion", "could not find region configuration", nil)},
		{"ap-south-1", nil},
	}

	for i, tc := range tcs {
		cfg := DynamoDBConfig{
			Region:            tc.region,
			Endpoint:          "http://localhost:2021",
			AccessKeyID:       "access-key-id",
			SecretAccessKey:   "secret-key",
			ConnRetryDuration: 5,
		}

		_, err := NewDynamoDB(log.NewLogger(), cfg)

		assert.IsType(t, tc.err, err, "TESTCASE[%d], failed.\n", i)
	}
}

func TestHealthCheck(t *testing.T) {
	tcs := []struct {
		accessKey string
		secretKey string
		status    string
	}{
		{"access-key-id", "secret-key", pkg.StatusUp},
		{"", "", pkg.StatusDown},
	}

	for i, tc := range tcs {
		cfg := DynamoDBConfig{
			Region:            "ap-south-1",
			Endpoint:          "http://localhost:2021",
			AccessKeyID:       tc.accessKey,
			SecretAccessKey:   tc.secretKey,
			ConnRetryDuration: 5,
		}
		db, _ := NewDynamoDB(log.NewLogger(), cfg)

		health := db.HealthCheck()
		expHealth := types.Health{Name: pkg.DynamoDB, Status: tc.status}

		assert.Equal(t, expHealth, health, "TEST[%d], failed.\n", i)
	}
}

func TestDynamoDB_PutItem(t *testing.T) {
	db := newDynamoDB(t)

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id":   {S: aws.String("1")},
			"name": {S: aws.String("test1")},
		},
		TableName: aws.String("test"),
	}

	_, err := db.PutItem(input)
	if err != nil {
		t.Errorf("PutItem got unexpected error, %v", err)
	}

	req, _ := db.PutItemRequest(input)
	if err = req.Send(); err != nil {
		t.Errorf("PutItemRequest got unexpected error, %v", err)
	}

	_, err = db.PutItemWithContext(context.Background(), input)
	if err != nil {
		t.Errorf("PutItemWithContext got unexpected error, %v", err)
	}
}

// nolint:dupl //read and delete performed on the same entity
func TestDynamoDB_GetItem(t *testing.T) {
	db := newDynamoDB(t)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String("1")},
		},
		TableName: aws.String("test"),
	}

	_, err := db.GetItem(input)
	if err != nil {
		t.Errorf("GetItem got unexpected error, %v", err)
	}

	req, _ := db.GetItemRequest(input)
	if err = req.Send(); err != nil {
		t.Errorf("GetItemRequest got unexpected error, %v", err)
	}

	_, err = db.GetItemWithContext(context.Background(), input)
	if err != nil {
		t.Errorf("GetItemWithContext got unexpected error, %v", err)
	}
}

// nolint:dupl //read and delete performed on the same entity
func TestDynamoDB_DeleteItem(t *testing.T) {
	db := newDynamoDB(t)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String("1")},
		},
		TableName: aws.String("test"),
	}

	_, err := db.DeleteItem(input)
	if err != nil {
		t.Errorf("DeleteItem got unexpected error, %v", err)
	}

	req, _ := db.DeleteItemRequest(input)
	if err = req.Send(); err != nil {
		t.Errorf("DeleteItemRequest got unexpected error, %v", err)
	}

	_, err = db.DeleteItemWithContext(context.Background(), input)
	if err != nil {
		t.Errorf("DeleteItemWithContext got unexpected error, %v", err)
	}
}

func TestDynamoDB_UpdateItem(t *testing.T) {
	db := newDynamoDB(t)

	input := &dynamodb.UpdateItemInput{
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			"name": {Value: &dynamodb.AttributeValue{S: aws.String("test name")}, Action: aws.String("PUT")},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String("1")},
		},
		TableName: aws.String("test"),
	}

	_, err := db.UpdateItem(input)
	if err != nil {
		t.Errorf("UpdateItem got unexpected error, %v", err)
	}

	req, _ := db.UpdateItemRequest(input)
	if err = req.Send(); err != nil {
		t.Errorf("UpdateItemRequest got unexpected error, %v", err)
	}

	_, err = db.UpdateItemWithContext(context.Background(), input)
	if err != nil {
		t.Errorf("UpdateItemWithContext got unexpected error, %v", err)
	}
}

func Test_genPutItemQuery(t *testing.T) {
	input := &dynamodb.PutItemInput{
		ConditionExpression: aws.String("NOT contains(id, :id)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {S: aws.String("1")},
		},
		Item: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String("1")},
		},
		TableName: aws.String("test"),
	}

	expQuery := []string{"PutItem", "Item Fields {id}", "ConditionExpression NOT contains(id, :id)", "test"}

	query := genPutItemQuery(input)

	assert.Equal(t, expQuery, query, "Test Failed")
}

func Test_genGetItemQuery(t *testing.T) {
	tcs := []struct {
		in       []*string
		expQuery []string
	}{
		{
			nil,
			[]string{"GetItem", "Key {id}", "test"},
		},
		{
			[]*string{aws.String("id"), aws.String("name")},
			[]string{"GetItem", "AttributesToGet {id, name}", "Key {id}", "test"},
		},
	}

	for i, tc := range tcs {
		input := &dynamodb.GetItemInput{
			AttributesToGet: tc.in,
			Key: map[string]*dynamodb.AttributeValue{
				"id": {S: aws.String("1")},
			},
			TableName: aws.String("test"),
		}

		query := genGetItemQuery(input)

		assert.Equal(t, tc.expQuery, query, "TESTCASE[%v]", i)
	}
}

func Test_genDeleteItemQuery(t *testing.T) {
	input := &dynamodb.DeleteItemInput{
		ConditionExpression: aws.String("NOT contains(email, :e_email)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":e_email": {S: aws.String("test@gmail.com")},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String("1")},
		},
		TableName: aws.String("test"),
	}

	expQuery := []string{"DeleteItem", "ConditionExpression NOT contains(email, :e_email)", "Key {id}", "test"}

	query := genDeleteItemQuery(input)

	assert.Equal(t, expQuery, query, "Test Failed")
}

func Test_genUpdateItemQuery(t *testing.T) {
	tcs := []struct {
		attributesToUpdate  map[string]*dynamodb.AttributeValueUpdate
		conditionExpression *string
		attributeValues     map[string]*dynamodb.AttributeValue
		updateExpression    *string
		expQuery            []string
	}{
		{
			map[string]*dynamodb.AttributeValueUpdate{
				"name": {Value: &dynamodb.AttributeValue{S: aws.String("test name")}, Action: aws.String("PUT")},
			},
			nil,
			nil,
			nil,
			[]string{"UpdateItem", "AttributesToUpdate {name}", "Key {id}", "test"},
		},
		{
			map[string]*dynamodb.AttributeValueUpdate{
				"name": {Value: &dynamodb.AttributeValue{S: aws.String("test name")}, Action: aws.String("PUT")},
			},
			aws.String("id != :id"),
			map[string]*dynamodb.AttributeValue{":id": {S: aws.String("1")}},
			nil,
			[]string{"UpdateItem", "AttributesToUpdate {name}", "ConditionExpression id != :id", "Key {id}", "test"},
		},
		{
			map[string]*dynamodb.AttributeValueUpdate{
				"name": {Value: &dynamodb.AttributeValue{S: aws.String("test name")}, Action: aws.String("PUT")},
			},
			nil,
			map[string]*dynamodb.AttributeValue{
				":e_email": {S: aws.String("test@gmail.com")},
			},
			aws.String("SET email = :e_email"),
			[]string{"UpdateItem", "AttributesToUpdate {name}", "UpdateExpression SET email = :e_email", "Key {id}", "test"},
		},
	}

	for i, tc := range tcs {
		input := &dynamodb.UpdateItemInput{
			AttributeUpdates:          tc.attributesToUpdate,
			ConditionExpression:       tc.conditionExpression,
			UpdateExpression:          tc.updateExpression,
			ExpressionAttributeValues: tc.attributeValues,
			Key: map[string]*dynamodb.AttributeValue{
				"id": {S: aws.String("1")},
			},
			TableName: aws.String("test"),
		}

		query := genUpdateItemQuery(input)

		assert.Equal(t, tc.expQuery, query, "TESTCASE[%v]", i)
	}
}

func Test_monitorQuery(t *testing.T) {
	db := newDynamoDB(t)

	b := new(bytes.Buffer)
	db.logger = log.NewMockLogger(b)

	input := &dynamodb.GetItemInput{
		AttributesToGet: []*string{aws.String("id"), aws.String("name"), aws.String("email")},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String("1")},
		},
		TableName: aws.String("test"),
	}

	expLog := "GetItem - with AttributesToGet {id, name, email}, Key {id}, on table test"

	_, _ = db.GetItem(input)

	assert.Contains(t, b.String(), expLog, "TEST Failed")
}
