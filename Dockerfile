# SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

FROM golang:1.13

RUN mkdir -p /peridot-jobrunner
WORKDIR /peridot-jobrunner

ADD . /peridot-jobrunner

RUN go get -v ./...
RUN go build
RUN go install github.com/swinslow/peridot-jobrunner
