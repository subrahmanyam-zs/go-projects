FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-cassandra/build/main /main
COPY examples/using-cassandra/configs /configs
EXPOSE 9094
CMD ["/main"]