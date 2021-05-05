package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zopsmart/gofr/pkg/datastore/pubsub"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
	"github.com/zopsmart/gofr/pkg/gofr/types"
)

type mockPubSub struct {
	id string
}

func TestAVROProducerHandler(t *testing.T) {
	k := gofr.New()
	m := mockPubSub{}
	k.PubSub = &m

	tests := []struct {
		name         string
		id           string
		expectedResp interface{}
		expectedErr  error
	}{
		{"error from publisher", "1", nil, errors.EntityNotFound{Entity: "", ID: "1"}},
		{"success", "123", nil, nil},
	}

	req := httptest.NewRequest("GET", "http://dummy", nil)
	context := gofr.NewContext(nil, request.NewHTTPRequest(req), k)

	for _, tt := range tests {
		context.SetPathParams(map[string]string{
			"id": tt.id,
		})

		m.id = tt.id

		gotResp, gotErr := Producer(context)
		assert.Equal(t, gotErr, tt.expectedErr)
		assert.Equal(t, gotResp, tt.expectedResp)
	}
}

func TestAVROConsumerHandler(t *testing.T) {
	k := gofr.New()

	k.PubSub = &mockPubSub{}

	ctx := gofr.NewContext(nil, nil, k)

	_, err := Consumer(ctx)
	assert.Equal(t, nil, err)
}

func (m *mockPubSub) HealthCheck() types.Health {
	return types.Health{}
}

func (m *mockPubSub) IsSet() bool {
	return false
}

func (m *mockPubSub) PublishEventWithOptions(key string, val interface{}, headers map[string]string, options *pubsub.PublishOptions) error {
	return nil
}

func (m *mockPubSub) PublishEvent(key string, val interface{}, headers map[string]string) error {
	if m.id == "1" {
		return errors.EntityNotFound{ID: "1"}
	}

	return nil
}

func (m *mockPubSub) Subscribe() (*pubsub.Message, error) {
	return &pubsub.Message{}, nil
}

func (m *mockPubSub) SubscribeWithCommit(commitFunc pubsub.CommitFunc) (*pubsub.Message, error) {
	return nil, nil
}

func (m *mockPubSub) Bind(v []byte, target interface{}) error {
	return nil
}

func (m *mockPubSub) Ping() error {
	return nil
}

func (m *mockPubSub) CommitOffset(offsets pubsub.TopicPartition) {
	return
}
