package config

import (
	"fmt"
	"net/http"

	"github.com/Shopify/sarama"
)

type Config struct {
	Version     string
	VersionDate string
	HTTPAddr    string         `json:"http_addr"`
	DebugAddr   string         `json:"debug_addr"`
	GRPCAddr    string         `json:"grpc_addr"`
	TopicRouter *TopicRouter   `json:"topic_router,omitempty"`
	Storage     *StorageConfig `json:"storage,omitempty"`
}

func (c *Config) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Version:%s\nDate:%s", c.Version, c.VersionDate)
}

type TopicRouter struct {
	DefaultTopic string `json:"default_topic"`
	// router defination and priority:
	//   {app_type}.env.{env} => topic
	//   {app_type}.{event} => topic
	//   {app_type} => topic
	//
	//   e.g.
	//   routemap:
	//   myapp.env.test => myapp_event_log_test
	//   myapp.pv => myapp_event_log_pv
	//   myapp => myapp_event_log_other
	//
	//   case app_type = myapp, env = test then myapp_event_log_test
	//   case app_type = myapp, env != test, event = pv then myapp_event_log_pv
	//   case app_type = myapp, env != test, event = mv then myapp_event_log_other
	RouteMap map[string]string `json:"route_map"`
}

type StorageConfig struct {
	DataType string       `json:"data_type"`       //json|protobuf, default protobuf
	Types    string       `json:"types,omitempty"` //kafka|stdout, multiple values are comma separated, default stdout
	Kafka    *KafkaConfig `json:"kafka,omitempty"`
}

type KafkaConfig struct {
	Addrs          string         `json:"addrs"`
	ProducerConfig *sarama.Config `json:"producer_config"`
}
