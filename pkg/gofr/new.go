package gofr

import (
	"context"
	"os"
	"strconv"
	"strings"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/avro"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/eventhub"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/kafka"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"developer.zopsmart.com/go/gofr/pkg/log"
	awssns "developer.zopsmart.com/go/gofr/pkg/notifier/aws-sns"

	"github.com/prometheus/client_golang/prometheus"

	"go.opentelemetry.io/otel"
)

func New() (k *Gofr) {
	var (
		logger       = log.NewLogger()
		configFolder string
	)

	if _, err := os.Stat("./configs"); err == nil {
		configFolder = "./configs"
	} else if _, err := os.Stat("../configs"); err == nil {
		configFolder = "../configs"
	} else {
		configFolder = "../../configs"
	}

	return NewWithConfig(config.NewGoDotEnvProvider(logger, configFolder))
}

//nolint:gocognit  // It's a sequence of initialization. Easier to understand this way.
func NewWithConfig(c Config) (k *Gofr) {
	// Here we do things based on what is provided by Config
	logger := log.NewLogger()

	gofr := &Gofr{
		Logger:         logger,
		Config:         c,
		DatabaseHealth: []HealthCheck{},
	}

	gofr.DataStore.Logger = logger

	appVers := c.Get("APP_VERSION")
	if appVers == "" {
		appVers = pkg.DefaultAppVersion

		logger.Warnf("APP_VERSION is not set. '%v' will be used in logs", pkg.DefaultAppVersion)
	}

	appName := c.Get("APP_NAME")
	if appName == "" {
		appName = pkg.DefaultAppName

		logger.Warnf("APP_NAME is not set.'%v' will be used in logs", pkg.DefaultAppName)
	}

	frameworkInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "zs_info",
		Help: "Gauge to count the pods running for each service and framework version",
	}, []string{"app", "framework"})
	_ = prometheus.Register(frameworkInfo)
	frameworkInfo.WithLabelValues(appName+"-"+appVers, "gofr-"+log.GofrVersion).Set(1)

	s := NewServer(c, gofr)
	gofr.Server = s

	// HTTP PORT
	p, err := strconv.Atoi(c.Get("HTTP_PORT"))
	s.HTTP.Port = p

	if err != nil || p <= 0 {
		s.HTTP.Port = 8000
	}

	// HTTPS Initialisation
	s.HTTPS.KeyFile = c.Get("KEY_FILE")
	s.HTTPS.CertificateFile = c.Get("CERTIFICATE_FILE")

	p, err = strconv.Atoi(c.Get("HTTPS_PORT"))
	s.HTTPS.Port = p

	if err != nil || p <= 0 {
		s.HTTPS.Port = 443
	}

	s.GRPC.server = NewGRPCServer()

	// set GRPC port from config
	p, err = strconv.Atoi(c.Get("GRPC_PORT"))
	if err == nil {
		s.GRPC.Port = p
	}

	// Set Metrics Port
	s.initializeMetricServerConfig(c)

	// If Tracing is set, Set tracing
	_ = enableTracing(c, logger)
	initializeDataStores(c, gofr)

	initializeNotifiers(c, gofr)

	return gofr
}

func (s *server) initializeMetricServerConfig(c Config) {
	// Set Metrics Port
	if val, err := strconv.Atoi(c.Get("METRICS_PORT")); err == nil && val >= 0 {
		s.MetricsPort = val
	}

	if route := c.Get("METRICS_ROUTE"); route != "" {
		s.MetricsRoute = "/" + strings.TrimPrefix(route, "/")
	}
}

func initializePubSub(c Config, k *Gofr) {
	pubsubBackend := c.Get("PUBSUB_BACKEND")
	if pubsubBackend == "" {
		return
	}

	switch pubsubBackend {
	case "KAFKA", "AVRO":
		initializeKafka(c, k)
	case "EVENTHUB":
		initializeEventhub(c, k)
	}
}

