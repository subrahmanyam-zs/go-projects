# AWS SNS
Amazon Simple Notification Service (Amazon SNS) is a fully managed messaging service for both application-to-application (A2A) and application-to-person (A2P) communication.

The A2A pub/sub functionality provides topics for high-throughput, push-based, many-to-many messaging between distributed systems, microservices, and event-driven serverless applications. 

Using Amazon SNS topics, your publisher systems can fanout messages to a large number of subscriber systems

## Prerequisities
* For using AWS-SNS, we need to set the following values in the configs:
```
SNS_ACCESS_KEY = <value>              # Security credentials from AWS
SNS_SECRET_ACCESS_KEY = <value>       # Security credentials from AWS
SNS_REGION = <ap-south-1>         # The region selected for SNS service
SNS_PROTOCOL = <EMAIL/HTTP/HTTPS> # Protocol of the endpoint to be registered on pub-sub
SNS_ENDPOINT = <xyz@gmail.com>    # Used while subscribing to a topic
SNS_TOPIC_ARN = <value>           # Topic to be used for pub-sub
NOTIFIER_BACKEND = SNS               # SNS is the pub-sub backend for AWS SNS
```

* Using ACCESS_KEY and SECRET_ACCESS_KEY , gofr makes a connection to aws sns.

* ### SUBSCRIBE
     SNS uses SNS_TOPIC_ARN, SNS_PROTOCOL and SNS_ENDPOINT in subcribing.
     SNS_PROTOCOL varies according to endpoint being subscribed like EMAIL for subscribing an email endpoint and HTTPS / HTTP for subscribing to a server endpoint.
    
- ### PUBLISH 
     SNS uses SNS_TOPIC_ARN to publish the data provided.
  
* To run the example simply provide these configurations either in environment variables or in the configs folder.
  And run <code>go run main.go</code>.
  