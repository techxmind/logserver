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

var (
	// version is compiled into logservice with the flag
	// go install -ldflags "-X main.version=$VERSION"
	version string

	// BuildDate is compiled into logservice with the flag
	// go install -ldflags "-X main.date=$VERSION_DATE"
	date string
)

func main() {
	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version:%s\nDate:%s", version, date)
		return
	}

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
