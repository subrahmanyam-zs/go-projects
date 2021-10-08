package handlers

import (
	"net/http"
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

func (m *mockPubSub) CommitOffset(offsets pubsub.TopicPartition) {
}

func (m *mockPubSub) PublishEventWithOptions(key string, val interface{}, headers map[string]string, options *pubsub.PublishOptions) error {
	if m.id == "1" {
		return errors.EntityNotFound{ID: "1"}
	}

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
	return &pubsub.Message{}, nil
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

func TestProducerHandler(t *testing.T) {
	app := gofr.New()
	m := mockPubSub{}
	app.PubSub = &m

	tests := []struct {
		desc string
		id   string
		resp interface{}
		err  error
	}{
		{"error from publisher", "1", nil, errors.EntityNotFound{ID: "1"}},
		{"success", "123", nil, nil},
	}

	req := httptest.NewRequest(http.MethodGet, "http://dummy", nil)
	ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), app)

	for i, tc := range tests {
		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		m.id = tc.id
		resp, err := Producer(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestConsumerHandler(t *testing.T) {
	app := gofr.New()
	app.PubSub = &mockPubSub{}

	ctx := gofr.NewContext(nil, nil, app)
	_, err := Consumer(ctx)

	assert.Equal(t, nil, err)
}
