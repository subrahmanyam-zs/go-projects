FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-http-service/build/main /main
COPY examples/using-http-service/configs /configs
EXPOSE 9091
CMD ["/main"]