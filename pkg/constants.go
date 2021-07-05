package pkg

import "developer.zopsmart.com/go/gofr/pkg/log"

const (
	StatusUp          = "UP"
	StatusDown        = "DOWN"
	StatusDegraded    = "DEGRADED"
	Cassandra         = "cassandra"
	Redis             = "redis"
	SQL               = "sql"
	Mongo             = "mongo"
	Kafka             = "kafka"
	ElasticSearch     = "elasticsearch"
	YCQL              = "ycql"
	EventHub          = "eventhub"
	DefaultAppName    = "gofr-app"
	DefaultAppVersion = "dev"
	Framework         = "gofr-" + log.GofrVersion
	AWSSNS            = "aws-sns"
)
