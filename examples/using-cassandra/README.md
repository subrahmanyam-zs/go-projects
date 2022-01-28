## Cassandra Example

To run the example follow the steps below :

1) Run the docker image of cassandra

   > `docker run --name gofr-cassandra -d -p 2003:9042 cassandra:latest`

2) Now on the project folder path `/zopsmart/gofr/` run the following command to load the schema

   > `docker exec -i gofr-cassandra cqlsh < .github/setups/keyspace.cql`

3) Now you can run the example on path `/zopsmart/gofr/examples/using-cassandra` by

   > `go run main.go`