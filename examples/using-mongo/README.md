# Example: using MongoDB

### Prerequisites
To run the the example you need to have MongoDB installed on your system or docker container.

1. To run the docker container run the following command:
- `docker run --name gofr-mongo -d -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin123 -p 2004:27017 mongo:latest`


2. Set the following environment variables:
    ```
    APP_NAME=gofr-mongo-example    
    MONGO_DB_HOST=localhost
    MONGO_DB_PORT=2004  // exposed port of container
    MONGO_DB_USER=admin
    MONGO_DB_PASS=admin123
    MONGO_DB_NAME=test
    ```
### Example
Run the application
```
    go run main.go
```

Endpoints:
```
    /customer  GET, POST, DELETE
```

Sample Request
```
    POST: http://localhost:9097/customer
    Headers:
       X-Correlation-ID=1s3d323adsd
       X-Zopsmart-Tenant=good4more
       True-Client-Ip=127.0.0.1
    Body:
        {
            "id":1,
            "name":"person",
            "age":12,
            "city":"Gujarat"
        }
```
Sample Response:
```
    &{New Customer Added! <nil>}
```