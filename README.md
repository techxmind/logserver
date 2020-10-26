# LogServer [![Build Status](https://travis-ci.org/techxmind/logserver.svg?branch=main)](https://travis-ci.org/techxmind/logserver)

User behavior event log collection service.

[event_log.proto](https://github.com/techxmind/logserver/blob/master/interface-defs/event_log.proto)

## Usage
```
go build service/cmd/logservice/main.go

logservice -option
Usage
  -debug.addr string
    	Debug and metrics listen address (default ":5060")
  -grpc.addr string
    	gRPC (HTTP) listen address (default ":5040")
  -http.addr string
    	HTTP listen address (default ":5050")
  -storage.data_type string
    	protobuf|json. Data type that store in storage
  -storage.kafka.addrs string
    	Kafka broker addresses, multiple values are comma separated
  -storage.types string
    	stdout|kafka. Multiple values are comma separated
```

## Custom
```
package main

import (
	"flag"

	"github.com/techxmind/logserver/config"
	"github.com/techxmind/logserver/service/svc/server"
)

func main() {
	flag.Parse()

    // init config or load config from file
    cfg := &config.Config{}
	server.Run(cfg)
}
```
