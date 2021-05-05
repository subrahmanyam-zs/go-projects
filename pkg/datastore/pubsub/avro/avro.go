package avro

import (
	"encoding/binary"
	"strings"

	"github.com/hamba/avro"
	"github.com/zopsmart/gofr/pkg/datastore/pubsub"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr/types"
)

type Avro struct {
	schemaVersion        string
	subject              string
	schemaRegistryClient SchemaRegistryClientInterface
	pubSub               pubsub.PublisherSubscriber
	schema               avro.Schema
	schemaID             int
	subSchema            avro.Schema
}

type Config struct {
	URL            string
	Version        string
	Subject        string
	SchemaUser     string
	SchemaPassword string
}

func NewWithConfig(c *Config, ps pubsub.PublisherSubscriber) (pubsub.PublisherSubscriber, error) {
	if c == nil || c.URL == "" {
		return nil, nil
	}

	if c.Version == "" {
		c.Version = "latest"
	}

	registryURLSlc := strings.Split(c.URL, ",")
	schemaRegistryClient := NewSchemaRegistryClient(registryURLSlc, c.SchemaUser, c.SchemaPassword)

	return New(ps, schemaRegistryClient, c.Version, c.Subject)
}

func New(ps pubsub.PublisherSubscriber, src SchemaRegistryClientInterface, version, sub string) (pubsub.PublisherSubscriber, error) {
	avroPubSub := &Avro{
		schemaVersion:        version,
		schemaRegistryClient: src,
		pubSub:               ps,
	}
	// Avro should be initialized even if subject is not provided
	if sub == "" {
		return avroPubSub, nil
	}

	schemaID, schemaStr, err := src.GetSchemaByVersion(sub, version)
	if err != nil {
		return nil, err
	}

	schema, err := avro.Parse(schemaStr)
	if err != nil {
		return nil, err
	}

	avroPubSub.schemaID = schemaID
	avroPubSub.schema = schema
	avroPubSub.subject = sub

	return avroPubSub, nil
}

func (a *Avro) PublishEventWithOptions(key string, value interface{}, headers map[string]string, options *pubsub.PublishOptions) error {
	// Missing schema will generate panic
	if a.schema == nil {
		return &errors.Response{Code: "Missing schema", Reason: "Avro is initialized without schema"}
	}

	binaryValue, err := avro.Marshal(a.schema, value)
	if err != nil {
		return err
	}

	encodedMsg := Encoder{
		SchemaID: a.schemaID,
		Content:  binaryValue,
	}

	binaryEncoded := encodedMsg.Encode()

	return a.pubSub.PublishEventWithOptions(key, binaryEncoded, headers, options)
}

func (a *Avro) PublishEvent(key string, value interface{}, headers map[string]string) error {
	return a.PublishEventWithOptions(key, value, headers, nil)
}

func (a *Avro) Subscribe() (*pubsub.Message, error) {
	msg, err := a.pubSub.Subscribe()
	if err != nil {
		return nil, err
	}

	return a.processMessage(msg)
}

func (a *Avro) SubscribeWithCommit(f pubsub.CommitFunc) (*pubsub.Message, error) {
	for {
		msg, err := a.pubSub.SubscribeWithCommit(nil)
		if err != nil {
			return nil, err
		}

		msg, err = a.processMessage(msg)
		if err != nil {
			return nil, err
		}

		isCommit, isContinue := f(msg)
		if isCommit {
			a.CommitOffset(pubsub.TopicPartition{
				Topic:     msg.Topic,
				Partition: msg.Partition,
				Offset:    msg.Offset,
			})
		}

		if !isContinue {
			return msg, nil
		}
	}
}

func (a *Avro) processMessage(msg *pubsub.Message) (*pubsub.Message, error) {
	value := []byte(msg.Value)
	schemaID := binary.BigEndian.Uint32(value[1:5])

	finalMsg := &pubsub.Message{
		SchemaID:  int(schemaID),
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Key:       msg.Key,
		Value:     msg.Value[5:],
		Headers:   msg.Headers,
	}

	schema, err := a.schemaRegistryClient.GetSchema(int(schemaID))
	if err != nil {
		return finalMsg, err
	}

	a.subSchema, _ = avro.Parse(schema)

	return finalMsg, err
}

// Encoder encodes schemaId and Avro message.
type Encoder struct {
	SchemaID int
	Content  []byte
}

// Note: the Confluent schema registry has special requirements for the Avro serialization rules,
// not only need to serialize the specific content, but also attach the Schema ID and Magic Byte.
// Ref: https://docs.confluent.io/current/schema-registry/serializer-formatter.html#wire-format
func (a *Encoder) Encode() []byte {
	var binaryMsg []byte

	// Confluent serialization format version number; currently always 0.
	binaryMsg = append(binaryMsg, byte(0))

	// 4-byte schema ID as returned by Schema Registry
	binarySchemaID := make([]byte, 4)

	binary.BigEndian.PutUint32(binarySchemaID, uint32(a.SchemaID))

	binaryMsg = append(binaryMsg, binarySchemaID...)

	// Avro serialized data in Avro's binary encoding
	binaryMsg = append(binaryMsg, a.Content...)

	return binaryMsg
}

func (a *Avro) Bind(message []byte, target interface{}) error {
	return avro.Unmarshal(a.subSchema, message, target)
}

func (a *Avro) Ping() error {
	return a.pubSub.Ping()
}

func (a *Avro) HealthCheck() types.Health {
	return a.pubSub.HealthCheck()
}

func (a *Avro) IsSet() bool {
	if a == nil {
		return false
	}

	if a.pubSub == nil || a.schema == nil || a.schemaRegistryClient == nil {
		return false
	}

	return true
}

func (a *Avro) CommitOffset(offsets pubsub.TopicPartition) {
	a.pubSub.CommitOffset(offsets)
}
