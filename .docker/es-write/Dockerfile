FROM alpine:3.12
MAINTAINER Michalski Luc <michalski.luc@gmail.com>

RUN apk add --no-cache --no-progress ca-certificates curl bash nano jq

WORKDIR /opt/service
COPY docker-entrypoint.sh .

CMD ["./docker-entrypoint.sh"]

