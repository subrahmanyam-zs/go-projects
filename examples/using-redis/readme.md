#### Headers
The request must contain the following header value pairs:

| Header                | Value                         |
|-----------------------|-------------------------------|
| X-Zopsmart-Tenant     | good4more                     |
| X-Correlation-ID      | 1s3d323adsd                   |
| True-Client-Ip        | 127.0.0.1                     |


#### To Create a Redis container
```
docker run --name gofr-redis -p 2002:6379 -d redis:latest
docker run --name container-redis -d redis
```

#### Sample request and response
_POST_

url : http://localhost:9091/config

json body:
```
{
    "id":"1",
    "name":"xyz"
}
```
response:
```
{
    "data": "Successful"
}
```

_GET_

url : http://localhost:9091/config/id

response:
```
{
   "data": {
        "id": "1"
    }
}
```
_DELETE_

url : http://localhost:9091/config/id

response:

```
{
    "data": "Deleted successfully"
}
```
