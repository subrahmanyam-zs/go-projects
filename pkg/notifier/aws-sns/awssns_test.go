package awssns

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/notifier"
)

func TestNew(t *testing.T) {
	cfg := config.NewGoDotEnvProvider(log.NewLogger(), "../../../configs")

	tests := []struct {
		desc    string
		c       Config
		wantErr bool
	}{
		{desc: "Success connection", c: Config{
			AccessKeyID:     cfg.Get("SNS_ACCESS_KEY"),
			SecretAccessKey: cfg.Get("SNS_SECRET_ACCESS_KEY"),
			Region:          "dummy",
		}},
		{desc: "Dummy connection", c: Config{
			AccessKeyID:     "dummy",
			SecretAccessKey: "dummy",
			Region:          "SNS_Region",
		}},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			_, err := New(&tc.c)

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestSNS_Subscribe(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := NewMockAWS(mockCtrl)

	cfg := &Config{AccessKeyID: "1", SecretAccessKey: "1", TopicArn: "arn", Endpoint: "api", Protocol: "proto"}

	svc := SNS{cfg: cfg, sns: mockService}

	tests := []struct {
		desc    string
		expOut  *notifier.Message
		wantErr error
	}{
		{
			desc:   "Success Case",
			expOut: &notifier.Message{Value: fmt.Sprintf(`{"SubscriptionArn":"%s"}`, svc.cfg.TopicArn)},
		},
		{
			desc:    "Failure Case",
			wantErr: errors.EntityNotFound{},
		},
	}

	for _, tc := range tests {
		mockService.EXPECT().Subscribe(&sns.SubscribeInput{Endpoint: &svc.cfg.Endpoint, Protocol: &svc.cfg.Protocol,
			ReturnSubscriptionArn: aws.Bool(true), TopicArn: &svc.cfg.TopicArn}).
			Return(&sns.SubscribeOutput{SubscriptionArn: &svc.cfg.TopicArn}, tc.wantErr)

		result, err := svc.Subscribe()

		assert.Equalf(t, tc.expOut, result, "%v Expected Output : %v , got %v", tc.desc, tc.expOut, result)

		assert.ErrorIsf(t, err, tc.wantErr, " %s Expected Error : %v , got %v", tc.desc, tc.wantErr, err)
	}
}

func TestSNS_SubscribeWithResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := NewMockAWS(mockCtrl)

	cfg := &Config{AccessKeyID: "1", SecretAccessKey: "1", TopicArn: "arn", Endpoint: "api", Protocol: "proto"}

	svc := SNS{cfg: cfg, sns: mockService}

	tests := []struct {
		desc    string
		expOut  *notifier.Message
		wantErr error
	}{
		{
			desc:   "Success Case",
			expOut: &notifier.Message{Value: fmt.Sprintf(`{"SubscriptionArn":"%s"}`, svc.cfg.TopicArn)},
		},
		{
			desc:    "Failure Case",
			wantErr: errors.EntityNotFound{},
		},
	}

	for _, tc := range tests {
		mockService.EXPECT().Subscribe(&sns.SubscribeInput{Endpoint: &svc.cfg.Endpoint, Protocol: &svc.cfg.Protocol,
			ReturnSubscriptionArn: aws.Bool(true), TopicArn: &svc.cfg.TopicArn}).
			Return(&sns.SubscribeOutput{SubscriptionArn: &svc.cfg.TopicArn}, tc.wantErr)
		var tar interface{}
		result, err := svc.SubscribeWithResponse(tar)

		assert.Equalf(t, tc.expOut, result, "%v Expected Output : %v , got %v", tc.desc, tc.expOut, result)

		assert.ErrorIsf(t, err, tc.wantErr, " %s Expected Error : %v , got %v", tc.desc, tc.wantErr, err)

	}
}

func TestSNS_PublishEvent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := NewMockAWS(mockCtrl)

	cfg := &Config{TopicArn: "arn"}

	svc := SNS{cfg: cfg, sns: mockService}

	tests := []struct {
		desc         string
		inputPublish *sns.PublishInput
		inputValue   interface{}
		wantErr      error
	}{
		{
			desc: "Success Case",
			inputPublish: &sns.PublishInput{Message: aws.String(`{"framework":"GOFR"}`),
				TopicArn: aws.String(svc.cfg.TopicArn)},
			inputValue: map[string]interface{}{"framework": "GOFR"},
		},
		{
			desc: "Failure Case",
			inputPublish: &sns.PublishInput{Message: aws.String(`{"framework":"GOFR"}`),
				TopicArn: aws.String(svc.cfg.TopicArn)},
			inputValue: map[string]interface{}{"framework": "GOFR"},
			wantErr:    errors.EntityNotFound{},
		},
	}

	for _, tc := range tests {
		mockService.EXPECT().Publish(tc.inputPublish).Return(&sns.PublishOutput{}, tc.wantErr)

		err := svc.Publish(tc.inputValue)

		assert.ErrorIsf(t, err, tc.wantErr, " %s Expected Error : %v , got %v", tc.desc, tc.wantErr, err)
	}
}

