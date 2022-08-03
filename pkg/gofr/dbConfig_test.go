package gofr

import (
	"crypto/tls"
	"io"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/datastore/kvdata"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/eventbridge"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/kafka"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
	awssns "developer.zopsmart.com/go/gofr/pkg/notifier/aws-sns"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func Test_cassandraConfigFromEnv(t *testing.T) {
	testCases := []struct {
		name           string
		configLoc      Config
		expectedConfig datastore.CassandraCfg
		prefix         string
		expectedError  bool
	}{
		// nolint:dupl // testcases are different but some values are same
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
		}, "", false},
		// nolint:dupl // testcases are different but some values are same
		{"success with prefix", &config.MockConfig{
			Data: map[string]string{"PRE_CASS_DB_HOST": "Host", "PRE_CASS_DB_PORT": "90012", "PRE_CASS_DB_USER": "cass",
				"PRE_CASS_DB_PASS": "cass123", "PRE_CASS_DB_KEYSPACE": "keyspace", "PRE_CASS_DB_CONSISTENCY": "cass_consistency",
				"PRE_RetryPolicy": "5", "PRE_CASS_DB_CERTIFICATE_FILE": "private node certificate path",
				"PRE_CASS_DB_KEY_FILE": "private node key path", "PRE_CASS_DB_ROOT_CERTIFICATE_FILE": "root certificate file path",
				"PRE_CASS_DB_INSECURE_SKIP_VERIFY": "false", "PRE_DATA_CENTER": "Cassandra"},
		}, datastore.CassandraCfg{
			Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace", Consistency: "cass_consistency",
			Timeout: 600, ConnectTimeout: 600, RetryPolicy: &gocql.SimpleRetryPolicy{NumRetries: 5}, TLSVersion: tls.VersionTLS12,
			HostVerification: false, CertificateFile: "private node certificate path", KeyFile: "private node key path",
			RootCertificateFile: "root certificate file path", ConnRetryDuration: 30, DataCenter: "Cassandra",
		}, "PRE", false},
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
		}, "", true},
		{
			"Failure due to HostVerification",
			&config.MockConfig{Data: map[string]string{"CASS_DB_HOST": "Host", "CASS_DB_PORT": "90012",
				"CASS_DB_USER": "cass", "CASS_DB_PASS": "cass123", "CASS_DB_KEYSPACE": "keyspace",
				"CASS_DB_CONSISTENCY": "cass_consistency", "RetryPolicy": "5"}},
			datastore.CassandraCfg{Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace",
				Consistency: "cass_consistency", Timeout: 600, ConnectTimeout: 600, RetryPolicy: &gocql.SimpleRetryPolicy{NumRetries: 5},
				HostVerification: false, TLSVersion: tls.VersionTLS12, ConnRetryDuration: 30,
			}, "", true,
		},
	}

	for _, tc := range testCases {
		cassandraConfig := cassandraConfigFromEnv(tc.configLoc, tc.prefix)
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
		prefix         string
		expectedConfig datastore.CassandraCfg
		expectedError  bool
	}{
		{"success", &config.MockConfig{
			Data: map[string]string{"CASS_DB_HOST": "Host", "CASS_DB_PORT": "90012", "CASS_DB_USER": "cass", "CASS_DB_PASS": "cass123",
				"CASS_DB_KEYSPACE": "keyspace", "CASS_DB_INSECURE_SKIP_VERIFY": "false", "CASS_DB_CERTIFICATE_FILE": "private node certificate path",
				"CASS_DB_KEY_FILE": "private node key path", "CASS_DB_ROOT_CERTIFICATE_FILE": "root certificate file path",
				"CASS_DB_HOST_VERIFICATION": "true", "DATA_CENTER": "US Central"},
		}, "", datastore.CassandraCfg{
			Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace", Timeout: 600,
			ConnectTimeout: 600, HostVerification: true, ConnRetryDuration: 30, CertificateFile: "private node certificate path",
			KeyFile: "private node key path", RootCertificateFile: "root certificate file path", DataCenter: "US Central",
		}, false},
		{"success with prefix", &config.MockConfig{
			Data: map[string]string{"PRE_CASS_DB_HOST": "Host", "PRE_CASS_DB_PORT": "90012", "PRE_CASS_DB_USER": "cass",
				"PRE_CASS_DB_PASS": "cass123", "PRE_CASS_DB_KEYSPACE": "keyspace", "PRE_CASS_DB_INSECURE_SKIP_VERIFY": "false",
				"PRE_CASS_DB_CERTIFICATE_FILE": "private node certificate path", "PRE_CASS_DB_KEY_FILE": "private node key path",
				"PRE_CASS_DB_ROOT_CERTIFICATE_FILE": "root certificate file path", "PRE_CASS_DB_HOST_VERIFICATION": "true",
				"PRE_DATA_CENTER": "US Central"},
		}, "PRE", datastore.CassandraCfg{
			Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace", Timeout: 600,
			ConnectTimeout: 600, HostVerification: true, ConnRetryDuration: 30, CertificateFile: "private node certificate path",
			KeyFile: "private node key path", RootCertificateFile: "root certificate file path", DataCenter: "US Central",
		}, false},
		{"Failure due to User", &config.MockConfig{
			Data: map[string]string{
				"CASS_DB_HOST": "Host", "CASS_DB_PORT": "90012", "CASS_DB_USER": "cassUser", "CASS_DB_PASS": "cass123", "CASS_DB_KEYSPACE": "keyspace",
				"CASS_DB_CERTIFICATE_FILE": "private node certificate path", "CASS_DB_KEY_FILE": "private node key path",
				"CASS_DB_ROOT_CERTIFICATE_FILE": "root certificate file path", "CASS_DB_INSECURE_SKIP_VERIFY": "true"},
		}, "", datastore.CassandraCfg{
			Hosts: "Host", Port: 90012, Username: "cass", Password: "cass123", Keyspace: "keyspace", Timeout: 600,
			ConnectTimeout: 600, HostVerification: false, ConnRetryDuration: 30, CertificateFile: "private node certificate path",
			KeyFile: "private node key path", RootCertificateFile: "root certificate file path", InsecureSkipVerify: true,
		}, true},
	}

	for i, tc := range testCases {
		cassandraConfig := getYcqlConfigs(tc.configLoc, tc.prefix)
		if !reflect.DeepEqual(cassandraConfig, tc.expectedConfig) {
			if tc.expectedError == false {
				t.Errorf("Testcase[%v] Failed:%vGot: %v,expected:%v", i, tc.name, cassandraConfig, tc.expectedConfig)
			}
		}
	}
}

