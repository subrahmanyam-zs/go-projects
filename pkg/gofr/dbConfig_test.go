package gofr

import (
	"crypto/tls"
	awssns "developer.zopsmart.com/go/gofr/pkg/notifier/aws-sns"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/kafka"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"github.com/gocql/gocql"
)

func Test_cassandraConfigFromEnv(t *testing.T) {
	testCases := []struct {
		name           string
		configLoc      Config
		expectedConfig datastore.CassandraCfg
		expectedError  bool
	}{
		{"success", &config.MockConfig{
			Data: map[string]string{"CASS_DB_HOST": "Host", "CASS_DB_PORT": "90012", "CASS_DB_USER": "cass", "CASS_DB_PASS": "cass123",
				"CASS_DB_KEYSPACE": "keyspace", "CASS_DB_CONSISTENCY": "cass_consistency", "RetryPolicy": "5",
				"CASS_DB_CERTIFICATE_FILE": "private node certificate path", "CASS_DB_KEY_FILE": "private node key path",
				"CASS_DB_ROOT_CERTIFICATE_FILE": "root certificate file path", "CASS_DB_INSECURE_SKIP_VERIFY": "false",
				"DATA_CENTER": "Cassandra"},
		}, datastore.CassandraCfg{
			Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace", Consistency: "cass_consistency",
			Timeout: 600, ConnectTimeout: 600, RetryPolicy: &gocql.SimpleRetryPolicy{NumRetries: 5}, TLSVersion: tls.VersionTLS12,
			HostVerification: false, CertificateFile: "private node certificate path", KeyFile: "private node key path",
			RootCertificateFile: "root certificate file path", ConnRetryDuration: 30, DataCenter: "Cassandra",
		}, false},
		{"Failure due to EnableSSl", &config.MockConfig{
			Data: map[string]string{"CASS_DB_HOST": "Host", "CASS_DB_PORT": "90012", "CASS_DB_USER": "cass", "CASS_DB_PASS": "cass123",
				"CASS_DB_KEYSPACE": "keyspace", "CASS_DB_CONSISTENCY": "cass_consistency", "RetryPolicy": "5",
				"CASS_DB_CERTIFICATE_FILE": "private node certificate path", "CASS_DB_KEY_FILE": "private node key path",
				"CASS_DB_HOST_VERIFICATION": "true", "CASS_DB_ROOT_CERTIFICATE_FILE": "root certificate file path",
				"CASS_DB_INSECURE_SKIP_VERIFY": "true"},
		}, datastore.CassandraCfg{
			Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace", Consistency: "cass_consistency",
			Timeout: 600, ConnectTimeout: 600, RetryPolicy: &gocql.SimpleRetryPolicy{NumRetries: 5}, HostVerification: true,
			TLSVersion: tls.VersionTLS12, ConnRetryDuration: 30, CertificateFile: "private node certificate path",
			KeyFile: "private node key path", RootCertificateFile: "root certificate file path", InsecureSkipVerify: true,
		}, true},
		{
			"Failure due to HostVerification", &config.MockConfig{
				Data: map[string]string{"CASS_DB_HOST": "Host", "CASS_DB_PORT": "90012", "CASS_DB_USER": "cass", "CASS_DB_PASS": "cass123",
					"CASS_DB_KEYSPACE": "keyspace", "CASS_DB_CONSISTENCY": "cass_consistency", "RetryPolicy": "5"},
			}, datastore.CassandraCfg{Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace",
				Consistency: "cass_consistency", Timeout: 600, ConnectTimeout: 600, RetryPolicy: &gocql.SimpleRetryPolicy{NumRetries: 5},
				HostVerification: false, TLSVersion: tls.VersionTLS12, ConnRetryDuration: 30,
			}, true,
		},
	}

	for _, tc := range testCases {
		cassandraConfig := cassandraConfigFromEnv(tc.configLoc)
		if !reflect.DeepEqual(cassandraConfig, &tc.expectedConfig) {
			if tc.expectedError == false {
				t.Errorf("Fail:%vGot: %v,expected:%v", tc.name, cassandraConfig, tc.expectedConfig)
			}
		}
	}
}

