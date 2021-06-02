package gofr

import (
	"crypto/tls"
	"strconv"
	"strings"

	"github.com/gocql/gocql"
	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/avro"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/eventhub"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/kafka"
)

// cassandraDBConfigFromEnv returns configuration from environment variables to client so it can connect to cassandra
func cassandraConfigFromEnv(c Config) *datastore.CassandraCfg {
	cassandraTimeout, err := strconv.Atoi(c.Get("CASS_DB_TIMEOUT"))
	if err != nil {
		// setting default timeout of 600 milliseconds
		cassandraTimeout = 600
	}

	cassandraConnTimeout, err := strconv.Atoi(c.Get("CASS_DB_CONN_TIMEOUT"))
	if err != nil {
		// setting default timeout of 600 milliseconds
		cassandraConnTimeout = 600
	}

	cassandraPort, err := strconv.Atoi(c.Get("CASS_DB_PORT"))
	if err != nil {
		// if any error, setting default port
		cassandraPort = 9042
	}

	const retries = 5

	cassandraConfig := datastore.CassandraCfg{
		Hosts:               c.Get("CASS_DB_HOST"),
		Port:                cassandraPort,
		Username:            c.Get("CASS_DB_USER"),
		Password:            c.Get("CASS_DB_PASS"),
		Keyspace:            c.Get("CASS_DB_KEYSPACE"),
		Consistency:         c.Get("CASS_DB_CONSISTENCY"),
		Timeout:             cassandraTimeout,
		ConnectTimeout:      cassandraConnTimeout,
		RetryPolicy:         &gocql.SimpleRetryPolicy{NumRetries: retries},
		TLSVersion:          setTLSVersion(c.Get("CASS_DB_TLS_VERSION")),
		HostVerification:    getBool(c.Get("CASS_DB_HOST_VERIFICATION")),
		ConnRetryDuration:   getRetryDuration(c.Get("CASS_CONN_RETRY")),
		CertificateFile:     c.Get("CASS_DB_CERTIFICATE_FILE"),
		KeyFile:             c.Get("CASS_DB_KEY_FILE"),
		RootCertificateFile: c.Get("CASS_DB_ROOT_CERTIFICATE_FILE"),
		InsecureSkipVerify:  getBool(c.Get("CASS_DB_INSECURE_SKIP_VERIFY")),
		DataCenter:          c.Get("DATA_CENTER"),
	}

	return &cassandraConfig
}

func getBool(val string) bool {
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return boolVal
}

func setTLSVersion(version string) uint16 {
	if version == "10" {
		return tls.VersionTLS10
	} else if version == "11" {
		return tls.VersionTLS11
	} else if version == "13" {
		return tls.VersionTLS13
	}

	return tls.VersionTLS12
}

// mongoDBConfigFromEnv returns configuration from environment variables to client so it can connect to MongoDB
func mongoDBConfigFromEnv(c Config) *datastore.MongoConfig {
	enableSSl, _ := strconv.ParseBool(c.Get("MONGO_DB_ENABLE_SSL"))
	retryWrites, _ := strconv.ParseBool(c.Get("MONGO_DB_RETRY_WRITES"))

	mongoConfig := datastore.MongoConfig{
		HostName:          c.Get("MONGO_DB_HOST"),
		Port:              c.Get("MONGO_DB_PORT"),
		Username:          c.Get("MONGO_DB_USER"),
		Password:          c.Get("MONGO_DB_PASS"),
		Database:          c.Get("MONGO_DB_NAME"),
		SSL:               enableSSl,
		RetryWrites:       retryWrites,
		ConnRetryDuration: getRetryDuration(c.Get("MONGO_CONN_RETRY")),
	}

	return &mongoConfig
}

// kafkaDBConfigFromEnv returns configuration from environment variables to client so it can connect to kafka
func kafkaConfigFromEnv(c Config) *kafka.Config {
	hosts := c.Get("KAFKA_HOSTS") // CSV string
	topic := c.Get("KAFKA_TOPIC") // CSV string
	retryFrequency, _ := strconv.Atoi(c.Get("KAFKA_RETRY_FREQUENCY"))
	maxRetry, _ := strconv.Atoi(c.GetOrDefault("KAFKA_MAX_RETRY", "10"))
	// consumer groupID generation using APP_NAME and APP_VERSION
	groupName := c.Get("KAFKA_CONSUMERGROUP_NAME")
	if groupName == "" {
		groupName = c.GetOrDefault("APP_NAME", pkg.DefaultAppName) + "-" + c.GetOrDefault("APP_VERSION", pkg.DefaultAppVersion) + "-consumer"
	}

	// converting the CSV string to slice of string
	topics := strings.Split(topic, ",")
	config := &kafka.Config{
		Brokers: hosts,
		SASL: kafka.SASLConfig{
			User:     c.Get("KAFKA_SASL_USER"),
			Password: c.Get("KAFKA_SASL_PASS"),
		},
		Topics:            topics,
		MaxRetry:          maxRetry,
		RetryFrequency:    retryFrequency,
		ConnRetryDuration: getRetryDuration(c.Get("KAFKA_CONN_RETRY")),
		InitialOffsets:    kafka.OffsetOldest,
		GroupID:           groupName,
	}

	offset := c.GetOrDefault("KAFKA_CONSUMER_OFFSET", "OLDEST")

	switch offset {
	case "NEWEST":
		config.InitialOffsets = kafka.OffsetNewest
	default:
		config.InitialOffsets = kafka.OffsetOldest
	}

	return config
}

// getElasticSearchConfigFromEnv returns configuration from environment variables to client so it can connect to elasticsearch
func getElasticSearchConfigFromEnv(c Config) datastore.ElasticSearchCfg {
	elasticSearchCfg := datastore.ElasticSearchCfg{
		Host: c.Get("ELASTIC_SEARCH_HOST"),
		User: c.Get("ELASTIC_SEARCH_USER"),
		Pass: c.Get("ELASTIC_SEARCH_PASS"),
	}

	elasticSearchCfg.Port, _ = strconv.Atoi(c.Get("ELASTIC_SEARCH_PORT"))
	elasticSearchCfg.ConnectionRetryDuration = getRetryDuration(c.Get("ELASTIC_SEARCH_CONN_RETRY"))

	return elasticSearchCfg
}

func avroConfigFromEnv(c Config) *avro.Config {
	return &avro.Config{
		URL:            c.Get("AVRO_SCHEMA_URL"),
		Version:        c.Get("AVRO_SCHEMA_VERSION"),
		Subject:        c.Get("AVRO_SUBJECT"),
		SchemaUser:     c.Get("AVRO_USER"),
		SchemaPassword: c.Get("AVRO_PASSWORD"),
	}
}

func eventhubConfigFromEnv(c Config) eventhub.Config {
	brokers := c.Get("EVENTHUB_NAMESPACE")
	topics := strings.Split(c.Get("EVENTHUB_NAME"), ",")

	return eventhub.Config{
		Namespace:         brokers,
		EventhubName:      topics[0],
		ClientID:          c.Get("AZURE_CLIENT_ID"),
		ClientSecret:      c.Get("AZURE_CLIENT_SECRET"),
		TenantID:          c.Get("AZURE_TENANT_ID"),
		SharedAccessName:  c.Get("EVENTHUB_SAS_NAME"),
		SharedAccessKey:   c.Get("EVENTHUB_SAS_KEY"),
		ConnRetryDuration: getRetryDuration(c.Get("EVENTHUB_CONN_RETRY")),
	}
}
