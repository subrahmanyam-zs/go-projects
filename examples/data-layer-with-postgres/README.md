## Run Postgres Example

Steps to run the example are :

1) Run the docker image

   > `docker run --name gofr-pgsql -e POSTGRES_DB=customers -e POSTGRES_PASSWORD=root123 -p 2006:5432 -d postgres:latest
   `
2) Now on the project path `zopsmart/gofr` run the following command to load the schema

   > `docker exec -i gofr-pgsql psql -U postgres customers < .github/setups/setup.sql`

3) Now run the following command in example folder path `zopsmart/gofr/examples/using-potgres` to run the example
   endpoints

   > `go run main.go`
