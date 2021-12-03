FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/data-layer-with-cassandra/build/main /main
COPY examples/data-layer-with-cassandra/configs /configs
EXPOSE 9094
CMD ["/main"]