## Redis Example
#### To run the example follow the steps below:
1) Run the docker image of redis

   > `docker run --name gofr-redis -p 2002:6379 -d redis:latest`
   
2) Now run the example on path `/zopsmart/gofr/examples/using-redis` by

   > `go run main.go`
   