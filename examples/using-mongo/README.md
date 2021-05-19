# Example: using MongoDB

### Prerequisites
To run the the example you need to have MongoDB installed on your system or docker container.

To run the docker container run the following command:
- `docker run --name gofr-mongo -d -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin123 -p 2004:27017 mongo:latest`

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
Body:
    {
       "id":1,
       "name":"person",
       "age":12,
       "city":"Banglore"
    }
```
Sample Response:
```
&{New Customer Added! <nil>}
```
