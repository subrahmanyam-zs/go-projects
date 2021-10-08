package gofr

import (
	"io"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func Test_kafkaRetry(t *testing.T) {
	var k Gofr

	logger := log.NewMockLogger(io.Discard)
	k.Logger = logger
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	kafkaConfig := kafkaConfigFromEnv(c)
	avroConfig := avroConfigFromEnv(c)
	kafkaConfig.ConnRetryDuration = 1
	// for the failed case
	kafkaConfig.Brokers = "invalid-host"

	go kafkaRetry(kafkaConfig, avroConfig, &k)

	for i := 0; i < 5; i++ {
		time.Sleep(3 * time.Second)

		if k.PubSub != nil && k.PubSub.IsSet() {
			t.Errorf("FAILED, expected: Kafka initialization to fail, got: kafka initialized")
			break
		}
	}
	// for the success case
	kafkaConfig.Brokers = c.Get("KAFKA_HOSTS")

	for i := 0; i < 5; i++ {
		time.Sleep(3 * time.Second)

		if k.PubSub.IsSet() {
			break
		}
	}

	if !k.PubSub.IsSet() {
		t.Errorf("FAILED, expected: Kafka initialized successfully, got: kafka initialization failed")
	}
}

func Test_eventhubRetry(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	var k Gofr

	logger := log.NewMockLogger(io.Discard)
	k.Logger = logger
	conf := config.NewGoDotEnvProvider(logger, "../../configs")

	c := &config.MockConfig{Data: map[string]string{
		"EVENTHUB_NAME":       "healthcheck",
		"EVENTHUB_NAMESPACE":  "",
		"AZURE_CLIENT_ID":     "incorrect",
		"AZURE_CLIENT_SECRET": conf.Get("AZURE_CLIENT_SECRET"),
		"AZURE_TENANT_ID":     conf.Get("AZURE_TENANT_ID"),
		"PUBSUB_BACKEND":      "EVENTHUB",
	}}

	eventhubConfig := eventhubConfigFromEnv(c)
	eventhubConfig.ConnRetryDuration = 1
	// for the failed case
	go eventhubRetry(&eventhubConfig, nil, &k)

	for i := 0; i < 5; i++ {
		time.Sleep(3 * time.Second)

		if k.PubSub != nil && k.PubSub.IsSet() {
			t.Errorf("FAILED, expected: Eventhub initialization to fail, got: Eventhub initialized")
			break
		}
	}
	// for the success case
	eventhubConfig.Namespace = "zsmisc-dev"
	eventhubConfig.ClientID = conf.Get("AZURE_CLIENT_ID")

	for i := 0; i < 5; i++ {
		time.Sleep(3 * time.Second)

		if k.PubSub.IsSet() {
			break
		}
	}

	if !k.PubSub.IsSet() {
		t.Errorf("FAILED, expected: Eventhub initialized successfully, got: Eventhub initialization failed")
	}
}

func Test_mongoRetry(t *testing.T) {
	var k Gofr

	logger := log.NewMockLogger(io.Discard)
	k.Logger = logger
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	enableSSl, _ := strconv.ParseBool(c.Get("MONGO_DB_ENABLE_SSL"))
	retryWrites, _ := strconv.ParseBool(c.Get("MONGO_DB_RETRY_WRITES"))
	// for the failed case
	mongoCfg := datastore.MongoConfig{
		HostName:          "invalid-host",
		Port:              c.Get("MONGO_DB_PORT"),
		Username:          c.Get("MONGO_DB_USER"),
		Password:          c.Get("MONGO_DB_PASS"),
		Database:          c.Get("MONGO_DB_NAME"),
		SSL:               enableSSl,
		RetryWrites:       retryWrites,
		ConnRetryDuration: 1,
	}

	go mongoRetry(&mongoCfg, &k)

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.MongoDB != nil && k.MongoDB.IsSet() {
			t.Errorf("FAILED, expected: MongoDB initialization to fail, got: MongoDB initialized")
			break
		}
	}

	// for the success case
	mongoCfg.HostName = c.Get("MONGO_DB_HOST")

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.MongoDB.IsSet() {
			break
		}
	}

	if !k.MongoDB.IsSet() {
		t.Errorf("FAILED, expected: MongoDB initialized successfully, got: MongoDB initialization failed")
	}
}