func TestSNS_PublishEventMarshalError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := NewMockAWS(mockCtrl)

	cfg := &Config{TopicArn: "arn"}

	svc := SNS{cfg: cfg, sns: mockService}

	tests := struct {
		desc       string
		inputValue interface{}
		wantErr    error
	}{
		desc:       "Success Case",
		inputValue: make(chan int),
		wantErr:    fmt.Errorf("json: unsupported type: chan int"),
	}

	err := svc.Publish(tests.inputValue)

	assert.Error(t, err, " %s Expected Error : %v , got %v", tests.desc, tests.wantErr, err)
}

func TestSNS_HealthCheck(t *testing.T) {
	cfg := config.NewGoDotEnvProvider(log.NewLogger(), "../../../configs")
	testcases := []struct {
		c    Config
		resp types.Health
	}{
		// Correct Credentials
		{c: Config{
			AccessKeyID:     cfg.Get("SNS_ACCESS_KEY"),
			SecretAccessKey: cfg.Get("SNS_SECRET_ACCESS_KEY"),
			Region:          cfg.Get("SNS_REGION")},
			resp: types.Health{Name: pkg.AWSSNS, Status: pkg.StatusUp},
		},
	}

	for i, v := range testcases {
		conn, _ := New(&v.c)

		resp := conn.HealthCheck()

		assert.Equalf(t, v.resp, resp, "[TESTCASE%d]Failed.Expected %v\tGot %v\n", i+1, v.resp, resp)
	}
}

func TestSNS_HealthCheckDown(t *testing.T) {
	var s *SNS
	expected := types.Health{
		Name:   pkg.AWSSNS,
		Status: pkg.StatusDown,
	}

	resp := s.HealthCheck()

	assert.Equalf(t, expected, resp, "Expected %v\tGot %v\n", expected, resp)
}

func TestSNS_IsSet(t *testing.T) {
	var s *SNS
	logger := log.NewMockLogger(ioutil.Discard)
	cfg := config.NewGoDotEnvProvider(logger, "../../../configs")
	conn, _ := New(&Config{
		AccessKeyID:     cfg.Get("SNS_ACCESS_KEY"),
		SecretAccessKey: cfg.Get("SNS_SECRET_ACCESS_KEY"),
		Region:          cfg.Get("SNS_REGION"),
	})

	testcases := []struct {
		notifier notifier.Notifier
		resp     bool
	}{
		{notifier: s},
		{notifier: &SNS{}},
		{notifier: conn, resp: true},
	}

	for i, v := range testcases {
		resp := v.notifier.IsSet()
		assert.Equalf(t, v.resp, resp, "[TESTCASE%d]Failed.Expected %v\tGot %v\n", i+1, v.resp, resp)
	}
}

func TestSNS_Bind(t *testing.T) {
	svc := &SNS{}
	message := map[string]interface{}{}

	val := []byte(`{"message":"Hi Gofr"}`)

	err := svc.Bind(val, &message)

	assert.NoError(t, err, "Error was not Expected while valid Unmarshalling")
}

func TestSNS_BindError(t *testing.T) {
	svc := &SNS{}
	message := map[string]interface{}{}

	val := []byte(`{`)

	err := svc.Bind(val, &message)

	assert.Error(t, err, "Expected Error but got nothing")
}
