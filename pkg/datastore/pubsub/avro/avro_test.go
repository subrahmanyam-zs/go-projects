package avro

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zopsmart/gofr/pkg/datastore/pubsub"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr/types"
)

type mockPubSub struct {
	Param string
}

func (m *mockPubSub) CommitOffset(offsets pubsub.TopicPartition) {
}

func (m *mockPubSub) PublishEventWithOptions(key string, val interface{}, headers map[string]string, options *pubsub.PublishOptions) error {
	return nil
}

func (m *mockPubSub) PublishEvent(key string, value interface{}, headers map[string]string) error {
	return nil
}

func (m *mockPubSub) Subscribe() (*pubsub.Message, error) {
	if m.Param == "error" {
		return nil, &errors.Response{Reason: "test error"}
	}

	binarySchemaID := []byte(`00000`)
	if m.Param == "id" {
		binarySchemaID = []byte(`00001`)
	}

	return &pubsub.Message{
		SchemaID: 1,
		Topic:    "test_topic",
		Key:      "test",
		Value:    string(binarySchemaID) + `{"name": "test"}`,
		Headers:  map[string]string{"name": "avro-test"},
	}, nil
}

func (m *mockPubSub) SubscribeWithCommit(commitFunc pubsub.CommitFunc) (*pubsub.Message, error) {
	return m.Subscribe()
}

func (m *mockPubSub) Bind(v []byte, target interface{}) error {
	return nil
}

func (m *mockPubSub) Ping() error {
	return nil
}

func (m *mockPubSub) HealthCheck() types.Health {
	return types.Health{}
}

func (m *mockPubSub) IsSet() bool {
	return true
}

type mockSchemaClient struct {
}

func (m *mockSchemaClient) GetSchemaByVersion(subject, version string) (id int, s string, err error) {
	if subject == "error" {
		return 0, "", &errors.Response{Reason: "test error"}
	}

	schema := `{"name": "name", "type": "string"}`

	return 1, schema, nil
}
func (m *mockSchemaClient) GetSchema(id int) (string, error) {
	if id == 808464433 {
		schema := `{"name": "name", "type": "string"}`
		return schema, nil
	}

	return "", &errors.Response{Reason: "test error"}
}

//nolint:gocognit // reducing the cognitive complexity so all the test cases can be considered
func TestAvro_Publish(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}

	tests := []struct {
		name             string
		args             args
		mockPubSub       *mockPubSub
		mockSchemaClient *mockSchemaClient
		subject          string
		pubErr           bool
		wantErr          bool
	}{
		{"error converting native to binary", args{key: "testKey", value: nil}, &mockPubSub{}, &mockSchemaClient{}, "test_topic", true, false},
		{"error fetching schema", args{key: "testKey", value: `{"name": "test"}`}, &mockPubSub{}, &mockSchemaClient{}, "error", true, true},
		{"success", args{key: "testKey", value: `{"name": "Rohan"}`}, &mockPubSub{}, &mockSchemaClient{}, "test_topic", false, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a, err := New(tt.mockPubSub, tt.mockSchemaClient, "latest", tt.subject)

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if err := a.PublishEvent(tt.args.key, tt.args.value, nil); (err != nil) != tt.pubErr {
					t.Errorf("PublishEvent() error = %v, wantErr %v", err, tt.pubErr)
				}
			}
		})
	}
}

func TestAvro_Subscribe(t *testing.T) {
	tests := []struct {
		name             string
		mockPubSub       pubsub.PublisherSubscriber
		mockSchemaClient *mockSchemaClient
		wantErr          bool
	}{
		{"error from subscribe", &mockPubSub{"error"}, &mockSchemaClient{}, true},
		{"unable to fetch schema", &mockPubSub{}, &mockSchemaClient{}, true},
		{"success", &mockPubSub{"id"}, &mockSchemaClient{}, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a, _ := New(tt.mockPubSub, tt.mockSchemaClient, "latest", "test_topic")

			msg, err := a.Subscribe()
			if (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}

			if msg != nil && len(msg.Headers) == 0 {
				t.Error("Subscribe() headers expected, got empty headers")
			}
		})
	}
}

func TestAvro_SubscribeWithCommit(t *testing.T) {
	commitFunc := func(msg *pubsub.Message) (bool, bool) {
		return true, false
	}

	tests := []struct {
		name             string
		mockPubSub       pubsub.PublisherSubscriber
		mockSchemaClient *mockSchemaClient
		wantErr          bool
	}{
		{"error from subscribe", &mockPubSub{"error"}, &mockSchemaClient{}, true},
		{"unable to fetch schema", &mockPubSub{}, &mockSchemaClient{}, true},
		{"success commit and stop consuming", &mockPubSub{"id"}, &mockSchemaClient{}, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a, _ := New(tt.mockPubSub, tt.mockSchemaClient, "latest", "test_topic")

			if _, err := a.SubscribeWithCommit(commitFunc); (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAvro_Ping(t *testing.T) {
	m := mockPubSub{}
	sc := mockSchemaClient{}
	a, _ := New(&m, &sc, "latest", "test_topic")

	if err := a.Ping(); err != nil {
		t.Errorf("FAILED, expected successful ping")
	}
}

func TestAvro_IsSet(t *testing.T) {
	var a *Avro
	testcases := []struct {
		a    *Avro
		resp bool
	}{
		{a, false},
		{&Avro{}, false},
	}

	for i, v := range testcases {
		resp := v.a.IsSet()
		if resp != v.resp {
			t.Errorf("[TESTCASE%d]Failed.Expected %v\tGot %v\n", i+1, v.resp, resp)
		}
	}
}

func Test_NewAvro(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respMap := map[string]interface{}{"subject": "gofr-value", "version": 2, "id": 293,
			"schema": `{"type":"record","name":"test","fields":[{"name":"ID","type":"string"}]}`}
		_ = json.NewEncoder(w).Encode(respMap)
	}))

	forbiddenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))

	var mockPubSub pubsub.PublisherSubscriber

	{ // Successful avro initialization
		tests := []struct {
			cfg *Config
		}{
			{&Config{URL: server.URL, Subject: "gofr-value"}}, // with avro subject
			{&Config{URL: server.URL, Subject: ""}},           // without avro subject
		}

		for i, tc := range tests {
			// create a avro client
			_, err := NewWithConfig(tc.cfg, mockPubSub)
			if err != nil {
				t.Errorf("Failed[%v]: got error: %v \n", i, err)
			}
		}
	}
	{ // Failure in avro initialization
		avroCfg := &Config{URL: "dummy-url", Subject: "gofr-value", Version: "latest"}
		// create a avro client
		_, err := NewWithConfig(avroCfg, mockPubSub)
		if err == nil {
			t.Errorf("Failed: want error got nil")
		}
	}
	{
		cfg := &Config{URL: forbiddenServer.URL, Subject: "gofr-value"}
		_, err := NewWithConfig(cfg, mockPubSub)
		if err == nil {
			t.Errorf("Failed: want error got nil")
		}
	}
}
