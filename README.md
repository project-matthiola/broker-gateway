# broker-gateway

[![Build Status](https://travis-ci.com/rudeigerc/broker-gateway.svg?token=m9esAaP4YUBsZ2yN5xJq&branch=master)](https://travis-ci.com/rudeigerc/broker-gateway)

The broker gateway of project **Matthiola**, a distributed commodities OTC electronic trading system, instructed by Morgan Stanley. 

## Architecture

- Receiver
- Matcher
- Server (HTTP Server)
- Broadcaster (WebSocket Server)

## Build

```shell
$ brew install dep
$ dep ensure
$ go build
$ ./broker-gateway --help
```

## Run

```shell
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

### Database

#### MySQL

```shell
$ brew install mysql
$ brew services start mysql
```

### Store

#### etcd

[etcd](https://github.com/coreos/etcd) is a distributed reliable key-value store for the most critical data of a distributed system.

```shell
$ brew install etcd
$ brew services start etcd
```

### Message Queue

#### NSQ

[NSQ](https://nsq.io/) is a realtime distributed messaging platform.

```shell
$ brew install nsq
```

- In one shell, start `nsqlookupd`:

```shell
$ nsqlookupd --broadcast-address=127.0.0.1
```

- In another shell, start `nsqd`:

```shell
$ nsqd --lookupd-tcp-address=127.0.0.1:4160 --broadcast-address=127.0.0.1
```

- In another shell, start `nsqadmin`:

```shell
$ nsqadmin --lookupd-http-address=127.0.0.1:4161
```

- In a web browser open [`http://127.0.0.1:4171`](http://127.0.0.1:4171) to view the nsqadmin UI and see statistics.

### Service Discovery

#### Consul

```shell
$ brew install consul
$ consul agent -dev
```

## Docs

See [docs](https://github.com/project-matthiola/docs) and [api-docs](https://github.com/project-matthiola/api-docs).

## License

MIT