package awssns

import (
	"encoding/json"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/notifier"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/prometheus/client_golang/prometheus"
)

type SNS struct {
	sns AWS
	cfg *Config
}

type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	// ConnRetryDuration for specifying connection retry duration
	ConnRetryDuration int
	TopicArn          string
	Endpoint          string
	Protocol          string
}

//nolint // The declared global variable can be accessed across multiple functions
var (
	notifierReceiveCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_notifier_receive_count",
		Help: "Total number of subscribe operation",
	}, []string{"topic"})

	notifierSuccessCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_notifier_success_count",
		Help: "Total number of successful subscribe operation",
	}, []string{"topic"})

	notifierFailureCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_notifier_failure_count",
		Help: "Total number of failed subscribe operation",
	}, []string{"topic"})
)

func New(c *Config) (notifier.Notifier, error) {
	_ = prometheus.Register(notifierReceiveCount)
	_ = prometheus.Register(notifierSuccessCount)
	_ = prometheus.Register(notifierFailureCount)

	sessionConfig := &aws.Config{
		Region:      aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, ""),
	}

	sess, _ := session.NewSession(sessionConfig)

	svc := sns.New(sess, nil)

	return &SNS{sns: svc, cfg: c}, nil
}

func (s *SNS) Publish(value interface{}) (err error) {
	data, ok := value.([]byte)
	if !ok {
		data, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(data)),
		TopicArn: aws.String(s.cfg.TopicArn),
	}

	_, err = s.sns.Publish(input)
	if err != nil {
		return err
	}

	return nil
}

func (s *SNS) Subscribe() (*notifier.Message, error) {
	// for every subscribe
	notifierReceiveCount.WithLabelValues(s.cfg.Endpoint).Inc()

	out, err := s.sns.Subscribe(&sns.SubscribeInput{
		Endpoint:              &s.cfg.Endpoint,
		Protocol:              &s.cfg.Protocol,
		ReturnSubscriptionArn: aws.Bool(true), // Return the ARN, even if user has yet to confirm
		TopicArn:              &s.cfg.TopicArn,
	})

	if err != nil {
		// for failed subscribe
		notifierFailureCount.WithLabelValues(s.cfg.Endpoint).Inc()
		return nil, err
	}

	msg, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}

	// for successful subscribe
	notifierSuccessCount.WithLabelValues(s.cfg.Endpoint).Inc()

	return &notifier.Message{Value: string(msg)}, nil
}

func (s *SNS) SubscribeWithResponse(target interface{}) (*notifier.Message, error) {
	message, err := s.Subscribe()
	if err != nil {
		return message, err
	}

	return message, s.Bind([]byte(message.Value), &target)
}

func (s *SNS) Bind(message []byte, target interface{}) error {
	return json.Unmarshal(message, target)
}

func (s *SNS) ping() error {
	sessionConfig := &aws.Config{
		Region:      aws.String(s.cfg.Region),
		Credentials: credentials.NewStaticCredentials(s.cfg.AccessKeyID, s.cfg.SecretAccessKey, ""),
	}

	sess, _ := session.NewSession(sessionConfig)

	svc := sns.New(sess, nil)
	if svc == nil {
		return errors.AWSSessionNotCreated
	}

	return nil
}

func (s *SNS) HealthCheck() types.Health {
	if s == nil {
		return types.Health{
			Name:   pkg.AWSSNS,
			Status: pkg.StatusDown,
		}
	}

	resp := types.Health{
		Name:   pkg.AWSSNS,
		Status: pkg.StatusDown,
		Host:   s.cfg.TopicArn,
	}

	if err := s.ping(); err != nil {
		return resp
	}

	resp.Status = pkg.StatusUp

	return resp
}

func (s *SNS) IsSet() bool {
	if s == nil {
		return false
	}

	return s.sns != nil
}
