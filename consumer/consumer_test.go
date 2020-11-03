package consumer_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/techxmind/logserver/consumer"
)

func ExampleNew() {
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
