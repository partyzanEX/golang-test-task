FROM golang:latest
RUN apt-get update
RUN apt-get install -y git libtool make g++ unzip

RUN mkdir -p $GOPATH/src/github.com/partyzanex/golang-test-task
COPY ./ $GOPATH/src/github.com/partyzanex/golang-test-task
WORKDIR $GOPATH/src/github.com/partyzanex/golang-test-task

RUN	make gotask

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=0 /go/src/github.com/partyzanex/golang-test-task/gotask ./gotask
COPY --from=0 /go/src/github.com/partyzanex/golang-test-task/config.json ./config.json

EXPOSE 3030