// initializeAvro initializes avro schema registry along with
// pubsub present in k.Pubsub, only if registryURL is present,
// else k.PubSub remains as is, either Kafka/Eventhub
func initializeAvro(c *avro.Config, k *Gofr) {
	pubsubKafka, _ := k.PubSub.(*kafka.Kafka)
	pubsubEventhub, _ := k.PubSub.(*eventhub.Eventhub)

	if pubsubKafka == nil && pubsubEventhub == nil {
		k.Logger.Error("Kafka/Eventhub not present, cannot use Avro")
		return
	}

	if c == nil {
		return
	}

	if c.URL == "" {
		k.Logger.Error("Schema registry URL is required for Avro")
	}

	ps, err := avro.NewWithConfig(c, k.PubSub)
	if err != nil {
		k.Logger.Errorf("Avro could not be initialized! SchemaRegistry: %v SchemaVersion: %v, Subject: %v, Error: %v",
			c.URL, c.Version, c.Subject, err)
	}

	if ps != nil {
		k.PubSub = ps
		k.Logger.Infof("Avro initialized! SchemaRegistry: %v SchemaVersion: %v, Subject: %v",
			c.URL, c.Version, c.Subject)
	}
}

func NewCMD() *Gofr {
	var (
		configFolder string
		err          error
	)

	if _, err = os.Stat("./configs"); err == nil {
		configFolder = "./configs"
	} else if _, err = os.Stat("../configs"); err == nil {
		configFolder = "../configs"
	} else {
		configFolder = "../../configs"
	}

	c := config.NewGoDotEnvProvider(log.NewLogger(), configFolder)
	// Here we do things based on what is provided by Config, eg LOG_LEVEL etc.
	logger := log.NewLogger()
	cmdApp := &cmdApp{Router: NewCMDRouter(), metricSvr: &metricServer{route: defaultMetricsRoute}}
	gofr := &Gofr{
		Logger: logger,
		cmd:    cmdApp,
		Config: c, // need to be set for using gofr.Config.Get() method
	}

	if appVers := c.Get("APP_VERSION"); appVers == "" {
		logger.Warnf("APP_VERSION is not set. '%v' will be used in logs", pkg.DefaultAppVersion)
	}

	if appName := c.Get("APP_NAME"); appName == "" {
		logger.Warnf("APP_NAME is not set.'%v' will be used in logs", pkg.DefaultAppName)
	}

	if cmdApp.metricSvr.port, err = strconv.Atoi(c.Get("METRIC_PORT")); err != nil {
		cmdApp.metricSvr.port = defaultMetricsPort
	}

	if route := c.Get("METRIC_ROUTE"); route != "" {
		route = strings.TrimPrefix(route, "/")
		cmdApp.metricSvr.route = "/" + route
	}

	cmdApp.context = NewContext(&responder.CMD{}, request.NewCMDRequest(), gofr)

	tracer := otel.Tracer("gofr")

	// Start tracing span
	ctx, span := tracer.Start(context.Background(), "CMD")

	cmdApp.context.Context = ctx
	cmdApp.tracingSpan = &span

	// If Tracing is set, Set tracing
	_ = enableTracing(c, logger)
	initializeDataStores(c, gofr)

	initializeNotifiers(c, gofr)

	return gofr
}

func enableTracing(c Config, logger log.Logger) error {
	// If Tracing is set, initialize tracing
	tp := TraceProvider(
		c.GetOrDefault("APP_NAME", "gofr"),
		c.Get("TRACER_EXPORTER"),
		c.Get("TRACER_NAME"),
		c.Get("TRACER_PORT"),
		logger,
	)
	if tp == nil {
		return errors.Error("could not create exporter")
	}

	otel.SetTracerProvider(tp)

	return nil
}