func Test_kafkaConfigFromEnv(t *testing.T) {
	testcases := []struct {
		config         *config.MockConfig
		expectedConfig kafka.Config
	}{
		{
			&config.MockConfig{
				Data: map[string]string{"KAFKA_HOSTS": "Host:2008,Host:2009", "KAFKA_TOPIC": "test-topics",
					"APP_NAME": "testing", "APP_VERSION": "dev", "KAFKA_AUTOCOMMIT_DISABLE": "false"}},
			kafka.Config{
				Brokers: "Host:2008,Host:2009", Topics: []string{"test-topics"}, ConnRetryDuration: 30,
				MaxRetry: 10, InitialOffsets: kafka.OffsetOldest, GroupID: "testing-dev-consumer", DisableAutoCommit: false},
		},
		{&config.MockConfig{
			Data: map[string]string{"KAFKA_HOSTS": "Host:2008,Host:2009", "KAFKA_TOPIC": "test-topics", "APP_NAME": "testing",
				"APP_VERSION": "dev", "KAFKA_CONSUMER_OFFSET": "NEWEST"}},
			kafka.Config{
				Brokers: "Host:2008,Host:2009", Topics: []string{"test-topics"}, ConnRetryDuration: 30,
				MaxRetry: 10, InitialOffsets: kafka.OffsetNewest, GroupID: "testing-dev-consumer", DisableAutoCommit: false},
		},
	}
	for i, tc := range testcases {
		res := kafkaConfigFromEnv(tc.config, "")
		if !reflect.DeepEqual(res, &tc.expectedConfig) {
			t.Errorf("Test case failed [%v]. Got: %v,expected:%v", i, tc.config, tc.expectedConfig)
		}
	}
}

