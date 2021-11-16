package person

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/model"
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type person struct {
	table string
}

// New factory function for person store
func New(table string) store.Person {
	return person{table: table}
}

func (p person) Create(c *gofr.Context, person model.Person) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id":    {S: aws.String(person.ID)},
			"name":  {S: aws.String(person.Name)},
			"email": {S: aws.String(person.Email)},
		},
		TableName: aws.String(p.table),
	}

	_, err := c.DynamoDB.PutItem(input)

	return err
}

func (p person) Get(c *gofr.Context, id string) (model.Person, error) {
	input := &dynamodb.GetItemInput{
		AttributesToGet: []*string{aws.String("id"), aws.String("name"), aws.String("email")},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: aws.String(p.table),
	}

	var person model.Person

	out, err := c.DynamoDB.GetItem(input)
	if err != nil {
		return person, errors.DB{Err: err}
	}

	err = dynamodbattribute.UnmarshalMap(out.Item, &person)
	if err != nil {
		return person, errors.DB{Err: err}
	}

	return person, nil
}

func (p person) Update(c *gofr.Context, person model.Person) error {
	input := &dynamodb.UpdateItemInput{
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			"name":  {Value: &dynamodb.AttributeValue{S: aws.String(person.Name)}, Action: aws.String("PUT")},
			"email": {Value: &dynamodb.AttributeValue{S: aws.String(person.Email)}, Action: aws.String("PUT")},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(person.ID)},
		},
		TableName: aws.String(p.table),
	}

	_, err := c.DynamoDB.UpdateItem(input)

	return err
}

func (p person) Delete(c *gofr.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: aws.String(p.table),
	}

	_, err := c.DynamoDB.DeleteItem(input)

	return err
}
