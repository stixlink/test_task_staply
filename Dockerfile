FROM golang:alpine
# specify the token in order to fetch the sources
ARG GITHUB_TOKEN

ADD ./ /go/src/guthub.com/stixlink/test_task_staply

WORKDIR /go/src/guthub.com/stixlink/test_task_staply

ENV GOPATH=/go

RUN apk add --no-cache --virtual .build-deps git && \
    git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/" && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep ensure -v

CMD ["go", "run", "*.go"]
