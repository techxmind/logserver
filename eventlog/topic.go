package eventlog

import (
	"strings"

	pb "github.com/techxmind/logserver/interface-defs"

	"github.com/techxmind/logserver/config"
)

func GetTopic(event *pb.EventLog, cfg *config.TopicRouter) string {
	var (
		routeKey string
	)

	if cfg.RouteMap == nil {
		return cfg.DefaultTopic
	}

	if event.Env != "" {
		routeKey = strings.ToLower(event.AppType + ".env." + event.Env)
		if topic, ok := cfg.RouteMap[routeKey]; ok {
			return topic
		}
	}

	routeKey = strings.ToLower(event.AppType + "." + event.Event)
	if topic, ok := cfg.RouteMap[routeKey]; ok {
		return topic
	}

	routeKey = strings.ToLower(event.AppType)
	if topic, ok := cfg.RouteMap[routeKey]; ok {
		return topic
	}

	return cfg.DefaultTopic
}
