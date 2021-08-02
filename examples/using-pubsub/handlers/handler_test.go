package handlers

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/metrics"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
)

type mockPubSub struct {
	id string
}

func (m *mockPubSub) SubscribeWithCommit(commitFunc pubsub.CommitFunc) (*pubsub.Message, error) {
	pubsubMessage := &pubsub.Message{
		Offset: 1,
		Topic:  "test-topic",
	}

	count := 0

	for {
		pubsubMessage.Offset = int64(count + 1)

		_, isContinue := commitFunc(pubsubMessage)
		if !isContinue {
			break
		}
	}

	if m.id == "error" {
		return nil, &errors.Response{Reason: "test-error"}
	}

	return &pubsub.Message{}, nil
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

func (m *mockPubSub) Bind(v []byte, target interface{}) error {
	return nil
}

func (m *mockPubSub) Ping() error {
	return nil
}

func (m *mockPubSub) IsSet() bool {
	return true
}

func (m *mockPubSub) HealthCheck() types.Health {
	return types.Health{}
}

func (m *mockPubSub) CommitOffset(offsets pubsub.TopicPartition) {
}

func TestProducerHandler(t *testing.T) {
	k := gofr.New()
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	k.Metric = mockMetric

	mockMetric.EXPECT().ObserveHistogram(gomock.Any(), gomock.Any()).Return(nil)

	m := mockPubSub{}
	k.PubSub = &m

	tests := []struct {
		name    string
		id      string
		want    interface{}
		wantErr bool
	}{
		{"error from publisher", "1", nil, true},
		{"success", "123", nil, false},
	}

	req := httptest.NewRequest("GET", "http://dummy", nil)
	context := gofr.NewContext(nil, request.NewHTTPRequest(req), k)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			context.SetPathParams(map[string]string{
				"id": tt.id,
			})

			m.id = tt.id

			got, err := Producer(context)
			if (err != nil) != tt.wantErr {
				t.Errorf("Producer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Producer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsumerHandler(t *testing.T) {
	k := gofr.New()
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	k.Metric = mockMetric

	mockMetric.EXPECT().ObserveSummary(gomock.Any(), gomock.Any()).Return(nil)

	k.PubSub = &mockPubSub{}

	ctx := gofr.NewContext(nil, nil, k)

	_, err := Consumer(ctx)
	if err != nil {
		t.Errorf("Consumer() error = %v, wantErr %v", err, nil)
		return
	}
}

func TestConsumerWithCommitHandler(t *testing.T) {
	k := gofr.New()

	testCases := []struct {
		name        string
		pubSub      pubsub.PublisherSubscriber
		expectedErr error
	}{
		{"success consuming messages", &mockPubSub{}, nil},
		{"error from pubsub subscribe", &mockPubSub{"error"}, &errors.Response{Reason: "test-error"}},
	}

	for _, tc := range testCases {
		k.PubSub = tc.pubSub
		ctx := gofr.NewContext(nil, nil, k)

		_, err := ConsumerWithCommit(ctx)

		assert.Equal(t, tc.expectedErr, err, tc.name)
	}
}