func Test_cassandraRetry(t *testing.T) {
	var k Gofr

	logger := log.NewMockLogger(io.Discard)
	k.Logger = logger
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	cassandraCfg := cassandraConfigFromEnv(c)
	cassandraCfg.ConnRetryDuration = 1
	// for the failed case
	cassandraCfg.Hosts = ""

	go cassandraRetry(cassandraCfg, &k)

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.Cassandra.Session != nil {
			t.Errorf("FAILED, expected: Cassandra initialization to fail, got: cassandra initialized")
			break
		}
	}
	// for the success case
	cassandraCfg.Hosts = c.Get("CASS_DB_HOST")

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.Cassandra.Session != nil {
			break
		}
	}

	if k.Cassandra.Session == nil {
		t.Errorf("FAILED, expected: Cassandra initialized successfully, got: cassandra initialization failed")
	}
}

func Test_ycqlRetry(t *testing.T) {
	var k Gofr

	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	cassandraCfg := getYcqlConfigs(c)
	cassandraCfg.Port, _ = strconv.Atoi(c.Get("YCQL_DB_PORT"))
	cassandraCfg.Password = c.Get("YCQL_DB_PASS")
	cassandraCfg.Username = c.Get("YCQL_DB_USER")
	cassandraCfg.ConnRetryDuration = 1
	// for the failed case
	cassandraCfg.Hosts = "invalid-url"
	k.Logger = logger

	go yclRetry(&cassandraCfg, &k)

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.YCQL.Session != nil {
			t.Errorf("FAILED, expected: Ycql initialization to fail, got: Ycql initialized")
			break
		}
	}

	// for the success case
	cassandraCfg.Hosts = c.Get("CASS_DB_HOST")

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.YCQL.Session != nil {
			break
		}
	}

	if k.YCQL.Session == nil {
		t.Errorf("FAILED, expected: Ycql initialized successfully, got: Ycql initialization failed")
	}
}

func Test_ormRetry(t *testing.T) {
	var k Gofr

	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	// for the failed case
	dc := datastore.DBConfig{
		HostName:          "invalid-url",
		Username:          c.Get("DB_USER"),
		Password:          c.Get("DB_PASSWORD"),
		Database:          c.Get("DB_NAME"),
		Port:              c.Get("DB_PORT"),
		Dialect:           c.Get("DB_DIALECT"),
		SSL:               c.Get("DB_SSL"),
		ORM:               c.Get("DB_ORM"),
		ConnRetryDuration: 1,
	}

	k.Logger = logger

	go ormRetry(&dc, &k)
	time.Sleep(5 * time.Second)
	// for the failed case
	if k.GORM() != nil {
		t.Errorf("FAILED, expected: Orm initialization to fail, got: orm initialized")
	}

	// for the success case
	dc.HostName = c.Get("DB_HOST")

	time.Sleep(5 * time.Second)

	if k.GORM() == nil || (k.GORM() != nil && k.GORM().DB().Ping() != nil) {
		t.Errorf("FAILED, expected: Orm initialized successfully, got: orm initialization failed")
	}
}

// Testing sqlx retry mechanism
func Test_sqlxRetry(t *testing.T) {
	var k Gofr

	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	dc := datastore.DBConfig{
		HostName:          "invalid-url",
		Username:          c.Get("DB_USER"),
		Password:          c.Get("DB_PASSWORD"),
		Database:          c.Get("DB_NAME"),
		Port:              c.Get("DB_PORT"),
		Dialect:           c.Get("DB_DIALECT"),
		SSL:               c.Get("DB_SSL"),
		ORM:               c.Get("DB_ORM"),
		ConnRetryDuration: 1,
	}

	// for the failed case
	k.Logger = logger

	go sqlxRetry(&dc, &k)

	// Failure case
	if k.SQLX() != nil {
		t.Errorf("FAILED, expected: SQLX initialization to fail, got: sqlx initialized")
	}
	// for the success case
	dc.HostName = c.Get("DB_HOST")

	time.Sleep(5 * time.Second)

	if k.SQLX() == nil || (k.SQLX() != nil && k.SQLX().Ping() != nil) {
		t.Errorf("FAILED, expected: SQLX initialized successfully, got: sqlx initialization failed")
	}
}

