package eventhub

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/Azure/azure-amqp-common-go/v3/aad"
	"github.com/Azure/azure-amqp-common-go/v3/sas"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/prometheus/client_golang/prometheus"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/avro"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type Config struct {
	Namespace        string
	EventhubName     string
	ClientID         string
	ClientSecret     string
	TenantID         string
	SharedAccessName string
	SharedAccessKey  string
	// Offsets is slice of Offset in which "PartitionID" and "StartOffset"
	// are the field needed to be set to start consuming from specific offset
	Offsets           []Offset
	ConnRetryDuration int
}

type AvroWithEventhubConfig struct {
	EventhubConfig Config
	AvroConfig     avro.Config
}

type Eventhub struct {
	Config
	hub                *eventhub.Hub
	partitionOffsetMap map[string]string // for persisting offsets
	initialiseOffset   sync.Once
}

type Offset struct {
	PartitionID string
	StartOffset string
}

//nolint // The declared global variable can be accessed across multiple functions
var (
	subscribeRecieveCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_pubsub_receive_count",
		Help: "Total number of subscribe operation",
	}, []string{"topic", "consumerGroup"})

	subscribeSuccessCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_pubsub_success_count",
		Help: "Total number of successful subscribe operation",
	}, []string{"topic", "consumerGroup"})

	subscribeFailureCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_pubsub_failure_count",
		Help: "Total number of failed subscribe operation",
	}, []string{"topic", "consumerGroup"})

	publishSuccessCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_pubsub_publish_success_count",
		Help: "Counter for the number of messages successfully published",
	}, []string{"topic", "consumerGroup"})

	publishFailureCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_pubsub_publish_failure_count",
		Help: "Counter for the number of failed publish operations",
	}, []string{"topic", "consumerGroup"})

	publishTotalCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_pubsub_publish_total_count",
		Help: "Counter for the total number of publish operations",
	}, []string{"topic", "consumerGroup"})
)

func New(c *Config) (pubsub.PublisherSubscriber, error) {
	_ = prometheus.Register(subscribeRecieveCount)
	_ = prometheus.Register(subscribeSuccessCount)
	_ = prometheus.Register(subscribeFailureCount)
	_ = prometheus.Register(publishSuccessCount)
	_ = prometheus.Register(publishFailureCount)
	_ = prometheus.Register(publishTotalCount)

	if c.SharedAccessKey != "" && c.SharedAccessName != "" {
		tokenProviderOption := sas.TokenProviderWithKey(c.SharedAccessName, c.SharedAccessKey)

		tokenProvider, err := sas.NewTokenProvider(tokenProviderOption)
		if err != nil {
			return &Eventhub{Config: *c}, err
		}

		hub, err := eventhub.NewHub(c.Namespace, c.EventhubName, tokenProvider)
		if err != nil {
			return &Eventhub{Config: *c}, err
		}

		return &Eventhub{hub: hub, Config: *c, partitionOffsetMap: make(map[string]string)}, nil
	}

	jwtProvider, err := aad.NewJWTProvider(jwtProvider(c))
	if err != nil {
		return &Eventhub{Config: *c}, err
	}

	hub, err := eventhub.NewHub(c.Namespace, c.EventhubName, jwtProvider)
	if err != nil {
		return &Eventhub{Config: *c}, err
	}

	return &Eventhub{hub: hub, Config: *c, partitionOffsetMap: make(map[string]string)}, nil
}

func jwtProvider(c *Config) aad.JWTProviderOption {
	return func(config *aad.TokenProviderConfiguration) error {
		config.TenantID = c.TenantID
		config.ClientID = c.ClientID
		config.ClientSecret = c.ClientSecret
		config.Env = &azure.PublicCloud

		return nil
	}
}

func (e *Eventhub) PublishEvent(key string, value interface{}, headers map[string]string) (err error) {
	return e.PublishEventWithOptions(key, value, headers, nil)
}

