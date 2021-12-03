FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/data-layer-with-elasticsearch/build/main /main
COPY examples/data-layer-with-elasticsearch/configs /configs
EXPOSE 8001
CMD ["/main"]