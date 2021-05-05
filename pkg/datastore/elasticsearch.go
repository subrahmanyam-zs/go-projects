package datastore

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v6"
	"github.com/zopsmart/gofr/pkg"
	"github.com/zopsmart/gofr/pkg/gofr/types"
)

type ElasticSearchCfg struct {
	Host                    string
	Port                    int
	User                    string
	Pass                    string
	ConnectionRetryDuration int
}

type Elasticsearch struct {
	url string
	*elastic.Client
	config ElasticSearchCfg
}

func NewElasticsearchClient(c *ElasticSearchCfg) (Elasticsearch, error) {
	var e Elasticsearch

	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	e.url = url
	//  options function list
	optionFunc := []elastic.ClientOptionFunc{
		elastic.SetSniff(false), elastic.SetURL(url)}

	if c.User != "" && c.Pass != "" {
		optionFunc = append(optionFunc, elastic.SetBasicAuth(c.User, c.Pass))
	}

	e.config = *c
	// establish connection with the elastic search host
	client, err := elastic.NewClient(optionFunc...)
	if err != nil {
		return e, err
	}

	e.Client = client

	return e, nil
}

func (e *Elasticsearch) Ping(ctx context.Context) (*elastic.PingResult, int, error) {
	return e.Client.Ping(e.url).Do(ctx)
}

func (e *Elasticsearch) HealthCheck() types.Health {
	resp := types.Health{
		Name:   pkg.ElasticSearch,
		Status: pkg.StatusDown,
		Host:   e.config.Host,
	}
	// The following check is for the condition when the connection to Elasticsearch has not been made during initialization
	if e.Client == nil {
		return resp
	}

	_, _, err := e.Ping(context.Background())
	if err != nil {
		return resp
	}

	resp.Status = pkg.StatusUp

	return resp
}
