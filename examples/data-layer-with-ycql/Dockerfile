FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-ycql/build/main /main
COPY examples/using-ycql/configs /configs
EXPOSE 9005
CMD ["/main"]