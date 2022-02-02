FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-dynamodb/build/main /main
COPY examples/using-dynamodb/configs /configs
EXPOSE 9091
CMD ["/main"]