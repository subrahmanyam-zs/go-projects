package kafka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/avro"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestNewKafkaProducer(t *testing.T) {
	tests := []struct {
		config *Config
		err    error
	}{
		{
			config: &Config{Brokers: "somehost"},
			err:    sarama.ErrOutOfBrokers,
		},
		{
			config: &Config{
				Topics:  []string{"some-topic"},
				Brokers: "localhost:2009",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		_, err := NewKafkaProducer(test.config)
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("FAILED, expected: %v, got: %v", test.err, err)
		}
	}
}

func TestNewKafkaConsumer(t *testing.T) {
	conf := sarama.NewConfig()
	conf.Consumer.Group.Session.Timeout = 1
	tests := []struct {
		config *Config
		err    error
	}{
		{
			config: &Config{
				Brokers: "localhost:2009",
				Config:  conf,
			},
			err: sarama.ConfigurationError("Consumer.Group.Session.Timeout must be >= 2ms"),
		},
		{
			config: &Config{
				Brokers: "localhost:2009",
				Topics:  []string{"some-topic"},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		_, err := NewKafkaConsumer(test.config)
		if !reflect.DeepEqual(test.err, err) {
			t.Errorf("FAILED, expected: %v, got: %v", test.err, err)
		}
	}
}

func TestNewKafkaFromEnv(t *testing.T) {
	logger := log.NewLogger()
	config.NewGoDotEnvProvider(logger, "../../../../configs")

	{
		// success case
		_, err := NewKafkaFromEnv()
		if err != nil {
			t.Errorf("FAILED, expected: %v, got: %v", nil, err)
		}
		{
			// error case due to invalid kafka host
			kafkaHosts := os.Getenv("KAFKA_HOSTS")
			os.Setenv("KAFKA_HOSTS", "localhost:9999")
			_, err := NewKafkaFromEnv()
			if err == nil {
				t.Errorf("Failed, expected: %v, got: %v ", brokersErr{}, nil)
			}
			os.Setenv("KAFKA_HOSTS", kafkaHosts)
		}
	}
}

func TestNewKafka(t *testing.T) {
	conf := sarama.NewConfig()
	conf.Consumer.Group.Session.Timeout = 1

	logger := log.NewMockLogger(io.Discard)

	testCases := []struct {
		k       Config
		wantErr bool
	}{
		{
			Config{
				Brokers:        "localhost:2008,localhost:2009",
				Topics:         []string{"test-topic"},
				MaxRetry:       4,
				RetryFrequency: 300,
			}, false,
		},
		{
			Config{
				Brokers: "localhost:2009",
				Topics:  []string{"test-topic"},
				Config:  conf,
			}, true,
		},
		{
			Config{
				Brokers: "localhost:0000",
				Topics:  []string{"test-topic"},
			}, true,
		},
	}

	for i, tt := range testCases {
		_, err := New(&tt.k, logger)
		if !tt.wantErr && err != nil {
			t.Errorf("FAILED[%v], expected: %v, got: %v", i+1, tt.wantErr, err)
		}

		if tt.wantErr && err == nil {
			t.Errorf("FAILED[%v], expected: %v, got: %v", i+1, tt.wantErr, err)
		}
	}
}

func Test_PubSub(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../../../configs")

	k, err := New(&Config{
		Brokers:        c.Get("KAFKA_HOSTS"),
		Topics:         []string{c.Get("KAFKA_TOPIC")},
		InitialOffsets: OffsetOldest,
		GroupID:        "testing-consumerGroup",
	}, logger)
	if err != nil {
		t.Errorf("Kafka connection failed : %v", err)
		return
	}

	Ping(t, k)
	PublishEvent(t, k)
	SubscribeWithCommit(t, k, c.Get("KAFKA_TOPIC"))
	Subscribe(t, k)
}

func PublishEvent(t *testing.T, k *Kafka) {
	type args struct {
		key   string
		value interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success writing message", args{"testKey", "testValue"}, false},
		{"success writing message", args{"testKey1", "testValue1"}, false},
		{"success writing message", args{"testKey2", "testValue2"}, false},
		{"success writing message", args{"testKey3", "testValue3"}, false},
		{"json error in message", args{"testKey", make(chan int)}, true},
	}

	for _, tt := range tests {
		tt := tt
		if err := k.PublishEvent(tt.args.key, tt.args.value, map[string]string{
			"header": "value",
		}); (err != nil) != tt.wantErr {
			t.Errorf("PublishEvent() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func SubscribeWithCommit(t *testing.T, k *Kafka, topic string) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)

	count := 0
	commitFunc := func(msg *pubsub.Message) (bool, bool) {
		logger.Infof("Message received: %v, Topic: %v", msg.Value, msg.Topic)

		if count < 1 {
			count++

			return true, true
		}

		return false, false
	}

	_, _ = k.SubscribeWithCommit(commitFunc)

	expectedMsg1 := fmt.Sprintf("Message received: \\\"testValue\\\", Topic: %v", topic)
	expectedMsg2 := fmt.Sprintf("Message received: \\\"testValue1\\\", Topic: %v", topic)

	if !strings.Contains(b.String(), expectedMsg1) {
		t.Errorf("FAILED expected: %v, got: %v", expectedMsg1, b.String())
	}

	if !strings.Contains(b.String(), expectedMsg2) {
		t.Errorf("FAILED expected: %v, got: %v", expectedMsg2, b.String())
	}

	msg, err := k.SubscribeWithCommit(nil)
	if err != nil {
		t.Errorf("FAILED, expected no error when commitFunc is not provided, got: %v", err)
	}

	k.CommitOffset(pubsub.TopicPartition{
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
	})
}

func Subscribe(t *testing.T, k *Kafka) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"success reading message", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			msg, err := k.Subscribe()
			if (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if msg == nil {
				t.Errorf("Subscribe(): expected message, got nil")
				return
			}

			if len(msg.Headers) != 1 {
				t.Errorf("Subscribe() only one message header should be present, found: %v", len(msg.Headers))
			}
		})
	}
}

func Ping(t *testing.T, k *Kafka) {
	err := k.Ping()
	if err != nil && err.Error() == "invalid brokers connection failed" {
		t.Errorf("FAILED, expected: successful ping, got: %v", err)
	}

	k.config.Brokers = "localhost"

	err = k.Ping()
	if err == nil {
		t.Error("FAILED, expected: unsuccessful ping, got: nil")
	}
}

func Test_convertKafkaConfig(t *testing.T) {
	expectedConfig := sarama.NewConfig()
	setDefaultConfig(expectedConfig)

	expectedConfig.Version = sarama.V2_0_0_0
	expectedConfig.Consumer.Group.Member.UserData = []byte("1")
	expectedConfig.Consumer.Offsets.Initial = OffsetOldest

	kafkaConfig := &Config{GroupID: "1", MaxRetry: 3, InitialOffsets: OffsetOldest}

	convertKafkaConfig(kafkaConfig)

	kafkaConfig.Config.Producer.Partitioner = nil
	expectedConfig.Producer.Partitioner = nil

	assert.Equal(t, expectedConfig, kafkaConfig.Config)
}

func TestKafkaHealthCheck(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../../../configs")
	topic := c.Get("KAFKA_TOPIC")
	topics := strings.Split(topic, ",")
	testCases := []struct {
		config   Config
		expected types.Health
	}{
		{Config{Brokers: c.Get("KAFKA_HOSTS"), Topics: topics},
			types.Health{
				Name:     pkg.Kafka,
				Status:   pkg.StatusUp,
				Host:     c.Get("KAFKA_HOSTS"),
				Database: topic}},
		{
			Config{
				Brokers: "random", Topics: topics},
			types.Health{
				Name:     pkg.Kafka,
				Status:   pkg.StatusDown,
				Host:     "random",
				Database: topic}},
	}

	for i, tc := range testCases {
		conn, _ := New(&tc.config, logger)
		output := conn.HealthCheck()

		if !reflect.DeepEqual(tc.expected, output) {
			t.Errorf("[TESTCASE%v]Failed. Got%v Expected%v", i+1, output, tc.expected)
		}
	}
}

func TestKafka_HealthCheckDown(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	{
		// invalid configs
		c := &Config{
			Brokers: "localhost:2003",
			Topics:  []string{"unknown-topic"},
		}

		expected := types.Health{
			Name:     pkg.Kafka,
			Status:   pkg.StatusDown,
			Host:     c.Brokers,
			Database: c.Topics[0],
		}

		con, _ := New(c, logger)
		healthCheck := con.HealthCheck()

		if !reflect.DeepEqual(expected, healthCheck) {
			t.Errorf("Got %v\tExpected %v\n", healthCheck, expected)
		}
	}

	{
		// nil kafka connection
		expected := types.Health{
			Name:   pkg.Kafka,
			Status: pkg.StatusDown,
		}
		var con *Kafka
		healthCheck := con.HealthCheck()

		if !reflect.DeepEqual(expected, healthCheck) {
			t.Errorf("Got %v\tExpected %v\n", healthCheck, expected)
		}
	}

	{
		// nil producer and consumer
		c := &Config{
			Brokers: "localhost:2003",
			Topics:  []string{"unknown-topic"},
		}

		expected := types.Health{
			Name:     pkg.Kafka,
			Status:   pkg.StatusDown,
			Host:     c.Brokers,
			Database: c.Topics[0],
		}

		con := new(Kafka)
		con.config = c

		healthCheck := con.HealthCheck()

		if !reflect.DeepEqual(expected, healthCheck) {
			t.Errorf("Got %v\tExpected %v\n", healthCheck, expected)
		}
	}
}

func TestIsSet(t *testing.T) {
	var k *Kafka

	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../../../configs")
	topic := c.Get("KAFKA_TOPIC")
	conn, _ := New(&Config{Brokers: c.Get("KAFKA_HOSTS"), Topics: strings.Split(topic, ",")}, logger)

	testcases := []struct {
		k    *Kafka
		resp bool
	}{
		{k, false},
		{&Kafka{}, false},
		{&Kafka{Producer: conn.Producer}, false},
		{&Kafka{Consumer: conn.Consumer}, false},
		{conn, true},
	}

	for i, v := range testcases {
		resp := v.k.IsSet()
		if resp != v.resp {
			t.Errorf("[TESTCASE%d]Failed.Expected %v\tGot %v\n", i+1, v.resp, resp)
		}
	}
}

func TestSubscribeError(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../../../configs")
	topic := "dummy-topic"
	conn, _ := New(&Config{Brokers: c.Get("KAFKA_HOSTS"), Topics: strings.Split(topic, ",")}, logger)

	_ = conn.Consumer.ConsumerGroup.Close()

	if _, err := conn.Subscribe(); err == nil {
		t.Errorf("FAILED, expected error from subcribe got nil")
	}
}

func TestSubscribeWithCommitError(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../../../configs")
	topic := "dummy-topic"
	conn, _ := New(&Config{Brokers: c.Get("KAFKA_HOSTS"), Topics: strings.Split(topic, ",")}, logger)

	_ = conn.Consumer.ConsumerGroup.Close()

	if _, err := conn.SubscribeWithCommit(nil); err == nil {
		t.Errorf("FAILED, expected error from subcribe got nil")
	}
}

func PublishMessage(t *testing.T, k *Kafka) {
	tests := []struct {
		key   string
		value interface{}
	}{
		{"testKey", "testValue"},
		{"testKey1", "testValue1"},
		{"testKey2", "testValue2"},
		{"testKey3", "testValue3"},
	}

	for i, tc := range tests {
		if err := k.PublishEvent(tc.key, tc.value, map[string]string{
			"header": "value",
		}); err != nil {
			t.Errorf("Failed[%v] expected error as nil\n got %v", i, err)
		}
	}
}

// Test_PubSubWithOffset check the subscribe operation with custom initial offset value.
func Test_PubSubWithOffset(t *testing.T) {
	topic := "test-custom-offset"
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../../configs")
	// prereqisite
	k, err := New(&Config{
		Brokers:        c.Get("KAFKA_HOSTS"),
		Topics:         []string{topic},
		InitialOffsets: OffsetOldest,
		GroupID:        "testing-consumerGroup",
	}, logger)

	if err != nil {
		t.Errorf("Kafka connection failed : %v", err)
		return
	}
	// In this we are first trying to publish some messages then we are consuming 1 message
	PublishMessage(t, k)
	Subscribe(t, k)
	// testing subscribe operation with offset
	k, err = New(&Config{
		Brokers:        c.Get("KAFKA_HOSTS"),
		Topics:         []string{topic},
		InitialOffsets: OffsetOldest,
		GroupID:        "testing-consumerGroup",
		Offsets:        []pubsub.TopicPartition{{Topic: topic, Partition: 0, Offset: 0}},
	}, logger)

	if err != nil {
		t.Errorf("Kafka connection failed : %v", err)
		return
	}
	// since we have already consumed one message so the offset value changed to 1. As we are setting the offset value as 0
	// so when we will try to consume message it will start consuming message from offset 0.
	// If setting of offset is unsuccessful then the subscribe operation will fail.
	SubscribeWithCommit(t, k, topic)
	Subscribe(t, k)
}

func Test_populateOffsetTopic(t *testing.T) {
	tests := []struct {
		config         *Config
		expectedConfig *Config
	}{
		{&Config{Topics: []string{"test-topic"}}, &Config{Topics: []string{"test-topic"}}},
		{&Config{Offsets: []pubsub.TopicPartition{}}, &Config{Offsets: []pubsub.TopicPartition{}}},
		{&Config{Offsets: []pubsub.TopicPartition{{Offset: 1}}}, &Config{Offsets: []pubsub.TopicPartition{{Offset: 1}}}},
		{&Config{Topics: []string{"test-topic"}, Offsets: []pubsub.TopicPartition{{Offset: 1, Topic: "test-custom-topic"}, {Offset: 2}}},
			&Config{Topics: []string{"test-topic"}, Offsets: []pubsub.TopicPartition{{Offset: 1, Topic: "test-custom-topic"},
				{Offset: 2, Topic: "test-topic"}}}},
	}

	for i, tc := range tests {
		populateOffsetTopic(tc.config)
		assert.Equal(t, tc.expectedConfig, tc.config, i)
	}
}

func TestNewKafkaWithAvro(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respMap := map[string]interface{}{"subject": "gofr-value", "version": 2, "id": 293,
			"schema": `{"type":"record","name":"test","fields":[{"name":"ID","type":"string"}]}`}
		_ = json.NewEncoder(w).Encode(respMap)
	}))

	testCases := []struct {
		config  AvroWithKafkaConfig
		wantErr bool
	}{
		{
			// Success cases, kafkaConfig and avroConfig both are right
			config: AvroWithKafkaConfig{
				KafkaConfig: Config{
					Brokers: "localhost:2008,localhost:2009", Topics: []string{"test-topic"},
				},
				AvroConfig: avro.Config{URL: server.URL, Version: "", Subject: "gofr-value"},
			}, wantErr: false,
		},
		{
			// failure due wrong kafkaConfig, so it wil not check the avroConfig
			config: AvroWithKafkaConfig{
				KafkaConfig: Config{
					Brokers: "localhost:0000", Topics: []string{"test-topic"},
				}, AvroConfig: avro.Config{URL: server.URL, Version: "", Subject: "gofr-value"},
			}, wantErr: true,
		},
		{
			// failure due to wrong avroConfig
			config: AvroWithKafkaConfig{
				KafkaConfig: Config{
					Brokers: "localhost:2008", Topics: []string{"test-topic"},
				}, AvroConfig: avro.Config{URL: "dummy-url.com", Subject: "gofr-value"},
			}, wantErr: true,
		},
	}

	for i, tc := range testCases {
		_, err := NewKafkaWithAvro(&tc.config, logger)
		if !tc.wantErr && err != nil {
			t.Errorf("FAILED[%v], expected: %v, got: %v", i+1, tc.wantErr, true)
		}

		if tc.wantErr && err == nil {
			t.Errorf("FAILED[%v], expected: %v, got: %v", i+1, tc.wantErr, false)
		}
	}
}

