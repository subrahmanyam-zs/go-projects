FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-solr/build/main /main
COPY examples/using-solr/configs /configs
EXPOSE 9099
CMD ["/main"]