// initializeDataStores initializes the Gofr struct with all the data stores for which
// correct config is set in the environment
func initializeDataStores(c Config, k *Gofr) {
	// Redis
	initializeRedis(c, k)

	// DB
	initializeDB(c, k)

	// Cassandra
	initializeCassandra(c, k)

	// Mongo DB
	initializeMongoDB(c, k)

	// PubSub
	initializePubSub(c, k)

	// Elasticsearch
	initializeElasticsearch(c, k)

	// Solr
	initializeSolr(c, k)
}

// initializeRedis initializes the Redis client in the Gofr struct if the Redis configuration is set
// in the environment, in case of an error, it logs the error
func initializeRedis(c Config, k *Gofr) {
	rc := datastore.RedisConfig{
		HostName:                c.Get("REDIS_HOST"),
		Port:                    c.Get("REDIS_PORT"),
		ConnectionRetryDuration: getRetryDuration(c.Get("REDIS_CONN_RETRY")),
	}

	if rc.HostName != "" || rc.Port != "" {
		var err error

		k.Redis, err = datastore.NewRedis(k.Logger, rc)
		k.DatabaseHealth = append(k.DatabaseHealth, k.RedisHealthCheck)

		if err != nil {
			k.Logger.Errorf("could not connect to Redis, HostName: %s, Port: %s, error: %v\n",
				rc.HostName, rc.Port, err)

			go redisRetry(&rc, k)

			return
		}

		k.Logger.Infof("Redis connected. HostName: %s, Port: %s", rc.HostName, rc.Port)
	}
}

// initializeDB initializes the ORM object in the Gofr struct if the DB configuration is set
// in the environment, in case of an error, it logs the error
// nolint:gocognit // breaking down function will reduce readability
func initializeDB(c Config, k *Gofr) {
	dc := datastore.DBConfig{
		HostName:          c.Get("DB_HOST"),
		Username:          c.Get("DB_USER"),
		Password:          c.Get("DB_PASSWORD"),
		Database:          c.Get("DB_NAME"),
		Port:              c.Get("DB_PORT"),
		Dialect:           c.Get("DB_DIALECT"),
		SSL:               c.Get("DB_SSL"),
		ORM:               c.Get("DB_ORM"),
		CertificateFile:   c.Get("DB_CERTIFICATE_FILE"),
		KeyFile:           c.Get("DB_KEY_FILE"),
		ConnRetryDuration: getRetryDuration(c.Get("DB_CONN_RETRY")),
	}

	if dc.HostName != "" && dc.Port != "" {
		if strings.EqualFold(dc.ORM, "SQLX") {
			db, err := datastore.NewSQLX(&dc)
			k.SetORM(db)

			k.DatabaseHealth = append(k.DatabaseHealth, k.SQLXHealthCheck)

			if err != nil {
				k.Logger.Errorf("could not connect to DB, HOST: %s, PORT: %s, Dialect: %s, error: %v\n",
					dc.HostName, dc.Port, dc.Dialect, err)

				go sqlxRetry(&dc, k)

				return
			}

			k.Logger.Infof("DB connected, HostName: %s, Port: %s, Database: %s", dc.HostName, dc.Port, dc.Database)

			return
		}

		db, err := datastore.NewORM(&dc)
		k.SetORM(db)

		k.DatabaseHealth = append(k.DatabaseHealth, k.SQLHealthCheck)

		if err != nil {
			k.Logger.Errorf("could not connect to DB, HostName: %s, Port: %s, Dialect: %s, error: %v\n",
				dc.HostName, dc.Port, dc.Dialect, err)

			go ormRetry(&dc, k)

			return
		}

		k.Logger.Infof("DB connected, HostName: %s, Port: %s, Database: %s", dc.HostName, dc.Port, dc.Database)
	}
}

func initializeMongoDB(c Config, k *Gofr) {
	hostName := c.Get("MONGO_DB_HOST")
	port := c.Get("MONGO_DB_PORT")

	if hostName != "" && port != "" {
		mongoConfig := mongoDBConfigFromEnv(c)

		var err error

		k.MongoDB, err = datastore.GetNewMongoDB(k.Logger, mongoConfig)
		k.DatabaseHealth = append(k.DatabaseHealth, k.MongoHealthCheck)

		if err != nil {
			go mongoRetry(mongoConfig, k)
		}
	}
}

