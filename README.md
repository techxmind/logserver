# LogServer [![Build Status](https://travis-ci.com/techxmind/logserver.svg?branch=main)](https://travis-ci.org/techxmind/logserver)

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

## Client

Support both gGRPC and HTTP protocol.

HTTP example:
Post single event, can use json or protobuf data format.
See interface-defs/event_log.proto `EventLog` for data structure.
```
curl -H "Content-Type: application/json" --data-binary @single-event.json http://logserver-host/s
curl -H "Content-Type: application/protobuf" --data-binary @single-event.pb http://logserver-host/s
```

Post multiple events, can use json or protobuf data format.
See interface-defs/event_log.proto `EventLogs` for data structure.
```
curl -H "Content-Type: application/json" --data-binary @multiple-events.json http://logserver-host/mul
curl -H "Content-Type: application/protobuf" --data-binary @multiple-events.pb http://logserver-host/mul
```
