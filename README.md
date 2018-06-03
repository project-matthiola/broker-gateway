# broker-gateway

[![Build Status](https://travis-ci.com/rudeigerc/broker-gateway.svg?token=m9esAaP4YUBsZ2yN5xJq&branch=master)](https://travis-ci.com/rudeigerc/broker-gateway)

The broker gateway of project **Matthiola**, a distributed commodities OTC electronic trading system, instructed by Morgan Stanley. 

## Architecture

- Receiver
- Matcher
- Server (HTTP Server)
- Broadcaster (WebSocket Server)

## Quick Start

### Docker Compose

```bash
$ docker-compose up
```

## Build

```bash
$ brew install dep
$ dep ensure
$ go build
$ ./broker-gateway --help
```

## Run

```bash
Usage:
  broker-gateway [command]

Available Commands:
  broadcaster Run WebSocket server
  help        Help about any command
  matcher     Run matcher
  receiver    Run receiver
  sender      Run sender
  server      Run HTTP server

Flags:
  -c, --config string   config file (default "config/config.toml")
  -h, --help            help for broker-gateway

Use "broker-gateway [command] --help" for more information about a command.
```

## Config

See `config/config.toml`.

## Requirement

### Microservices

#### Micro

[Micro](https://github.com/micro/micro) is a toolkit for cloud-native development. It helps you build future-proof application platforms and services.

```bash
$ go get -u github.com/micro/micro
```

### Service Discovery

#### Consul

[Consul](https://github.com/hashicorp/consul) is a tool for service discovery and configuration. Consul is distributed, highly available, and extremely scalable.

```bash
$ brew install consul
$ consul agent -dev
```

### API Gateway

#### Go API

[Go API](https://github.com/micro/go-api) is a pluggable API framework.

It builds on go-micro and includes a set of packages for composing HTTP based APIs.

```bash
# The HTTP handler with web socket support included
$ micro api --namespace=github.com.rudeigerc.broker-gateway --handler=web
```

- HTTP Server `/server`
- WebSocket Server `/broadcaster`

### Database

#### MySQL

```bash
$ brew install mysql
$ brew services start mysql
```

### Store

#### etcd

[etcd](https://github.com/coreos/etcd) is a distributed reliable key-value store for the most critical data of a distributed system.

```bash
$ brew install etcd
$ brew services start etcd
```

### Message Queue

#### NSQ

[NSQ](https://nsq.io/) is a realtime distributed messaging platform.

```bash
$ brew install nsq
```

- In one shell, start `nsqlookupd`:

```bash
$ nsqlookupd --broadcast-address=127.0.0.1
```

- In another shell, start `nsqd`:

```bash
$ nsqd --lookupd-tcp-address=127.0.0.1:4160 --broadcast-address=127.0.0.1
```

- In another shell, start `nsqadmin`:

```bash
$ nsqadmin --lookupd-http-address=127.0.0.1:4161
```

- In a web browser open [`http://127.0.0.1:4171`](http://127.0.0.1:4171) to view the nsqadmin UI and see statistics.

## Docs

See [docs](https://github.com/project-matthiola/docs) and [api-docs](https://github.com/project-matthiola/api-docs).

## License

MIT