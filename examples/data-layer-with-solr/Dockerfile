FROM alpine:latest

RUN mkdir -p /src/build
WORKDIR  /src/build

RUN apk add --no-cache tzdata ca-certificates

COPY /build/main /main
COPY /configs /configs

EXPOSE 9099
CMD ["/main"]