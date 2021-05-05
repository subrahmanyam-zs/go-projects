## WebSocket Message Example

The request must contain the following header value pairs for the web-socket connection to be intialized:

| Header                | Value                         |
|-----------------------|-------------------------------|
| Connection            | upgrade                       |
| Upgrade               | websocket                     |
| Sec-Websocket-Version | 13                            |
| Sec-Websocket-Key     | 'base64 encoded random bytes' |

`Sec-Websocket-Key` is meant to prevent proxies from caching the request, by sending a random key. It does not provide any authentication.

####Follow the steps below to run the example:<br/>
```
1. Run the example using the command: go run main.go.
2. In browser, go to http://localhost:9101.
3. Send a message using the textbox and the send button which appears in the browser.
4. Messages sent from the browser appear on the terminal console.
```