func initializeKafka(c Config, k *Gofr) {
	hosts := c.Get("KAFKA_HOSTS")
	topic := c.Get("KAFKA_TOPIC")

	if hosts != "" && topic != "" {
		var err error

		kafkaConfig := kafkaConfigFromEnv(c)
		avroConfig := avroConfigFromEnv(c)

		k.PubSub, err = kafka.New(kafkaConfig, k.Logger)
		k.DatabaseHealth = append(k.DatabaseHealth, k.PubSubHealthCheck)

		if err != nil {
			k.Logger.Errorf("Kafka could not be initialized, Hosts: %v, Topic: %v, error: %v\n",
				hosts, topic, err)

			go kafkaRetry(kafkaConfig, avroConfig, k)

			return
		}

		k.Logger.Infof("Kafka initialized. Hosts: %v, Topic: %v\n", hosts, topic)

		// initialize Avro using Kafka pubsub if the schema url is specified
		if avroConfig.URL != "" {
			initializeAvro(avroConfig, k)
		}
	}
}

func initializeEventhub(c Config, k *Gofr) {
	hosts := c.Get("EVENTHUB_NAMESPACE")
	topic := c.Get("EVENTHUB_NAME")

	if hosts != "" && topic != "" {
		var err error

		avroConfig := avroConfigFromEnv(c)
		eventhubConfig := eventhubConfigFromEnv(c)

		k.PubSub, err = eventhub.New(&eventhubConfig)
		k.DatabaseHealth = append(k.DatabaseHealth, k.PubSubHealthCheck)

		if err != nil {
			k.Logger.Errorf("Azure Eventhub could not be initialized, Namespace: %v, Eventhub: %v, error: %v\n",
				hosts, topic, err)

			go eventhubRetry(&eventhubConfig, avroConfig, k)

			return
		}

		k.Logger.Infof("Azure Eventhub initialized, Namespace: %v, Eventhub: %v\n", hosts, topic)

		// initialize Avro using eventhub pubsub if the schema url is specified
		if avroConfig.URL != "" {
			initializeAvro(avroConfig, k)
		}
	}
}

// initializeCassandra initializes the Cassandra/ YCQL client in the Gofr struct if the Cassandra configuration is set
// in the environment, in case of an error, it logs the error
// nolint:gocognit // reducing the function further is not required
func initializeCassandra(c Config, k *Gofr) {
	validDialects := map[string]bool{
		"cassandra": true,
		"ycql":      true,
	}

	host := c.Get("CASS_DB_HOST")
	port := c.Get("CASS_DB_PORT")
	dialect := strings.ToLower(c.Get("CASS_DB_DIALECT"))

	if host == "" || port == "" {
		return
	}

	if dialect == "" {
		dialect = "cassandra"
	}

	// Checks if dialect is valid
	if _, ok := validDialects[dialect]; !ok {
		k.Logger.Errorf("invalid dialect: supported dialects are - cassandra, ycql")
		return
	}

	var err error

	switch dialect {
	case "ycql":
		ycqlconfig := getYcqlConfigs(c)

		k.YCQL, err = datastore.GetNewYCQL(k.Logger, &ycqlconfig)
		k.DatabaseHealth = append(k.DatabaseHealth, k.YCQLHealthCheck)

		if err != nil {
			go yclRetry(&ycqlconfig, k)

			return
		}

	default:
		cassandraCfg := cassandraConfigFromEnv(c)

		k.Cassandra, err = datastore.GetNewCassandra(k.Logger, cassandraCfg)
		k.DatabaseHealth = append(k.DatabaseHealth, k.CQLHealthCheck)

		if err != nil {
			k.Logger.Errorf("could not connect to Cassandra, Hosts: %s, Port: %s, Error: %v\n", host, port, err)

			go cassandraRetry(cassandraCfg, k)

			return
		}
	}
}

