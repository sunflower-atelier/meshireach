FROM golang:1.12.5-alpine3.9
COPY ./src/api /go/src/api
WORKDIR /go/src/api
RUN apk update \
  && apk add git \
  && go get -u github.com/golang/dep/cmd/dep \
  && dep ensure

CMD go run main.go

