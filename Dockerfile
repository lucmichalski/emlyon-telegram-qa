FROM golang:1-alpine AS builder
MAINTAINER Michalski Luc <michalski.luc@gmail.com>

ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

# NB. CGO enabled for sqlite3 driver
ARG CGO=1
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/lucmichalski/emlyon-telegram-qa

COPY . /go/src/github.com/lucmichalski/emlyon-telegram-qa

# gcc/g++ are required to build SASS libraries for extended version
RUN apk update && \
    apk add --no-cache git ca-certificates

RUN go build -ldflags "-extldflags=-static -extldflags=-lm" -o /go/bin/emlyon-telegram-qa

FROM alpine:3.12

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/emlyon-telegram-qa /usr/bin/emlyon-telegram-qa

ENTRYPOINT ["emlyon-telegram-qa"]
