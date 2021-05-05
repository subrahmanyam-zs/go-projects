package customer

import (
	"context"

	"github.com/zopsmart/gofr/examples/using-mongo/entity"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct{}

// Get returns the list of models from mongodb based on the filter passed in the request
func (m Customer) Get(c *gofr.Context, name string) ([]*entity.Customer, error) {
	results := make([]*entity.Customer, 0)

	// fetch the Mongo collection
	collection := c.MongoDB.Collection("customers")

	filter := bson.D{}

	if name != "" {
		nameFilter := primitive.E{
			Key:   "name",
			Value: name,
		}
		filter = append(filter, nameFilter)
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return results, errors.DB{Err: err}
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var model entity.Customer
		if err := cur.Decode(&model); err != nil {
			return results, errors.DB{Err: err}
		}

		results = append(results, &model)
	}

	return results, nil
}

// Create extracts JSON content from request body and unmarshal it as Customer and then put it into db
func (m Customer) Create(c *gofr.Context, model *entity.Customer) error {
	// fetch the Mongo collection
	collection := c.MongoDB.Collection("customers")

	_, err := collection.InsertOne(context.TODO(), model)

	return err
}

// Delete deletes a record from MongoDB, returns delete count and the error if it fails to delete
func (m Customer) Delete(c *gofr.Context, name string) (int, error) {
	// fetch the Mongo collection
	collection := c.MongoDB.Collection("customers")
	filter := bson.D{}

	filter = append(filter, primitive.E{
		Key:   "name",
		Value: name,
	})

	deleted, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return 0, errors.DB{Err: err}
	}

	return int(deleted.DeletedCount), nil
}