func getYcqlConfigs(c Config) datastore.CassandraCfg {
	timeout, err := strconv.Atoi(c.Get("CASS_DB_TIMEOUT"))
	if err != nil {
		// setting default timeout of 600 milliseconds
		timeout = 600
	}

	cassandraConnTimeout, err := strconv.Atoi(c.Get("CASS_DB_CONN_TIMEOUT"))
	if err != nil {
		// setting default timeout of 600 milliseconds
		cassandraConnTimeout = 600
	}

	port, err := strconv.Atoi(c.Get("CASS_DB_PORT"))
	if err != nil || port == 0 {
		// if any error, setting default
		port = 9042
	}

	return datastore.CassandraCfg{
		Hosts:               c.Get("CASS_DB_HOST"),
		Port:                port,
		Username:            c.Get("CASS_DB_USER"),
		Password:            c.Get("CASS_DB_PASS"),
		Keyspace:            c.Get("CASS_DB_KEYSPACE"),
		Timeout:             timeout,
		ConnectTimeout:      cassandraConnTimeout,
		ConnRetryDuration:   getRetryDuration(c.Get("CASS_CONN_RETRY")),
		CertificateFile:     c.Get("CASS_DB_CERTIFICATE_FILE"),
		KeyFile:             c.Get("CASS_DB_KEY_FILE"),
		RootCertificateFile: c.Get("CASS_DB_ROOT_CERTIFICATE_FILE"),
		HostVerification:    getBool(c.Get("CASS_DB_HOST_VERIFICATION")),
		InsecureSkipVerify:  getBool(c.Get("CASS_DB_INSECURE_SKIP_VERIFY")),
		DataCenter:          c.Get("DATA_CENTER"),
	}
}

func initializeElasticsearch(c Config, k *Gofr) {
	elasticSearchCfg := getElasticSearchConfigFromEnv(c)

	if elasticSearchCfg.Host == "" || elasticSearchCfg.Port == 0 {
		return
	}

	var err error

	k.Elasticsearch, err = datastore.NewElasticsearchClient(&elasticSearchCfg)
	k.DatabaseHealth = append(k.DatabaseHealth, k.ElasticsearchHealthCheck)

	if err != nil {
		k.Logger.Errorf("could not connect to Elasticsearch, HOST: %s, PORT: %v, Error: %v\n", elasticSearchCfg.Host, elasticSearchCfg.Port, err)

		go elasticSearchRetry(&elasticSearchCfg, k)

		return
	}

	k.Logger.Infof("connected to Elasticsearch, HOST: %s, PORT: %v\n", elasticSearchCfg.Host, elasticSearchCfg.Port)
}

func initializeSolr(c Config, k *Gofr) {
	host := c.Get("SOLR_HOST")
	port := c.Get("SOLR_PORT")

	if host == "" || port == "" {
		return
	}

	k.Solr = datastore.NewSolrClient(host, port)
	k.Logger.Infof("Solr connected. Host: %s, Port: %s \n", host, port)
}

func initializeNotifiers(c Config, k *Gofr) {
	notifierBackend := c.Get("NOTIFIER_BACKEND")
	if notifierBackend == "" {
		return
	}

	switch notifierBackend {
	case "SNS":
		initializeAwsSNS(c, k)
	}
}
func initializeAwsSNS(c Config, k *Gofr) {

	awsConfig := awsSNSConfigFromEnv(c)
	var err error
	k.Notifier, err = awssns.New(&awsConfig)
	k.DatabaseHealth = append(k.DatabaseHealth, k.Notifier.HealthCheck)
	if err != nil {
		k.Logger.Errorf("AWS SNS could not be initialized, error: %v\n", err)

		go awsSNSRetry(&awsConfig, k)

		return
	}

	k.Logger.Infof("AWS SNS initialized")
}
