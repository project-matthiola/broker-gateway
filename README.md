# broker-gateway

The broker gateway of project **Matthiola**, a distributed commodities OTC electronic trading system, instructed by Morgan Stanley. 

## Architecture

- Receiver
- Matcher
- Executor
- Server
  - Router
  - Service
  - Mapper

## Build

```shell
$ brew install dep
$ dep ensure
$ go build
$ ./broker-gateway
```

### NSQ

[NSQ](https://nsq.io/) is a realtime distributed messaging platform.

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

- In a web browser open [`http://127.0.0.1:4171/`](http://127.0.0.1:4171/) to view the nsqadmin UI and see statistics.

## Run

```shell
Usage:
  broker-gateway [command]

Available Commands:
  help        Help about any command
  matcher     Run matcher
  receiver    Run receiver
  sender      Run sender
  server      Run server

Flags:
  -c, --config string   config file (default "config/config.toml")
  -h, --help            help for broker-gateway

Use "broker-gateway [command] --help" for more information about a command.
```

## License

MIT