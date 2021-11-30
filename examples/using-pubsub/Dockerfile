FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-pubsub/build/main /main
COPY examples/using-pubsub/configs /configs
EXPOSE 9112
CMD ["/main"]