func (e *Eventhub) PublishEventWithOptions(key string, value interface{}, headers map[string]string,
	options *pubsub.PublishOptions) (err error) {
	publishTotalCount.WithLabelValues(e.EventhubName, "").Inc()

	data, ok := value.([]byte)
	if !ok {
		data, err = json.Marshal(value)
		if err != nil {
			publishFailureCount.WithLabelValues(e.EventhubName, "").Inc()

			return err
		}
	}

	event := eventhub.NewEvent(data)

	err = e.hub.Send(context.TODO(), event, eventhub.SendWithMessageID(key))
	if err != nil {
		publishFailureCount.WithLabelValues(e.EventhubName, "").Inc()

		return err
	}

	publishSuccessCount.WithLabelValues(e.EventhubName, "").Inc()

	return nil
}

func (e *Eventhub) Subscribe() (*pubsub.Message, error) {
	// for every subscribe
	subscribeRecieveCount.WithLabelValues(e.EventhubName, "").Inc()

	msg := make(chan *pubsub.Message)

	handler := func(ctx context.Context, event *eventhub.Event) error {
		var partition int

		if event.SystemProperties.PartitionID != nil {
			partition = int(*event.SystemProperties.PartitionID)
		}

		msg <- &pubsub.Message{
			Value:     string(event.Data),
			Partition: partition,
			Offset:    *event.SystemProperties.Offset,
			Topic:     e.EventhubName,
		}

		e.partitionOffsetMap[strconv.Itoa(partition)] = strconv.Itoa(int(*event.SystemProperties.Offset))

		return nil
	}

	ctx := context.TODO()

	runtimeInfo, err := e.hub.GetRuntimeInformation(ctx)
	if err != nil {
		// for failed subscribe
		subscribeFailureCount.WithLabelValues(e.EventhubName, "").Inc()
		return nil, err
	}

	// Set the initial offset value for subscribe
	if e.Offsets != nil {
		e.initialiseOffset.Do(func() {
			for _, offset := range e.Offsets {
				e.partitionOffsetMap[offset.PartitionID] = offset.StartOffset
			}
		})
	}

	for _, partitionID := range runtimeInfo.PartitionIDs {
		offset := e.partitionOffsetMap[partitionID]

		_, err := e.hub.Receive(ctx, partitionID, handler, eventhub.ReceiveWithStartingOffset(offset))
		if err != nil {
			// for failed subscribe
			subscribeFailureCount.WithLabelValues(e.EventhubName, "").Inc()
			return nil, err
		}
	}
	// for successful subscribe
	subscribeSuccessCount.WithLabelValues(e.EventhubName, "").Inc()

	return <-msg, nil
}

func (e *Eventhub) SubscribeWithCommit(f pubsub.CommitFunc) (*pubsub.Message, error) {
	return e.Subscribe()
}

func (e *Eventhub) Bind(message []byte, target interface{}) error {
	return json.Unmarshal(message, target)
}

func (e *Eventhub) Ping() error {
	_, err := e.hub.GetRuntimeInformation(context.TODO())
	return err
}

func (e *Eventhub) HealthCheck() types.Health {
	// handling nil object
	if e == nil {
		return types.Health{
			Name:   pkg.EventHub,
			Status: pkg.StatusDown,
		}
	}

	resp := types.Health{
		Name:     pkg.EventHub,
		Status:   pkg.StatusDown,
		Host:     e.Namespace,
		Database: e.EventhubName,
	}

	// configs is present but not connected
	if e.hub == nil {
		return resp
	}

	if err := e.Ping(); err != nil {
		return resp
	}

	resp.Status = pkg.StatusUp

	return resp
}

func (e *Eventhub) CommitOffset(offsets pubsub.TopicPartition) {
}

func (e *Eventhub) IsSet() bool {
	if e == nil {
		return false
	}

	return e.hub != nil
}

// NewEventHubWithAvro initialize EventHub with Avro when EventHubConfig and AvroConfig are right
//nolint:interfacer //`logger` can be `github.com/stretchr/testify/assert.TestingT`
func NewEventHubWithAvro(config *AvroWithEventhubConfig, logger log.Logger) (pubsub.PublisherSubscriber, error) {
	eventHub, err := New(&config.EventhubConfig)
	if err != nil {
		logger.Errorf("Eventhub cannot be initialized, err: %v", err)
		return nil, err
	}

	p, err := avro.NewWithConfig(&config.AvroConfig, eventHub)
	if err != nil {
		logger.Errorf("Avro cannot be initialized, err: %v", err)
		return nil, err
	}

	return p, nil
}
