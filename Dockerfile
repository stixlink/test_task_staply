FROM golang:latest as builder
COPY . $GOPATH/src/github.com/stixlink/test_task_staply/
WORKDIR $GOPATH/src/github.com/stixlink/test_task_staply/

RUN apt-get update
RUN apt-get install libmagickwand-dev -y
RUN apt-get install imagemagick -y
#get dependancies
RUN go get -d -v
#build the binary
RUN go build -o /go/bin/test_task_staply
RUN chmod +x /go/bin/test_task_staply

ENTRYPOINT ["/go/bin/test_task_staply"]

