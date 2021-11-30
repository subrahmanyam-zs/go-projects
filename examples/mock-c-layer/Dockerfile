FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/mock-c-layer/build/main /main
COPY examples/mock-c-layer/configs /configs
EXPOSE 9090
CMD ["/main"]