// Test to check getYcqlConfigs function
func Test_GetYcqlConfigs(t *testing.T) {
	testCases := []struct {
		name           string
		configLoc      Config
		expectedConfig datastore.CassandraCfg
		expectedError  bool
	}{
		{"success", &config.MockConfig{
			Data: map[string]string{"CASS_DB_HOST": "Host", "CASS_DB_PORT": "90012", "CASS_DB_USER": "cass", "CASS_DB_PASS": "cass123",
				"CASS_DB_KEYSPACE": "keyspace", "CASS_DB_INSECURE_SKIP_VERIFY": "false", "CASS_DB_CERTIFICATE_FILE": "private node certificate path",
				"CASS_DB_KEY_FILE": "private node key path", "CASS_DB_ROOT_CERTIFICATE_FILE": "root certificate file path",
				"CASS_DB_HOST_VERIFICATION": "true", "DATA_CENTER": "US Central"},
		}, datastore.CassandraCfg{
			Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace", Timeout: 600,
			ConnectTimeout: 600, HostVerification: true, ConnRetryDuration: 30, CertificateFile: "private node certificate path",
			KeyFile: "private node key path", RootCertificateFile: "root certificate file path", DataCenter: "US Central",
		}, false},
		{"Failure due to User", &config.MockConfig{
			Data: map[string]string{
				"CASS_DB_HOST": "Host", "CASS_DB_PORT": "90012", "CASS_DB_USER": "cassUser", "CASS_DB_PASS": "cass123", "CASS_DB_KEYSPACE": "keyspace",
				"CASS_DB_CERTIFICATE_FILE": "private node certificate path", "CASS_DB_KEY_FILE": "private node key path",
				"CASS_DB_ROOT_CERTIFICATE_FILE": "root certificate file path", "CASS_DB_INSECURE_SKIP_VERIFY": "true"},
		}, datastore.CassandraCfg{
			Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace", Timeout: 600,
			ConnectTimeout: 600, HostVerification: false, ConnRetryDuration: 30, CertificateFile: "private node certificate path",
			KeyFile: "private node key path", RootCertificateFile: "root certificate file path", InsecureSkipVerify: true,
		}, true},
	}

	for i, tc := range testCases {
		cassandraConfig := getYcqlConfigs(tc.configLoc)
		if !reflect.DeepEqual(cassandraConfig, tc.expectedConfig) {
			if tc.expectedError == false {
				t.Errorf("Testcase[%v] Failed:%vGot: %v,expected:%v", i, tc.name, cassandraConfig, tc.expectedConfig)
			}
		}
	}
}

func Test_kafkaConfigFromEnv(t *testing.T) {
	expectedConfig := kafka.Config{
		Brokers:           "Host:2008,Host:2009",
		Topics:            []string{"test-topics"},
		ConnRetryDuration: 30,
		MaxRetry:          10,
		InitialOffsets:    kafka.OffsetOldest,
		GroupID:           "testing-dev-consumer",
	}
	kafkaConfig := kafkaConfigFromEnv(&config.MockConfig{
		Data: map[string]string{
			"KAFKA_HOSTS": "Host:2008,Host:2009",
			"KAFKA_TOPIC": "test-topics",
			"APP_NAME":    "testing",
			"APP_VERSION": "dev",
		},
	})

	if !reflect.DeepEqual(kafkaConfig, &expectedConfig) {
		t.Errorf("Got: %v,expected:%v", kafkaConfig, expectedConfig)
	}

	kafkaConfig = kafkaConfigFromEnv(&config.MockConfig{
		Data: map[string]string{
			"KAFKA_HOSTS":           "Host:2008,Host:2009",
			"KAFKA_TOPIC":           "test-topics",
			"APP_NAME":              "testing",
			"APP_VERSION":           "dev",
			"KAFKA_CONSUMER_OFFSET": "NEWEST",
		},
	})

	expectedConfig.InitialOffsets = kafka.OffsetNewest

	if !reflect.DeepEqual(kafkaConfig, &expectedConfig) {
		t.Errorf("Got: %v,expected:%v", kafkaConfig, expectedConfig)
	}
}

