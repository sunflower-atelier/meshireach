FROM golang:1.12.5-alpine3.9
COPY ./ /go/meshireach
WORKDIR /go/meshireach
RUN apk update \
  && apk add --nocache git\
  && go get github.com/pilu/fresh \
  && go get -u golang.org/x/lint/golint \
  && go mod download
CMD ["fresh"]

