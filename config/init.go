package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	DefaultConfig *Config
)

func init() {
	DefaultConfig = &Config{
		TopicRouter: &TopicRouter{
			DefaultTopic: "event-log",
		},
		Storage: &StorageConfig{
			DataType: "protobuf",
			Types:    "stdout",
			Kafka:    &KafkaConfig{},
		},
	}

	flag.StringVar(&DefaultConfig.DebugAddr, "debug.addr", ":5060", "Debug and metrics listen address")
	flag.StringVar(&DefaultConfig.HTTPAddr, "http.addr", ":5050", "HTTP listen address")
	flag.StringVar(&DefaultConfig.GRPCAddr, "grpc.addr", ":5040", "gRPC (HTTP) listen address")

	flag.StringVar(
		&DefaultConfig.Storage.DataType,
		"storage.data_type",
		"",
		"protobuf|json. Data type that store in storage",
	)
	flag.StringVar(
		&DefaultConfig.Storage.Kafka.Addrs,
		"storage.kafka.addrs",
		"",
		"Kafka broker addresses, multiple values are comma separated",
	)
	flag.StringVar(
		&DefaultConfig.Storage.Types,
		"storage.types",
		"",
		"stdout|kafka. Multiple values are comma separated",
	)

	// Use environment variables, if set. Flags have priority over Env vars.
	if addr := os.Getenv("DEBUG_ADDR"); addr != "" {
		DefaultConfig.DebugAddr = addr
	}
	if port := os.Getenv("PORT"); port != "" {
		DefaultConfig.HTTPAddr = fmt.Sprintf(":%s", port)
	}
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		DefaultConfig.HTTPAddr = addr
	}
	if addr := os.Getenv("GRPC_ADDR"); addr != "" {
		DefaultConfig.GRPCAddr = addr
	}
	if dataType := os.Getenv("STORAGE_DATA_TYPE"); dataType != "" {
		DefaultConfig.Storage.DataType = dataType
	}
	if types := os.Getenv("STORAGE_TYPES"); types != "" {
		DefaultConfig.Storage.Types = types
	}
	if addrs := os.Getenv("STORAGE_KAFKA_ADDRS"); addrs != "" {
		DefaultConfig.Storage.Kafka.Addrs = addrs
	}
}
