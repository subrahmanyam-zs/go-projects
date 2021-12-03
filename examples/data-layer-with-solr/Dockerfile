FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/data-layer-with-solr/build/main /main
COPY examples/data-layer-with-solr/configs /configs
EXPOSE 9099
CMD ["/main"]