func Test_mongoDBConfigFromEnv(t *testing.T) {
	testCases := []struct {
		name           string
		configLoc      Config
		prefix         string
		expectedConfig datastore.MongoConfig
		expectedError  bool
	}{
		{
			"success",
			&config.MockConfig{Data: map[string]string{"MONGO_DB_HOST": "Host", "MONGO_DB_PORT": "27001",
				"MONGO_DB_USER": "Rohan", "MONGO_DB_PASS": "Rohan123", "MONGO_DB_NAME": "testDb"}},
			"",
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
			"success with prefix",
			&config.MockConfig{Data: map[string]string{"PRE_MONGO_DB_HOST": "Host", "PRE_MONGO_DB_PORT": "27001",
				"PRE_MONGO_DB_USER": "Rohan", "PRE_MONGO_DB_PASS": "Rohan123", "PRE_MONGO_DB_NAME": "testDb"}},
			"PRE",
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
			"failure due to SSL",
			&config.MockConfig{Data: map[string]string{"MONGO_DB_HOST": "Host", "MONGO_DB_PORT": "27001",
				"MONGO_DB_USER": "Rohan", "MONGO_DB_PASS": "rohan123", "MONGO_DB_NAME": "testDb"}},
			"",
			datastore.MongoConfig{HostName: "Host", Port: "27001", Username: "Rohan", Password: "Rohan123",
				Database: "testDb", SSL: false, RetryWrites: false, ConnRetryDuration: 30},
			true,
		},
		{
			"failure due to RetryWrites",
			&config.MockConfig{Data: map[string]string{"MONGO_DB_HOST": "Host", "MONGO_DB_PORT": "27001",
				"MONGO_DB_USER": "Rohan", "MONGO_DB_PASS": "rohan123", "MONGO_DB_NAME": "testDb"}},
			"",
			datastore.MongoConfig{
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
		mongoConfig := mongoDBConfigFromEnv(tc.configLoc, tc.prefix)

		if !reflect.DeepEqual(mongoConfig, &tc.expectedConfig) {
			if tc.expectedError == false {
				t.Errorf("Got: %v,expected:%v", mongoConfig, tc.expectedConfig)
			}
		}
	}
}

func Test_dynamoDBConfigFromEnv(t *testing.T) {
	expConfig := datastore.DynamoDBConfig{
		Region:            "ap-south-1",
		Endpoint:          "http://localhost:2021",
		AccessKeyID:       "access-key-id",
		SecretAccessKey:   "access-key",
		ConnRetryDuration: 2,
	}
	testcases := []struct {
		inputConfig *config.MockConfig
		prefix      string
	}{
		{&config.MockConfig{Data: map[string]string{"DYNAMODB_REGION": "ap-south-1",
			"DYNAMODB_ENDPOINT_URL": "http://localhost:2021", "DYNAMODB_ACCESS_KEY_ID": "access-key-id",
			"DYNAMODB_SECRET_ACCESS_KEY": "access-key", "DYNAMODB_CONN_RETRY": "2"},
		}, ""},
		{&config.MockConfig{Data: map[string]string{"PRE_DYNAMODB_REGION": "ap-south-1",
			"PRE_DYNAMODB_ENDPOINT_URL": "http://localhost:2021", "PRE_DYNAMODB_ACCESS_KEY_ID": "access-key-id",
			"PRE_DYNAMODB_SECRET_ACCESS_KEY": "access-key", "PRE_DYNAMODB_CONN_RETRY": "2"},
		}, "PRE"},
	}

	for i, tc := range testcases {
		cfg := dynamoDBConfigFromEnv(tc.inputConfig, tc.prefix)
		if !reflect.DeepEqual(cfg, expConfig) {
			t.Errorf("Test case failed [%v], Got: %v,expected:%v", i, cfg, expConfig)
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

func Test_getElasticSearchConfigFromEnv(t *testing.T) {
	testcases := []struct {
		input  Config
		output datastore.ElasticSearchCfg
		prefix string
	}{
		{
			&config.MockConfig{Data: map[string]string{"ELASTIC_SEARCH_HOST": "localhost",
				"ELASTIC_SEARCH_PORT": "2012", "ELASTIC_SEARCH_CONN_RETRY": "20"}},
			datastore.ElasticSearchCfg{Host: "localhost", Ports: []int{2012}, ConnectionRetryDuration: 20},
			"",
		},
		{
			&config.MockConfig{Data: map[string]string{"PRE_ELASTIC_SEARCH_HOST": "localhost",
				"PRE_ELASTIC_SEARCH_PORT": "2012", "PRE_ELASTIC_SEARCH_CONN_RETRY": "20"}},
			datastore.ElasticSearchCfg{Host: "localhost", Ports: []int{2012}, ConnectionRetryDuration: 20},
			"PRE",
		},
		{
			&config.MockConfig{Data: map[string]string{"ELASTIC_SEARCH_HOST": "localhost",
				"ELASTIC_SEARCH_PORT": "2012,2011,2010", "ELASTIC_SEARCH_CONN_RETRY": "20"}},
			datastore.ElasticSearchCfg{Host: "localhost", Ports: []int{2012, 2011, 2010}, ConnectionRetryDuration: 20},
			"",
		},
		{
			&config.MockConfig{Data: map[string]string{"ELASTIC_SEARCH_HOST": "localhost",
				"ELASTIC_SEARCH_PORT": "2012,2011,abc,2010", "ELASTIC_SEARCH_CONN_RETRY": "20"}},
			datastore.ElasticSearchCfg{Host: "localhost", Ports: []int{2012, 2011, 2010}, ConnectionRetryDuration: 20},
			"",
		},
		{
			&config.MockConfig{Data: map[string]string{"ELASTIC_SEARCH_CONN_RETRY": "20", "ELASTIC_CLOUD_ID": "sample-cloud-id"}},
			datastore.ElasticSearchCfg{Ports: []int{}, CloudID: "sample-cloud-id", ConnectionRetryDuration: 20},
			"",
		},
	}

	for i, tc := range testcases {
		output := elasticSearchConfigFromEnv(tc.input, tc.prefix)

		if !reflect.DeepEqual(output, tc.output) {
			t.Errorf("[TESTCASE%v] Failed.\nExpected:%v\nGot: %v", i+1, tc.output, output)
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
			"PRE_SNS_ACCESS_KEY":        "AKIswe",
			"PRE_SNS_SECRET_ACCESS_KEY": "Vccvsqwesdd",
			"PRE_SNS_REGION":            "us-east-1",
			"PRE_SNS_PROTOCOL":          "email",
			"PRE_SNS_ENDPOINT":          "xyz@zopsmart.com",
			"PRE_SNS_TOPIC_ARN":         "arn:aws:aws-sns:us-east-1:123456789:TestTopic1",
		},
	}, "PRE")

	assert.Equal(t, expectedConfig, snsConfig)
}

func Test_sqlDBConfigFromEnv(t *testing.T) {
	var (
		mc1 = &config.MockConfig{Data: map[string]string{"DB_HOST": "localhost", "DB_USER": "root", "DB_PASSWORD": "root123",
			"DB_NAME": "mysql", "DB_PORT": "3306", "DB_DIALECT": "mysql", "DB_MAX_OPEN_CONN": "10", "DB_MAX_IDLE_CONN": "10",
			"DB_CONN_RETRY": "5", "DB_MAX_CONN_LIFETIME": "100"}}
		mc2 = &config.MockConfig{Data: map[string]string{"DB_HOST": "localhost", "DB_USER": "root", "DB_PASSWORD": "root123",
			"DB_NAME": "mysql", "DB_PORT": "3306", "DB_DIALECT": "mysql", "DB_MAX_OPEN_CONN": "abc", "DB_MAX_IDLE_CONN": "20",
			"DB_CONN_RETRY": "5", "DB_MAX_CONN_LIFETIME": "100"}}
		mc3 = &config.MockConfig{Data: map[string]string{"DB_HOST": "localhost", "DB_USER": "root", "DB_PASSWORD": "root123",
			"DB_NAME": "mysql", "DB_PORT": "3306", "DB_DIALECT": "mysql", "DB_MAX_OPEN_CONN": "56.78", "DB_MAX_IDLE_CONN": "20.22",
			"DB_CONN_RETRY": "5", "DB_MAX_CONN_LIFETIME": "100.30"}}
		mc4 = &config.MockConfig{Data: map[string]string{"PRE_DB_HOST": "localhost", "PRE_DB_USER": "root", "PRE_DB_PASSWORD": "root123",
			"PRE_DB_NAME": "mysql", "PRE_DB_PORT": "3306", "PRE_DB_DIALECT": "mysql", "PRE_DB_MAX_OPEN_CONN": "10", "PRE_DB_MAX_IDLE_CONN": "10",
			"PRE_DB_CONN_RETRY": "5", "PRE_DB_MAX_CONN_LIFETIME": "100"}}
		c1 = &datastore.DBConfig{HostName: "localhost", Username: "root",
			Password: "root123", Database: "mysql", Port: "3306", Dialect: "mysql", ConnRetryDuration: 5, MaxOpenConn: 10,
			MaxIdleConn: 10, MaxConnLife: 100}
		c2 = &datastore.DBConfig{HostName: "localhost", Username: "root",
			Password: "root123", Database: "mysql", Port: "3306", Dialect: "mysql", ConnRetryDuration: 5, MaxOpenConn: 0,
			MaxIdleConn: 20, MaxConnLife: 100}
		c3 = &datastore.DBConfig{HostName: "localhost", Username: "root",
			Password: "root123", Database: "mysql", Port: "3306", Dialect: "mysql", ConnRetryDuration: 5, MaxOpenConn: 0,
			MaxIdleConn: 0, MaxConnLife: 0}
	)

	testcases := []struct {
		desc     string
		input    *config.MockConfig
		expDBCfg *datastore.DBConfig
		prefix   string
	}{
		{"valid configs", mc1, c1, ""},
		{"valid configs with prefix", mc4, c1, "PRE"},
		{"invalid config for DB_MAX_OPEN_CONN", mc2, c2, ""},
		{"invalid configs for sql connection pool", mc3, c3, ""},
	}

	for i, tc := range testcases {
		cfg := sqlDBConfigFromEnv(tc.input, tc.prefix)

		assert.Equal(t, tc.expDBCfg, cfg, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func Test_eventBridgeConfigFromEnv(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	c := &config.MockConfig{
		Data: map[string]string{
			"EVENT_BRIDGE_REGION":           "us-east-1",
			"EVENT_BRIDGE_BUS":              "Gofr",
			"EVENT_BRIDGE_SOURCE":           "Gofr-application",
			"EVENT_BRIDGE_RETRY_FREQUENCY":  "5",
			"EVENTBRIDGE_ACCESS_KEY_ID":     "test",
			"EVENTBRIDGE_SECRET_ACCESS_KEY": "test",
		}}

	cfg := eventbridgeConfigFromEnv(c, logger, "")
	expCfg := &eventbridge.Config{
		ConnRetryDuration: 5,
		EventBus:          "Gofr",
		EventSource:       "Gofr-application",
		Region:            "us-east-1",
		AccessKeyID:       "test",
		SecretAccessKey:   "test",
	}

	assert.Equal(t, expCfg, cfg)
}

func Test_kvDataConfigFromEnv(t *testing.T) {
	mockCfg1 := &config.MockConfig{
		Data: map[string]string{
			"KV_URL":                "http://localhost:2021",
			"KV_CSP_APP_KEY_FWK":    "test key",
			"KV_CSP_SHARED_KEY_FWK": "test key",
		},
	}

	expConfig1 := kvdata.Config{
		URL:       "http://localhost:2021",
		AppKey:    "test key",
		SharedKey: "test key",
	}

	mockCfg2 := &config.MockConfig{
		Data: map[string]string{
			"KV_URL":            "http://localhost:2021",
			"KV_CSP_APP_KEY":    "test",
			"KV_CSP_SHARED_KEY": "test",
		},
	}

	expConfig2 := kvdata.Config{
		URL:       "http://localhost:2021",
		AppKey:    "test",
		SharedKey: "test",
	}

	testcases := []struct {
		input  *config.MockConfig
		expOut kvdata.Config
	}{
		{mockCfg1, expConfig1},
		{mockCfg2, expConfig2},
	}
	for i, tc := range testcases {
		cfg := kvDataConfigFromEnv(tc.input)
		if !reflect.DeepEqual(cfg, tc.expOut) {
			t.Errorf("Test case failed [%v]. Got: %v,expected:%v", i, cfg, tc.expOut)
		}
	}
}
