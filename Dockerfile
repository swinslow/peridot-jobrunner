FROM golang:1.12

RUN mkdir -p /peridot-jobrunner
WORKDIR /peridot-jobrunner

ADD . /peridot-jobrunner

RUN go get -v ./...
RUN go build
RUN go install github.com/swinslow/peridot-jobrunner
