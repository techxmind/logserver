package eventlog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/techxmind/logserver/config"
	pb "github.com/techxmind/logserver/interface-defs"
)

func TestGetTopic(t *testing.T) {
	topicRouter := &config.TopicRouter{
		DefaultTopic: "default",
		RouteMap: map[string]string{
			"myapp.env.test": "myapp_event_test",
			"myapp.pv":       "myapp_event_pv",
			"myapp":          "myapp_event_default",
		},
	}
	ev := func(appType, env, event string) *pb.EventLog {
		return &pb.EventLog{
			AppType: appType,
			Env:     env,
			Event:   event,
		}
	}
	tests := []struct {
		event       *pb.EventLog
		expectTopic string
	}{
		{
			ev("myapp", "test", "pv"),
			"myapp_event_test",
		},
		{
			ev("myapp", "prod", "pv"),
			"myapp_event_pv",
		},
		{
			ev("myapp", "", "pv"),
			"myapp_event_pv",
		},
		{
			ev("myapp2", "", "pv"),
			"default",
		},
		{
			ev("myapp", "", "mv"),
			"myapp_event_default",
		},
	}

	ast := assert.New(t)
	for _, test := range tests {
		ast.Equal(test.expectTopic, GetTopic(test.event, topicRouter))
	}
}
