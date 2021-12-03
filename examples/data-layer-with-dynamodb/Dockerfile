FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/data-layer-with-dynamodb/build/main /main
COPY examples/data-layer-with-dynamodb/configs /configs
EXPOSE 9091
CMD ["/main"]