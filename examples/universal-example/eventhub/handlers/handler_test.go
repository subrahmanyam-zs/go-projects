package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
)

type mockPubSub struct {
	id string
}

func TestProducerHandler(t *testing.T) {
	k := gofr.New()

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

	for _, tc := range tests {
		context.SetPathParams(map[string]string{
			"id": tc.id,
		})

		gotResp, gotErr := New(&mockPubSub{id: tc.id}).Producer(context)
		assert.Equal(t, tc.expectedErr, gotErr)
		assert.Equal(t, tc.expectedResp, gotResp)
	}
}

func TestConsumerHandler(t *testing.T) {
	k := gofr.New()

	ctx := gofr.NewContext(nil, nil, k)
	tests := []struct {
		id          string
		expectedErr error
	}{
		// Success Case
		{"", nil},
		// Failure Case
		{"1", errors.EntityNotFound{Entity: "", ID: "1"}},
	}

	for _, tc := range tests {
		_, gotErr := New(&mockPubSub{id: tc.id}).Consumer(ctx)
		assert.Equal(t, tc.expectedErr, gotErr)
	}
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
	if m.id == "1" {
		return nil, errors.EntityNotFound{ID: "1"}
	}

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

func (m *mockPubSub) HealthCheck() types.Health {
	return types.Health{}
}

func (m *mockPubSub) IsSet() bool {
	return false
}
