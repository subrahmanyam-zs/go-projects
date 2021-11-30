FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/sample-https/build/main /main
COPY examples/sample-https/configs /configs
EXPOSE 1443
EXPOSE 9007
CMD ["/main"]