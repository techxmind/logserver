package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pb "github.com/techxmind/logserver/interface-defs"
)

func TestJsonMarshaler(t *testing.T) {
	var (
		sb  = &strings.Builder{}
		jw  MarshalerFunc
		msg = &SinkMessage{
			Event: &pb.EventLog{
				EventId: "my-event-id",
			},
		}
	)

	jw = JSONMarshaler

	if _, err := jw.Marshal(msg); err != nil {
		assert.Errorf(t, err, "JSONMarshaler.Marshal")
	}

	var e pb.EventLog
	if err := json.Unmarshal([]byte(sb.String()), &e); err != nil {
		assert.Error(t, err)
	} else {
		assert.Equal(t, msg.Event.EventId, e.EventId)
	}
}

func TestCSVMarshaler(t *testing.T) {
	var (
		msg = &SinkMessage{
			Event: &pb.EventLog{
				EventId:     "my-event-id",
				UserAgent:   "this is \"useragent\"",
				ScreenWidth: 100,
				ExtendInfo: map[string]string{
					"foo": "bar",
				},
			},
		}
	)

	cw, err := NewCSVMarshaler([]string{"event_id", "user_agent", "screen_width", "extend_info.foo"})
	require.Nil(t, err)

	p, err := cw.Marshal(msg)
	require.Nil(t, err)

	assert.Equal(t, `my-event-id,"this is ""useragent""",100,bar`+"\n", string(p))
}

type testAck struct {
	id  string
	ids []string
}

func (a *testAck) Chain(ack SinkAck) SinkAck {
	dst, _ := ack.(*testAck)
	a.ids = append(a.ids, dst.id)
	return a
}

func TestSink(t *testing.T) {
	var (
		sb = new(strings.Builder)
	)

	cw, err := NewCSVMarshaler([]string{"event_id", "user_agent"})
	require.Nil(t, err)

	sink := NewSink(sb, cw, WithInputBufferSize(1))
	ctx, cancel := context.WithCancel(context.Background())
	ackedIds := make(map[string]int)
	ackCount := 0

	go func() {
		for {
			select {
			case iack := <-sink.Ack():
				ack, _ := iack.(*testAck)
				ackedIds[ack.id] += 1
				ackCount += 1
				for _, id := range ack.ids {
					ackedIds[id] += 1
					ackCount += 1
				}
			case err := <-sink.Errors():
				assert.Error(t, err, "sink error")
			case <-ctx.Done():
				return
			}
		}
	}()

	for i := 0; i < 5; i++ {
		sink.Input() <- &SinkMessage{
			Event: &pb.EventLog{
				EventId:   fmt.Sprintf("event-id-%d", i),
				UserAgent: fmt.Sprintf("ua-%d", i),
			},
			Ack: &testAck{
				id:  fmt.Sprintf("ack-%d", i),
				ids: make([]string, 0),
			},
		}
	}
	// wait for sink processing
	time.Sleep(time.Microsecond * 50)
	sink.Close()
	// wait for the final posiable ack
	time.Sleep(time.Microsecond * 50)
	cancel()
	assert.Equal(t, 5, len(ackedIds))
	assert.Equal(t, 5, ackCount)
	assert.Equal(
		t,
		`event-id-0,ua-0
event-id-1,ua-1
event-id-2,ua-2
event-id-3,ua-3
event-id-4,ua-4
`, sb.String())
}
