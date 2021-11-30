FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/sample-grpc/build/main /main
COPY examples/sample-grpc/configs /configs
EXPOSE 9093
CMD ["/main"]