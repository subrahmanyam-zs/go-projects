## Solr Example

To run the example follow the steps below :

1) Run the docker image of solr

   > `docker run --name gofr-solr -p 2020:8983 solr:8 -DzkRun`

2) Now on the project folder path `gofr/` run the following command to load the schema

   > `docker exec -i gofr-solr sh < .github/setups/solrSchema.sh;`

3) Now you can run the example on path `gofr/examples/using-solr` by

   > `go run main.go`
