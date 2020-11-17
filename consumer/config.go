package consumer

import (
	"flag"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/techxmind/logserver/logger"
)

var DefaultConfig *Config

func init() {
	DefaultConfig = &Config{
		GroupID: "default",
		Sink: &SinkConfig{
			Marshaler: "json",
			Target:    "stdout",
		},
	}

	flag.StringVar(
		&DefaultConfig.Addrs,
		"addrs",
		"127.0.0.1:9092",
		"Kafka broker addresses, multiple values are comma separated",
	)
	flag.StringVar(
		&DefaultConfig.GroupID,
		"group_id",
		"default",
		"Kafka consumer group id",
	)
	flag.StringVar(
		&DefaultConfig.KafkaVersion,
		"kafka_version",
		"",
		"Kafka version",
	)
	flag.StringVar(
		&DefaultConfig.Topics,
		"topics",
		"event-log",
		"Topics to consume, multiple values are comma separated",
	)
	flag.StringVar(
		&DefaultConfig.Offset,
		"offset",
		"newest",
		"initial offset: newest | oldest",
	)
	flag.StringVar(
		&DefaultConfig.Sink.Marshaler,
		"sink.marshaler",
		"json",
		"Output marshaler, json|csv",
	)
	flag.StringVar(
		&DefaultConfig.Sink.MarshalerArgs,
		"sink.marshaler_args",
		"",
		"Output marshaler args, csv marshaler requires headers definitions. Header names are comma separated",
	)
	flag.StringVar(
		&DefaultConfig.Sink.Target,
		"sink.target",
		"stdout",
		"Output target, stdout|file",
	)
	flag.StringVar(
		&DefaultConfig.Sink.TargetArgs,
		"sink.target_args",
		"",
		"Output target args, file target requires filename specified. RollingFile format: filename:maxsize[:maxage:suffix]",
	)
	flag.IntVar(
		&DefaultConfig.Sink.InputBufferSize,
		"sink.input_buffer_size",
		100,
		"Sink input buffer size",
	)
	flag.IntVar(
		&DefaultConfig.Sink.OutputBufferSize,
		"sink.output_buffer_size",
		4096,
		"Sink output buffer size, 0 means no output buffer",
	)
}

type Config struct {
	GroupID      string      `json:"group_id"`
	KafkaVersion string      `json:"kafka_version"`
	Addrs        string      `json:"addrs"`
	Topics       string      `json:"topics"`
	Offset       string      `json:"offset"`
	Sink         *SinkConfig `json:"sink"`
}

type SinkConfig struct {
	Marshaler        string `json:"marshaler"`
	MarshalerArgs    string `json:"marshaler_args"`
	Target           string `json:"target"`
	TargetArgs       string `json:"target_args"`
	OutputBufferSize int    `json:"output_buffer_size"`
	InputBufferSize  int    `json:"input_buffer_size"`
}

func (cfg *Config) GetKafkaVersion() sarama.KafkaVersion {
	v, err := sarama.ParseKafkaVersion(cfg.KafkaVersion)
	if err != nil {
		logger.Fatalf("invalid kafka version:%s", cfg.KafkaVersion)
	}
	return v
}

func (cfg *Config) GetTopics() []string {
	return splitValues(cfg.Topics)
}

func (cfg *Config) GetAddrs() []string {
	return splitValues(cfg.Addrs)
}

func splitValues(str string) []string {
	var vals = make([]string, 0, 1)

	for _, val := range strings.Split(str, ",") {
		val = strings.TrimSpace(val)
		if val != "" {
			vals = append(vals, val)
		}
	}

	return vals
}