func Test_Printf(t *testing.T) {
	tests := []struct {
		format    string
		input     []interface{}
		expOutput string
	}{
		{"%s %s %s", []interface{}{"log", struct{ Name string }{"data"}, map[string]interface{}{"key": "value"}},
			"log {data} map[key:value]"},
		{"print data %v %s", []interface{}{123, map[string]interface{}{"key": "value"}},
			"print data 123 map[key:value]"},
	}

	for i, tc := range tests {
		b := new(bytes.Buffer)
		kl := kafkaLogger{logger: log.NewMockLogger(b)}

		kl.Printf(tc.format, tc.input...)

		if !strings.Contains(b.String(), tc.expOutput) {
			t.Errorf("failed[%v] expected %v\n got %v", i, tc.expOutput, b.String())
		}
	}
}

func Test_Print(t *testing.T) {
	input := []interface{}{"Print the sys log,", "Print kafka Log", map[string]interface{}{"key": "value"}}
	expOutput := "Print the sys log, Print kafka Log"

	b := new(bytes.Buffer)
	kl := kafkaLogger{logger: log.NewMockLogger(b)}

	kl.Print(input...)

	if !strings.Contains(b.String(), expOutput) {
		t.Errorf("failed expected %v\n got %v", expOutput, b.String())
	}
}

func Test_Println(t *testing.T) {
	input := []interface{}{"Print the sys log,", "Print kafka Log", map[string]interface{}{"key": "value"}}
	expOutput := "Print the sys log, Print kafka Log"

	b := new(bytes.Buffer)
	kl := kafkaLogger{logger: log.NewMockLogger(b)}

	kl.Println(input...)

	if !strings.Contains(b.String(), expOutput) {
		t.Errorf("failed expected %v\n got %v", expOutput, b.String())
	}
}

func TestKafka_SubscribeNilMessage(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../../configs")

	topic := "test-topic"
	conn, _ := New(&Config{Brokers: c.Get("KAFKA_HOSTS"), Topics: strings.Split(topic, ",")}, logger)

	// close the channel to get the msg as nil
	close(conn.Consumer.ConsumerGroupHandler.msg)

	msg, err := conn.subscribeMessage()
	if msg != nil {
		t.Errorf("Failed: Expected Message: %v, Got: %v", nil, msg)
	}

	assert.Equal(t, ErrConsumeMsg, err)
}
