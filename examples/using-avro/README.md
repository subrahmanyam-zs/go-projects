#AVRO PubSub

##Instructions

Avro pub sub uses Kafka to write or read messages.

In order for this example to run:

- If PUBSUB_BACKEND is AVRO

    1. Provide the mandatory configs required for avro: AVRO_SCHEMA_URL
    2. Avro uses KAFKA streams, so mandatory configs for Kafka needs to be provided: KAFKA_HOSTS and KAFKA_TOPIC
    
- If PUBSUB_BACKEND is KAFKA

    1. KAFKA_HOSTS and KAFKA_TOPIC are the mandatory configs