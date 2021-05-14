## Server Setup
- Run the docker command:
 ```shell
 docker run --name gofr-yugabyte -d -p2021:7000 -p2010:9000 -p2023:5433 -p2011:9042 -v ~/yb_data:/home/yugabyte/var yugabytedb/yugabyte:latest bin/yugabyted start --daemon=false
```
_This will get the YugaByte service up and running for us._
> If you do not wish to use the predefined configs, you can adjust them in the `configs/.env` file

- Now, run `main_test.go` to load the **_test_** keyspace and **_shop_** table schema in the container for us.
```shell
go test
```

You'll see something like this in the terminal:
```shell
PASS
ok      github.com/zopsmart/gofr/examples/using-ycql    5.363s
```

_This indicates that the integration test for the `sample-ycql` have passed._

- Now run the sample server by typing: 
```shell
go run main.go
```
_This starts the server at **PORT 9005**, as configured in `main.go`_
<hr>

##  Sample Requests
You can use the following **cURL** commands to send request to the available endpoints:

  ### `GET /shop`  Sample Requests 

  - Example 1: Only name as a query parameter  
  
_Request:_
  ```shell
  curl --location --request GET 'localhost:9005/shop?name=Kalash'
  ```
  _Response:_
  ```shell
  {"data":[{"id":7,"name":"Kalash","location":"Jehanabad","state":"Bihar"}]}
  ```
  
  - Example 2: Name and location as query parameter  
  
_Request:_
  ```shell
  curl --location --request GET 'localhost:9005/shop?name=Kalash&location=Jehanabad'
  ```
  _Response:_
  ```shell
  {"data":[{"id":7,"name":"Kalash","location":"Jehanabad","state":"Bihar"}]}
  ```
  
  - Example 3: All Parameters in Query  
  
_Request:_  
  ```shell
  curl --location --request GET 'localhost:9005/shop?name=Kalash&location=Jehanabad&id=7&state=Bihar'
  ```
  _Response:_
  ```shell
  {"data":[{"id":7,"name":"Kalash","location":"Jehanabad","state":"Bihar"}]}
  ```
  
### `POST /shop`  Sample Requests   
  
  - Example : Valid Request  

  _Request:_  
  ```shell
  curl --location --request POST 'localhost:9005/shop' \
  --header 'Content-Type: application/json' \
  --data-raw '{"id":4, "name": "UBCity", "location":"HSR", "State":"karnataka"}'
  ```
  _Response:_  
   ```shell
  {"data":[{"id":4,"name":"UBCity","location":"HSR","state":"karnataka"}]}
  ```
  
### `PUT /shop/{id}`  Sample Requests   
  
  - Example 1: Valid Update

  _Request:_  
  ```shell
  curl --location --request PUT 'localhost:9005/shop/4' \
  --header 'Content-Type: application/json' \
  --data-raw '{"id":4, "name": "UBCity", "location":"HSR, Sector-5", "State":"karnataka"}'
  ```
  _Response:_  
  ```shell
  {"data":[{"id":4,"name":"UBCity","location":"HSR, Sector-5","state":"karnataka"}]}
  ```
  
  - Example 2: Invalid Update on non-existing ID  

  _Request:_  
  ```shell
  curl --location --request PUT 'localhost:9005/shop/1' \
  --header 'Content-Type: application/json' \
  --data-raw '{"id":4, "name": "UBCity", "location":"HSR, Sector-5", "State":"karnataka"}'
  ```
  _Response:_
  ```shell
  {"errors":[{"code":"Entity Not Found","reason":"No 'person' found for Id: '1'","datetime":{"value":"2021-05-13T23:52:51Z","timezone":"IST"}}]}
  ```

### `DELETE /shop/{id}`  Sample Requests   
  
  - Example 1: Valid Delete  

  _Request:_  
   ```shell
  curl --location --request DELETE 'localhost:9005/shop/4'
  ```
  _Response:_
  ```shell
  #no output
  ```
  
  - Example 2: Invalid Delete on non-existing ID  

   _Request:_  
  ```shell
  curl --location --request DELETE 'localhost:9005/shop/1'
  ```
  _Response:_
  ```shell
  {"errors":[{"code":"Entity Not Found","reason":"No 'person' found for Id: '1'","datetime":{"value":"2021-05-13T23:55:10Z","timezone":"IST"}}]}
  ```

> Validate Headers has been set to false in `main.go`, via `k.Server.ValidateHeaders = false`. If enabled, sample values for header fields could be used:
> `X-Correlation-ID = 1s3d323adsd`, `X-Zopsmart-Tenant = good4more`, `True-Client-Ip = 127.0.0.1`