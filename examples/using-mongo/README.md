## MongoDB Example

Steps to run the example using docker image of MongoDB: 

1) Run the docker image of MongoDB 
    
    > `docker run --name gofr-mongo -d -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin123 -p 2004:27017 mongo:latest`

2) Run the example on path `/zopsmart/gofr/examples/using-mongo` by 
        
    > `go run main.go
