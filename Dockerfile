FROM golang:latest
RUN apt-get update
RUN apt-get install -y git autoconf automake libtool curl make g++ unzip

#RUN go get -u github.com/golang/dep/...
#RUN mkdir ~/.ssh;\
#    touch ~/.ssh/known_hosts;\
#    ssh-keyscan gitlab.com >> ~/.ssh/known_hosts;

RUN mkdir -p $GOPATH/src/github.com/partyzanex/golang-test-task
COPY ./ $GOPATH/src/github.com/partyzanex/golang-test-task
WORKDIR $GOPATH/src/github.com/partyzanex/golang-test-task

RUN	make gotask

FROM alpine

COPY --from=0 /go/src/github.com/partyzanex/golang-test-task/crawl ./crawl

EXPOSE 3030