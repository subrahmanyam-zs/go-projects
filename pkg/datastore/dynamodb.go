package datastore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

// DynamoDBConfig configuration for DynamoDB connection
type DynamoDBConfig struct {
	Region            string
	Endpoint          string
	AccessKeyID       string
	SecretAccessKey   string
	ConnRetryDuration int
}

type DynamoDB struct {
	*dynamodb.DynamoDB
	logger log.Logger
	config DynamoDBConfig
}

// NewDynamoDB connects to DynamoDB and returns the connection
func NewDynamoDB(logger log.Logger, c DynamoDBConfig) (DynamoDB, error) {
	sessionConfig := &aws.Config{
		Region:      aws.String(c.Region),
		Logger:      logger,
		Endpoint:    aws.String(c.Endpoint),
		Credentials: credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, ""),
	}

	sess, err := session.NewSession(sessionConfig)
	if err != nil {
		return DynamoDB{}, err
	}

	db := dynamodb.New(sess)

	// check the db connection
	_, err = db.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return DynamoDB{}, err
	}

	return DynamoDB{DynamoDB: db, logger: logger, config: c}, nil
}

// HealthCheck checks health of the Dya
func (d DynamoDB) HealthCheck() types.Health {
	resp := types.Health{
		Name:   pkg.DynamoDB,
		Status: pkg.StatusDown,
	}

	// check if DynamoDB instance has been created during initialization
	if d.DynamoDB == nil {
		return resp
	}

	_, err := d.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return resp
	}

	resp.Status = pkg.StatusUp

	return resp
}
