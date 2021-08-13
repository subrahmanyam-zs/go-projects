package person

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/model"
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/store"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func TestMain(m *testing.M) {
	k := gofr.New()

	tableName := "person"
	deleteTableInput := &dynamodb.DeleteTableInput{TableName: aws.String(tableName)}

	_, err := k.DynamoDB.DeleteTable(deleteTableInput)
	if err != nil {
		k.Logger.Errorf("error in deleting table, %v", err)
	}

	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: aws.String("S")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: aws.String("HASH")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ReadCapacityUnits: aws.Int64(10), WriteCapacityUnits: aws.Int64(5)},
		TableName:             aws.String(tableName),
	}

	_, err = k.DynamoDB.CreateTable(createTableInput)
	if err != nil {
		k.Logger.Errorf("Failed creation of table %v, %v", tableName, err)
	}

	os.Exit(m.Run())
}

func initializeTest(t *testing.T) (*gofr.Context, store.Person) {
	k := gofr.New()

	// RefreshTables
	seeder := datastore.NewSeeder(&k.DataStore, "../../db")
	seeder.RefreshDynamoDB(t, "person")

	return gofr.NewContext(nil, nil, k), New("person")
}

func TestGet(t *testing.T) {
	expOut := model.Person{ID: "1", Name: "Ponting", Email: "Ponting@gmail.com"}

	ctx, p := initializeTest(t)

	out, err := p.Get(ctx, "1")
	if err != nil {
		t.Errorf("Expected no error, \nGot %v\n", err)
	}

	if !reflect.DeepEqual(expOut, out) {
		t.Errorf("Expected output %v\nGot %v\n", expOut, out)
	}
}

func TestGet_Error(t *testing.T) {
	k := gofr.New()

	ctx := gofr.NewContext(nil, nil, k)
	p := New("dummy")

	_, err := p.Get(ctx, "1")

	assert.IsType(t, errors.DB{}, err)
}

func TestCreate(t *testing.T) {
	input := model.Person{ID: "7", Name: "john", Email: "john@gmail.com"}

	ctx, p := initializeTest(t)

	err := p.Create(ctx, input)
	if err != nil {
		t.Errorf("Expected no error\nGot %v\n", err)
	}
}

func TestUpdate(t *testing.T) {
	input := model.Person{ID: "1", Name: "Ponting", Email: "Ponting.gates@gmail.com"}

	ctx, p := initializeTest(t)

	err := p.Update(ctx, input)
	if err != nil {
		t.Errorf("Expected no error\nGot %v\n", err)
	}
}

func TestDelete(t *testing.T) {
	ctx, p := initializeTest(t)

	err := p.Delete(ctx, "1")
	if err != nil {
		t.Errorf("Failed\tExpected %v\nGot %v\n", nil, err)
	}
}
