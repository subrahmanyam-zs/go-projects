FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-mongo/build/main /main
COPY examples/using-mongo/configs /configs
EXPOSE 9097
CMD ["/main"]