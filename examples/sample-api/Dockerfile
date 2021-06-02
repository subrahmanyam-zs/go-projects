FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/sample-api/build/main /main
COPY examples/sample-api/configs /configs
EXPOSE 9000
CMD ["/main"]