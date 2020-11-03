package main

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
	consumer, err := consumer.New(ctx, consumer.DefaultConfig)
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
