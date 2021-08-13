## DynamoDB with Gofr

DynamoDB is available as web-service, DynamoDB local and docker container. For these example we will user docker container.

</br>Run docker Container using following command:
> docker run -p 8000:8000 --name gofr-dynamodb -d amazon/dynamodb-local

***Note :*** we have exposed dynamodb service at port 8000 on localhost.So we can perform operations on endpoint-url `http://localhost:8000`. 

<hr> </hr>

## Steps to run the example

1) Create table by running the following command in `/examples/using-dynamodb` path

   > go test

   ***Note :*** It should PASS all the test cases <br></br>

2) Run the application

   > go run main.go

