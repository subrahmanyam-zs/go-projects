FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates
COPY examples/using-postgres/build/main /main
COPY examples/using-postgres/configs /configs
EXPOSE 9092
CMD ["/main"]