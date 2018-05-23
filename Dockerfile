FROM golang:1.10.2-alpine3.7

ENV GIN_MODE debug

ADD . /go/src/github.com/rudeigerc/broker-gateway

WORKDIR /go/src/github.com/rudeigerc/broker-gateway

RUN go build

ENTRYPOINT ["./broker-gateway"]

EXPOSE 8080