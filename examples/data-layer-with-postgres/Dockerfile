FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/data-layer-with-postgres/build/main /main
COPY examples/data-layer-with-postgres/configs /configs
EXPOSE 9092
CMD ["/main"]