FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-http2/build/main /main
COPY examples/using-http2/configs /configs
EXPOSE 9017
CMD ["/main"]