FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-avro/build/main /main
COPY examples/using-avro/configs /configs
EXPOSE 9111
CMD ["/main"]