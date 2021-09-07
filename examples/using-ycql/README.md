## Yugabyte Example

To run the example follow the steps below :

1) Run the docker image of cassandra

   > `docker run --name gofr-yugabyte -p 2011:9042 -d yugabytedb/yugabyte:latest bin/yugabyted start --daemon=false`

2) Now on the project folder path `gofr/` run the following command to load the schema

   > `docker exec -i gofr-yugabyte ycqlsh < .github/setups/keyspace.ycql`

3) Now you can run the example on path `gofr/examples/using-ycql` by

   > `go run main.go`