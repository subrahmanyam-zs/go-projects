FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-elasticsearch/build/main /main
COPY examples/using-elasticsearch/configs /configs
EXPOSE 8001
CMD ["/main"]