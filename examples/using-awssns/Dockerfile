FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-awssns/build/main /main
COPY examples/using-awssns/configs /configs
EXPOSE 8080
CMD ["/main"]