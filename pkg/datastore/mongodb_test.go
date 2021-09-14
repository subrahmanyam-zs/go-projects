package datastore

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestGetNewMongoDB_ContextErr(t *testing.T) {
	mongoConfig := MongoConfig{"fake_host", "9999", "admin", "admin123", "test", false, false, 30}
	expErr := context.DeadlineExceeded

	_, err := GetNewMongoDB(log.NewLogger(), &mongoConfig)
	if err != nil && !strings.Contains(err.Error(), expErr.Error()) {
		t.Errorf("Error in testcase. Expected: %v, Got: %v", expErr, err)
	}
}

func TestGetNewMongoDB_ConnectionError(t *testing.T) {
	mongoConfig := MongoConfig{"", "", "", "", "test", false, false, 30}
	expErr := errors.New("error validating uri: username required if URI contains user info")

	_, err := GetNewMongoDB(log.NewLogger(), &mongoConfig)
	if err != nil && !strings.Contains(err.Error(), expErr.Error()) {
		t.Errorf("Error in testcase. Expected: %v, Got: %v", expErr, err)
	}
}

func TestGetMongoDBFromEnv_Success(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	_ = config.NewGoDotEnvProvider(logger, "../../configs")

	// Checking for connection with default env vars
	_, err := GetMongoDBFromEnv(logger)
	if err != nil {
		t.Error(err)
	}
}

func TestGetMongoDBFromEnv_Error(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	testcases := []struct {
		envKey    string
		newEnvVal string
		expErr    error
	}{
		{"MONGO_DB_HOST", "fake_host", context.DeadlineExceeded},
		{"MONGO_DB_ENABLE_SSL", "true", context.DeadlineExceeded},
		{"MONGO_DB_ENABLE_SSL", "non_bool", &strconv.NumError{
			Func: "ParseBool",
			Num:  "non_bool",
			Err:  errors.New("invalid syntax"),
		}},
	}

	for i := range testcases {
		oldEnvVal := c.Get(testcases[i].envKey)

		t.Setenv(testcases[i].envKey, testcases[i].newEnvVal)

		// Checking for connection with default env vars
		_, err := GetMongoDBFromEnv(logger)
		if err != nil && !strings.Contains(err.Error(), testcases[i].expErr.Error()) {
			t.Errorf("Expected %v but got %v", testcases[i].expErr, err)
		}

		t.Setenv(testcases[i].envKey, oldEnvVal)
	}
}

func TestGetMongoConfigFromEnv_SSL_RetryWrites(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")
	oldEnableSSL := c.Get("MONGO_DB_ENABLE_SSL")
	oldretryWrites := c.Get("MONGO_DB_RETRY_WRITES")

	testcases := []struct {
		enableSSL   string
		retryWrites string
		expErr      error
	}{
		{"non_bool", "false", &strconv.NumError{
			Func: "ParseBool",
			Num:  "non_bool",
			Err:  errors.New("invalid syntax"),
		}},
		{"false", "non_bool", &strconv.NumError{
			Func: "ParseBool",
			Num:  "non_bool",
			Err:  errors.New("invalid syntax"),
		}},
	}

	for i := range testcases {
		t.Setenv("MONGO_DB_ENABLE_SSL", testcases[i].enableSSL)

		t.Setenv("MONGO_DB_RETRY_WRITES", testcases[i].retryWrites)

		_, err := getMongoConfigFromEnv()
		if !reflect.DeepEqual(err, testcases[i].expErr) {
			t.Errorf("Expected: %v, Got:%v", testcases[i].expErr, err)
		}
	}

	t.Setenv("MONGO_DB_ENABLE_SSL", oldEnableSSL)

	t.Setenv("MONGO_DB_RETRY_WRITES", oldretryWrites)
}

func Test_getMongoConnectionString(t *testing.T) {
	testcases := []struct {
		config        MongoConfig
		expConnString string
	}{
		{
			MongoConfig{"any_host", "9999", "admin", "admin123", "test", true, false, 30},
			"mongodb://admin:admin123@any_host:9999/?ssl=true&retrywrites=false",
		},
		{
			MongoConfig{"", "", "", "", "test", false, true, 30},
			"mongodb://:@:/?ssl=false&retrywrites=true",
		},
	}

	for i := range testcases {
		connStr := getMongoConnectionString(&testcases[i].config)
		if connStr != testcases[i].expConnString {
			t.Errorf("Testcase[%v] failed. Expected: %v, \nGot: %v", i, testcases[i].expConnString, connStr)
		}
	}
}

func TestDataStore_HealthCheck(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	testCases := []struct {
		config   MongoConfig
		expected types.Health
	}{
		{MongoConfig{HostName: c.Get("MONGO_DB_HOST"), Port: c.Get("MONGO_DB_PORT"),
			Username: c.Get("MONGO_DB_USER"), Password: c.Get("MONGO_DB_PASS"), Database: c.Get("MONGO_DB_NAME"),
		},
			types.Health{Name: "mongo", Status: "UP", Host: c.Get("MONGO_DB_HOST"), Database: "test"}},
		{MongoConfig{HostName: "random", Port: c.Get("MONGO_DB_PORT"), Username: c.Get("MONGO_DB_USER"),
			Password: c.Get("MONGO_DB_PASS"), Database: c.Get("MONGO_DB_NAME")},
			types.Health{Name: pkg.Mongo, Status: pkg.StatusDown, Host: "random", Database: "test"},
		},
	}

	for i, tc := range testCases {
		conn, _ := GetNewMongoDB(logger, &tc.config)

		output := conn.HealthCheck()
		if !reflect.DeepEqual(output, tc.expected) {
			t.Errorf("[FAILED]%v,Got %v,expecetd %v", i, output, tc.expected)
		}
	}
}

// TestDataStore_HealthCheck_Down tests the health check response when db was connected but goes down
func TestDataStore_HealthCheck_Down(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	conf := &MongoConfig{HostName: c.Get("MONGO_DB_HOST"), Port: c.Get("MONGO_DB_PORT"),
		Username: c.Get("MONGO_DB_USER"), Password: c.Get("MONGO_DB_PASS"), Database: c.Get("MONGO_DB_NAME"),
	}
	expectedResponse := types.Health{
		Name:     pkg.Mongo,
		Status:   pkg.StatusDown,
		Host:     conf.HostName,
		Database: conf.Database,
	}

	m := mongodb{
		Database: nil,
		config:   conf,
	}

	resp := m.HealthCheck()
	if !reflect.DeepEqual(resp, expectedResponse) {
		t.Errorf("expected %v\tgot %v", expectedResponse, resp)
	}
}
