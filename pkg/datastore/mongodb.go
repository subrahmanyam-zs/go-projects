package datastore

import (
	ctx "context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/zopsmart/gofr/pkg"
	"github.com/zopsmart/gofr/pkg/gofr/types"
	"github.com/zopsmart/gofr/pkg/log"
	"github.com/zopsmart/gofr/pkg/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig holds the configurations for Mongo Connectivity
type MongoConfig struct {
	HostName          string
	Port              string
	Username          string
	Password          string
	Database          string
	SSL               bool
	RetryWrites       bool
	ConnRetryDuration int
}

// MongoDB is an interface for accessing the base functionality
type MongoDB interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
	Aggregate(ctx ctx.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	RunCommand(ctx ctx.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult
	RunCommandCursor(ctx ctx.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (*mongo.Cursor, error)
	HealthCheck() types.Health
	IsSet() bool
}

type mongodb struct {
	*mongo.Database
	config *MongoConfig
}

func getMongoConfigFromEnv() (*MongoConfig, error) {
	getBoolEnv := func(varName string) (bool, error) {
		val := os.Getenv(varName)
		if val == "" {
			return false, nil
		}

		return strconv.ParseBool(val)
	}

	enableSSL, err := getBoolEnv("MONGO_DB_ENABLE_SSL")
	if err != nil {
		return nil, err
	}

	retryWrites, err := getBoolEnv("MONGO_DB_RETRY_WRITES")
	if err != nil {
		return nil, err
	}

	mongoConfig := MongoConfig{
		HostName:    os.Getenv("MONGO_DB_HOST"),
		Port:        os.Getenv("MONGO_DB_PORT"),
		Username:    os.Getenv("MONGO_DB_USER"),
		Password:    os.Getenv("MONGO_DB_PASS"),
		Database:    os.Getenv("MONGO_DB_NAME"),
		SSL:         enableSSL,
		RetryWrites: retryWrites,
	}

	return &mongoConfig, nil
}

// GetMongoDBFromEnv returns client to connect to MongoDB using configuration from environment variables
//Deprecated: Instead use datastore.GetNewMongoDB
func GetMongoDBFromEnv(logger log.Logger) (MongoDB, error) {
	// pushing deprecated feature count to prometheus
	middleware.PushDeprecatedFeature("GetMongoDBFromEnv")

	mongoConfig, err := getMongoConfigFromEnv()
	if err != nil {
		return mongodb{config: mongoConfig}, err
	}

	return GetNewMongoDB(logger, mongoConfig)
}

func getMongoConnectionString(config *MongoConfig) string {
	mongoConnectionString := fmt.Sprintf("mongodb://%v:%v@%v:%v/?ssl=%v&retrywrites=%v",
		config.Username,
		config.Password,
		config.HostName,
		config.Port,
		config.SSL,
		config.RetryWrites,
	)

	return mongoConnectionString
}

// GetNewMongoDB connects to MongoDB and returns the connection with the specified database in the configuration
func GetNewMongoDB(logger log.Logger, config *MongoConfig) (MongoDB, error) {
	mongoConnectionString := getMongoConnectionString(config)
	// set client options
	clientOptions := options.Client().ApplyURI(mongoConnectionString)

	const defaultMongoTimeout = 3
	ctxWithTimeout, cancel := ctx.WithTimeout(ctx.Background(), time.Duration(defaultMongoTimeout)*time.Second)

	defer cancel()

	// connect to MongoDB
	client, err := mongo.Connect(ctxWithTimeout, clientOptions)

	if err != nil {
		logger.Errorf("could not connect to Mongo DB, HostName: %s, Port: %s, error: %v\n",
			config.HostName, config.Port, err)
		return mongodb{config: config}, err
	}

	// check the connection since Calling Connect does not block for server discovery. If you wish to know if a
	// MongoDB server has been found and connected to, use the Ping method
	err = client.Ping(ctxWithTimeout, nil)
	if err != nil {
		logger.Errorf("error while pinging to Mongo DB, HostName: %s, Port: %s, error: %v\n",
			config.HostName, config.Port, err)

		_ = client.Disconnect(ctxWithTimeout)

		return mongodb{config: config}, err
	}

	logger.Infof("Mongo DB connected. HostName: %s, Port: %s, Database: %s\n", config.HostName, config.Port, config.Database)

	db := mongodb{Database: client.Database(config.Database), config: config}

	return db, err
}

// IsSet is used to check if the connection to Mongo is made or not
func (m mongodb) IsSet() bool {
	return m.Database != nil // if connection is not nil, it will return true, if no connection, then false
}

func (m mongodb) HealthCheck() types.Health {
	resp := types.Health{
		Name:     pkg.Mongo,
		Status:   pkg.StatusDown,
		Host:     m.config.HostName,
		Database: m.config.Database,
	}
	// The following check is for the condition when the connection to MongoDB has not been made during initialization
	if m.Database == nil {
		return resp
	}

	err := m.Database.RunCommand(ctx.Background(), bson.D{struct {
		Key   string
		Value interface{}
	}{Key: "ping", Value: 1}}, nil).Err()

	if err != nil {
		return resp
	}

	resp.Status = pkg.StatusUp

	return resp
}
