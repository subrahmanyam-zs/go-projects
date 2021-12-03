FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/data-layer-with-mongo/build/main /main
COPY examples/data-layer-with-mongo/configs /configs
EXPOSE 9097
CMD ["/main"]