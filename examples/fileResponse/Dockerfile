FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/fileResponse/build/main /main
COPY examples/fileResponse/configs /configs
EXPOSE 8080
CMD ["/main"]