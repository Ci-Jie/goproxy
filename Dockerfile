FROM golang:1.17.3

WORKDIR /goproxy
COPY . /goproxy
RUN CGO_ENABLED=0 make

FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk add bash

WORKDIR /root
COPY --from=0 /goproxy/goproxy /usr/bin/
COPY etc/goproxy.yaml /etc/goproxy/goproxy.yaml
CMD ["goproxy", "start"]
