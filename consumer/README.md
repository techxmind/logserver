# EvnetLog Consumer

## CMD tool
```
go get github.com/techxmind/logserver/consumer/cmd/logconsumer

logconsumer --help

Usage:
  -addrs string
    	Kafka broker addresses, multiple values are comma separated (default "127.0.0.1:9092")
  -group_id string
    	Kafka consumer group id (default "default")
  -kafka_version string
    	Kafka version
  -log-level value
    	minimum enabled logging level. debug|info|warn|error|dpanic|panic|fatal
  -sink.input_buffer_size int
    	Sink input buffer size (default 100)
  -sink.marshaler string
    	Output marshaler, json|csv (default "json")
  -sink.marshaler_args string
    	Output marshaler args, csv marshaler requires headers definitions. Header names are comma separated
  -sink.output_buffer_size int
    	Sink output buffer size, 0 means no output buffer (default 4096)
  -sink.target string
    	Output target, stdout|file (default "stdout")
  -sink.target_args string
    	Output target args, file target requires filename specified. RollingFile format: filename:maxsize[:maxage:suffix]
  -topics string
    	Topics to consume, multiple values are comma separated (default "event-log")
```

## CUSTOM
```
import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/techxmind/logserver/consumer"
)

func main() {
	flag.Parse()

	ctx := context.Background()

    cfg := &consumer.Config{
		Addrs:        "localhost:9092",
		KafkaVersion: "2.6.0",
		GroupID:      "test",
		Topics:       "event_log",
		Sink: &consumer.SinkConfig{
			Marshaler:  "json",
			Target:     "file",
			TargetArgs: "event_log.json:1024000:300", //rollingfile filename:maxsize:maxage
		},
	}

	consumer, err := consumer.New(ctx, cfg)
	if err != nil {
		fmt.Printf("consumer.New err:%s\n", err)
		return
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		consumer.Close()
	}()

	consumer.Start()
}
```
