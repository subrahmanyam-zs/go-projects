package datastore

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestGetNewDynamoDB(t *testing.T) {
	tcs := []struct {
		region string
		err    error
	}{
		{"", awserr.New("MissingRegion", "could not find region configuration", nil)},
		{"ap-south-1", nil},
	}

	for i, tc := range tcs {
		cfg := DynamoDBConfig{
			Region:            tc.region,
			Endpoint:          "http://localhost:8000",
			AccessKeyID:       "access-key-id",
			SecretAccessKey:   "secret-key",
			ConnRetryDuration: 5,
		}

		_, err := NewDynamoDB(log.NewLogger(), cfg)

		assert.IsType(t, tc.err, err, "TESTCASE[%d], failed.\n", i)
	}
}

func TestHealthCheck(t *testing.T) {
	tcs := []struct {
		accessKey string
		secretKey string
		status    string
	}{
		{"access-key-id", "secret-key", pkg.StatusUp},
		{"", "", pkg.StatusDown},
	}

	for i, tc := range tcs {
		cfg := DynamoDBConfig{
			Region:            "ap-south-1",
			Endpoint:          "http://localhost:8000",
			AccessKeyID:       tc.accessKey,
			SecretAccessKey:   tc.secretKey,
			ConnRetryDuration: 5,
		}
		db, _ := NewDynamoDB(log.NewLogger(), cfg)

		health := db.HealthCheck()
		expHealth := types.Health{Name: pkg.DynamoDB, Status: tc.status}

		assert.Equal(t, expHealth, health, "TEST[%d], failed.\n", i)
	}
}
