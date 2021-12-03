FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/data-layer-with-ycql/build/main /main
COPY examples/data-layer-with-ycql/configs /configs
EXPOSE 9005
CMD ["/main"]