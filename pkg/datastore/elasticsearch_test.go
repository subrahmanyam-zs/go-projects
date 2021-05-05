package datastore

import (
	"io"
	"reflect"
	"strconv"
	"testing"

	"github.com/zopsmart/gofr/pkg/gofr/config"
	"github.com/zopsmart/gofr/pkg/gofr/types"
	"github.com/zopsmart/gofr/pkg/log"
)

func TestNewElasticsearchClient(t *testing.T) {
	var elasticSearchCfg ElasticSearchCfg

	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")
	elasticSearchCfg.Host = c.Get("ELASTIC_SEARCH_HOST")
	elasticSearchCfg.Port, _ = strconv.Atoi(c.Get("ELASTIC_SEARCH_PORT"))
	elasticSearchCfg.User = c.Get("ELASTIC_SEARCH_USER")
	elasticSearchCfg.Pass = c.Get("ELASTIC_SEARCH_PASS")
	_, err := NewElasticsearchClient(&elasticSearchCfg)

	if err != nil {
		t.Errorf("Failed.Got %v\tExpected %v", err, nil)
	}
}

func TestNewElasticsearchClientError(t *testing.T) {
	elasticSearchCfg := ElasticSearchCfg{Host: "localhost", Port: 92}
	_, err := NewElasticsearchClient(&elasticSearchCfg)

	if err == nil {
		t.Errorf("Failed.Expected err got nil")
	}
}

func TestDataStore_ElasticsearchHealthCheck(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")

	port, _ := strconv.Atoi(c.Get("ELASTIC_SEARCH_PORT"))
	testCases := []struct {
		c        ElasticSearchCfg
		expected types.Health
	}{
		{
			ElasticSearchCfg{
				Host: c.Get("ELASTIC_SEARCH_HOST"), Port: port, User: c.Get("ELASTIC_SEARCH_USER"), Pass: c.Get("ELASTIC_SEARCH_PASS"),
			},
			types.Health{
				Name: "elasticsearch", Status: "UP", Host: c.Get("ELASTIC_SEARCH_HOST"),
			}},
		{
			ElasticSearchCfg{
				Host: "random", Port: port, User: c.Get("ELASTIC_SEARCH_USER"), Pass: c.Get("ELASTIC_SEARCH_PASS"),
			},
			types.Health{Name: "elasticsearch", Status: "DOWN", Host: "random"},
		},
	}

	for i, tc := range testCases {
		conn, _ := NewElasticsearchClient(&tc.c)
		output := conn.HealthCheck()

		if !reflect.DeepEqual(output, tc.expected) {
			t.Errorf("[FAILED]%v expected: %v, got: %v", i, tc.expected, output)
		}
	}
}
