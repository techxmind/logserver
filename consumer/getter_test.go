package consumer

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pb "github.com/techxmind/logserver/interface-defs"
)

func TestNewStandardName(t *testing.T) {
	ast := assert.New(t)

	tests := []struct {
		name     string
		hasError bool
		standard string
		key      string
	}{
		{"EventId", false, "eventid", ""},
		{"event_id", false, "eventid", ""},
		{"eventId", false, "eventid", ""},
		{"ScreenHeight", false, "screenheight", ""},
		{"ScreenHeight", false, "screenheight", ""},
		{"extend_info", false, "extendinfo", ""},
		{"extend_info.custom", false, "extendinfo", "custom"},
		{"extend_info.custom.other", true, "", ""},
		{"event_id.custom", true, "", ""},
		{"not_exists", true, "", ""},
	}

	for _, test := range tests {
		s, err := newStandardName(test.name)
		hasError := err != nil
		if test.hasError != hasError {
			expect := "no err"
			if test.hasError {
				expect = "err"
			}
			t.Errorf("%s expect %s", test.name, expect)
		}
		if err == nil {
			ast.Equal(test.standard, s.standard)
			ast.Equal(test.key, s.key)
		}
	}
}

func TestGetValue(t *testing.T) {
	ast := assert.New(t)

	now := time.Now()
	e := &pb.EventLog{
		EventId:      "my_event_id",
		EventTime:    now.UnixNano() / int64(1e6),
		LoggedTime:   now.UnixNano() / int64(1e6),
		ScreenHeight: 300,
		Network:      1,
		ExtendInfo: map[string]string{
			"foo": "bar",
		},
	}

	var jsonEncode = func(obj interface{}) string {
		v, _ := json.Marshal(obj)
		return string(v)
	}

	tests := []struct {
		name   string
		expect string
	}{
		{"event_id", e.EventId},
		{"not_exists", ""},
		{"screen_height", strconv.Itoa(int(e.ScreenHeight))},
		{"network", strconv.Itoa(int(e.Network))},
		{"extend_info", jsonEncode(e.ExtendInfo)},
		{"extend_info.foo", e.ExtendInfo["foo"]},
		{"event_time", strconv.Itoa(int(now.UnixNano() / int64(1e6)))},
		{"event_time_str", now.Format("2006-01-02T15:04:05Z")},
		{"logged_time_str", now.Format("2006-01-02T15:04:05Z")},
	}

	for _, test := range tests {
		if name, err := newStandardName(test.name); err != nil {
			ast.Errorf(err, "newStandardName err:%s", test.name)
		} else {
			ast.Equal(test.expect, name.MustGetValue(e), "name="+test.name)
		}
	}
}
