FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-eventhub/build/main /main
COPY examples/using-eventhub/configs /configs
EXPOSE 9113
CMD ["/main"]