func Test_mongoDBConfigFromEnv(t *testing.T) {
	testCases := []struct {
		name           string
		configLoc      Config
		expectedConfig datastore.MongoConfig
		expectedError  bool
	}{
		{
			"success", &config.MockConfig{
				Data: map[string]string{
					"MONGO_DB_HOST": "Host",
					"MONGO_DB_PORT": "27001",
					"MONGO_DB_USER": "Rohan",
					"MONGO_DB_PASS": "Rohan123",
					"MONGO_DB_NAME": "testDb",
				},
			},
			datastore.MongoConfig{
				HostName:          "Host",
				Port:              "27001",
				Username:          "Rohan",
				Password:          "Rohan123",
				Database:          "testDb",
				SSL:               false,
				RetryWrites:       false,
				ConnRetryDuration: 30,
			},
			false,
		},
		{
			"failure due to SSL", &config.MockConfig{
				Data: map[string]string{
					"MONGO_DB_HOST": "Host",
					"MONGO_DB_PORT": "27001",
					"MONGO_DB_USER": "Rohan",
					"MONGO_DB_PASS": "rohan123",
					"MONGO_DB_NAME": "testDb",
				},
			}, datastore.MongoConfig{
				HostName:          "Host",
				Port:              "27001",
				Username:          "Rohan",
				Password:          "Rohan123",
				Database:          "testDb",
				SSL:               false,
				RetryWrites:       false,
				ConnRetryDuration: 30,
			},
			true,
		},
		{
			"failure due to RetryWrites", &config.MockConfig{
				Data: map[string]string{
					"MONGO_DB_HOST": "Host",
					"MONGO_DB_PORT": "27001",
					"MONGO_DB_USER": "Rohan",
					"MONGO_DB_PASS": "rohan123",
					"MONGO_DB_NAME": "testDb",
				},
			}, datastore.MongoConfig{
				HostName:          "Host",
				Port:              "27001",
				Username:          "Rohan",
				Password:          "Rohan123",
				Database:          "testDb",
				SSL:               false,
				RetryWrites:       true,
				ConnRetryDuration: 30,
			},
			true,
		},
	}

	for _, tc := range testCases {
		mongoConfig := mongoDBConfigFromEnv(tc.configLoc)

		if !reflect.DeepEqual(mongoConfig, &tc.expectedConfig) {
			if tc.expectedError == false {
				t.Errorf("Got: %v,expected:%v", mongoConfig, tc.expectedConfig)
			}
		}
	}
}

func Test_GetBoolEnv(t *testing.T) {
	testcases := []struct {
		env    string
		output bool
	}{
		{"true", true},
		{"false", false},
		{"", false},
		{"abc", false},
	}

	for _, tc := range testcases {
		output := getBool(tc.env)
		if output != tc.output {
			t.Errorf("Expected boolean %t Got %t", tc.output, output)
		}
	}
}

func Test_AWSSNSConfigFromEnv(t *testing.T) {
	expectedConfig := awssns.Config{
		AccessKeyID:     "AKIswe",
		SecretAccessKey: "Vccvsqwesdd",
		Region:          "us-east-1",
		Protocol:        "email",
		Endpoint:        "xyz@zopsmart.com",
		TopicArn:        "arn:aws:aws-sns:us-east-1:123456789:TestTopic1",
	}
	snsConfig := awsSNSConfigFromEnv(&config.MockConfig{
		Data: map[string]string{
			"SNS_ACCESS_KEY":        "AKIswe",
			"SNS_SECRET_ACCESS_KEY": "Vccvsqwesdd",
			"SNS_REGION":            "us-east-1",
			"SNS_PROTOCOL":          "email",
			"SNS_ENDPOINT":          "xyz@zopsmart.com",
			"SNS_TOPIC_ARN":         "arn:aws:aws-sns:us-east-1:123456789:TestTopic1",
		},
	})

	assert.Equal(t, expectedConfig,snsConfig)
}