func Test_redisRetry(t *testing.T) {
	var k Gofr

	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	redisConfig := datastore.RedisConfig{
		HostName:                "invalid-url",
		Port:                    c.Get("REDIS_PORT"),
		ConnectionRetryDuration: 1,
	}

	redisConfig.Options = new(redis.Options)
	redisConfig.Options.Addr = redisConfig.HostName + ":" + redisConfig.Port

	// for the failed case
	k.Logger = logger

	go redisRetry(&redisConfig, &k)

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.Redis != nil && k.Redis.IsSet() {
			t.Errorf("FAILED, expected: Redis initialization to fail, got: redis initialized")
			break
		}
	}
	// for the success case
	redisConfig.HostName = c.Get("REDIS_HOST")
	redisConfig.Options.Addr = redisConfig.HostName + ":" + redisConfig.Port

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.Redis.IsSet() {
			break
		}
	}

	if !k.Redis.IsSet() {
		t.Errorf("FAILED, expected: Redis initialized successfully, got: redis initialization failed")
	}
}

//nolint:gocognit //breaks the readability of the code
func Test_elasticSearchRetry(t *testing.T) {
	k := Gofr{Logger: log.NewMockLogger(io.Discard)}

	testcases := []struct {
		name        string
		config      datastore.ElasticSearchCfg
		expectedErr bool
	}{
		{"success", datastore.ElasticSearchCfg{Ports: []int{2012}, ConnectionRetryDuration: 1, Host: "localhost"},
			false},
		{"failure", datastore.ElasticSearchCfg{Ports: []int{2012}, ConnectionRetryDuration: 1, Host: "localhost",
			CloudID: "elastic-cloud-id"}, true},
	}

	for _, tc := range testcases {
		go elasticSearchRetry(&tc.config, &k)

		for i := 0; i < 5; i++ {
			time.Sleep(2 * time.Second)

			if k.Elasticsearch.Client != nil {
				break
			}
		}

		if !tc.expectedErr && k.Elasticsearch.Client == nil {
			t.Errorf("%s\nFAILED, Expected: successful initialization, Got: initialization failed", tc.name)
		}

		if tc.expectedErr && k.Elasticsearch.Client != nil {
			t.Errorf("%s\nFAILED, Expected: failed initialization, Got: initialization successful", tc.name)
		}
	}
}

func Test_AWSSNSRetry(t *testing.T) {
	var k Gofr

	logger := log.NewMockLogger(ioutil.Discard)
	k.Logger = logger
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	awsSNSConfig := awsSNSConfigFromEnv(c)
	awsSNSConfig.ConnRetryDuration = 1

	go awsSNSRetry(&awsSNSConfig, &k)

	for i := 0; i < 5; i++ {
		time.Sleep(3 * time.Second)

		if k.Notifier.IsSet() {
			break
		}
	}

	assert.True(t, k.Notifier.IsSet(), "FAILED, expected: AwsSNS initialized successfully, got: AwsSNS initialization failed")
}

func Test_dynamoRetry(t *testing.T) {
	k := New()

	dynamoConfig := datastore.DynamoDBConfig{
		ConnRetryDuration: 1,
	}

	go dynamoRetry(dynamoConfig, k)

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.DynamoDB.DynamoDB != nil {
			t.Errorf("FAILED, expected: DynamoDB initialization to fail, got: DynamoDB initialized")
			break
		}
	}

	k = NewWithConfig(config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs"))

	dynamoConfig = datastore.DynamoDBConfig{
		Region:            "ap-south-1",
		Endpoint:          "http://localhost:2021",
		SecretAccessKey:   "sample-secret-access-key",
		AccessKeyID:       "sample-access-key",
		ConnRetryDuration: 1,
	}

	go dynamoRetry(dynamoConfig, k)

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		if k.DynamoDB.DynamoDB != nil {
			break
		}
	}

	if k.DynamoDB.DynamoDB == nil {
		t.Errorf("FAILED, expected: DynamoDB initialized successfully, got: DynamoDB initialization failed")
	}
}
