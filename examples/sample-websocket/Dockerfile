FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/sample-websocket/build/main /main
COPY examples/sample-websocket/configs /configs
EXPOSE 9101
CMD ["/main"]