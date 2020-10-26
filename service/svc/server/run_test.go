package server

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/techxmind/go-utils/stringutil"
	"github.com/techxmind/logserver/config"
	pb "github.com/techxmind/logserver/interface-defs"
	"github.com/techxmind/logserver/service/handlers"
	"github.com/techxmind/logserver/service/svc"
	"github.com/techxmind/logserver/storage"
)

var (
	_testStorage *storage.MockStorage
)

func TestMain(m *testing.M) {
	_testStorage, _ = storage.NewMockStorage(100)
	handlers.StorageGet = testStorage
	config.DefaultConfig.TopicRouter = &config.TopicRouter{
		DefaultTopic: "test-event",
		RouteMap: map[string]string{
			"myapp":     "test-event1",
			"myapp1.pv": "test-event2",
		},
	}
	config.DefaultConfig.Storage.DataType = "json"
	ret := m.Run()
	os.Exit(ret)
}

func testStorage(cfg *config.StorageConfig) (storage.Storager, error) {
	return _testStorage, nil
}

func TestHttpJSONRequest(t *testing.T) {
	ast := assert.New(t)

	service := handlers.NewService()
	endpoints := NewEndpoints(service)
	handler := svc.MakeHTTPHandler(endpoints)

	eventLog := testEventLog()
	postData, _ := json.Marshal(eventLog)
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/s", bytes.NewReader(postData))
	request.Header.Add("Accept", "test/plain")
	request.Header.Add("X-Forwarded-For", "192.168.0.101, 124.78.41.83") // 上海IP
	request.Header.Add("User-Agent", "test-ua")
	request.Header.Add("Referer", "http://techxmind.com")
	request.Header.Add("Content-Length", strconv.Itoa(len(postData)))

	handler.ServeHTTP(writer, request)

	ast.Equal("{}", writer.Body.String())

	msg := <-_testStorage.Successes()
	ast.Equal("test-event1", msg.Topic)
	v, _ := msg.Value.Marshal()
	t.Log(string(v))
}

func TestHttpProtobufRequest(t *testing.T) {
	ast := assert.New(t)

	service := handlers.NewService()
	endpoints := NewEndpoints(service)
	handler := svc.MakeHTTPHandler(endpoints)

	eventLog := testEventLog()
	eventLog.AppType = "myapp2"
	postData, _ := eventLog.Marshal()
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/s", bytes.NewReader(postData))
	request.Header.Add("Accept", "test/plain")
	request.Header.Add("X-Forwarded-For", "192.168.0.101, 124.78.41.83") // 上海IP
	request.Header.Add("User-Agent", "test-ua")
	request.Header.Add("Referer", "http://techxmind.com")
	request.Header.Add("Content-Type", "application/protobuf")
	request.Header.Add("Content-Length", strconv.Itoa(len(postData)))

	handler.ServeHTTP(writer, request)

	ast.Equal("{}", writer.Body.String())

	msg := <-_testStorage.Successes()
	ast.Equal("test-event", msg.Topic)
}

func TestHttpMulJSONRequest(t *testing.T) {
	ast := assert.New(t)

	service := handlers.NewService()
	endpoints := NewEndpoints(service)
	handler := svc.MakeHTTPHandler(endpoints)

	eventLogs := testEventLogs(3)
	postData, _ := json.Marshal(eventLogs)
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/mul", bytes.NewReader(postData))
	request.Header.Add("Accept", "test/plain")
	request.Header.Add("X-Forwarded-For", "192.168.0.101, 124.78.41.83") // 上海IP
	request.Header.Add("User-Agent", "test-ua")
	request.Header.Add("Referer", "http://techxmind.com")
	request.Header.Add("Content-Length", strconv.Itoa(len(postData)))

	handler.ServeHTTP(writer, request)

	ast.Equal("{}", writer.Body.String())

	msg := <-_testStorage.Successes()
	ast.Equal("test-event2", msg.Topic)
	msg = <-_testStorage.Successes()
	ast.Equal("test-event2", msg.Topic)
	msg = <-_testStorage.Successes()
	ast.Equal("test-event2", msg.Topic)
}

func testEventLogs(n int) *pb.EventLogs {
	logs := &pb.EventLogs{
		Common: testEventLogCommon(),
		Events: make([]*pb.EventLog, 0),
	}

	for i := 0; i < n; i++ {
		logs.Events = append(logs.Events, testSimpleEventLog())
	}

	return logs
}

func testEventLog() *pb.EventLog {
	return &pb.EventLog{
		EventId:    stringutil.Rand(10),
		EventTime:  time.Now().UnixNano() / int64(1e6),
		SessionId:  "s-" + stringutil.Rand(5),
		Udid:       "udid-" + stringutil.Rand(10),
		Tkid:       "tkid-" + stringutil.Rand(10),
		Mid:        "mid-" + stringutil.Rand(10),
		Platform:   "server",
		AppVersion: "1.0.0",
		AppChannel: "yyb",
		AppType:    "myapp",
		Event:      "PV",
		// device info
		Os:               "ios",
		OsVersion:        "10.1.2",
		DeviceModel:      "iphone",
		DeviceVendor:     "apple",
		DeviceBrand:      "apple",
		ScreenSize:       "900x768",
		ScreenWidth:      900,
		ScreenHeight:     768,
		ScreenResolution: "1",
		Imei:             "",
		AndroidId:        "",
		Idfa:             "idfa-xxxxx",
		Oaid:             "",
		Carrier:          0,
		Network:          0,
		Lon:              "lon-val",
		Lat:              "lat-val",
		Mac:              "mac-val",
		// event info
		PageId:     "order",
		PvId:       stringutil.Rand(10),
		LayoutId:   "",
		PageKey:    "100",
		ModuleId:   "",
		ActionType: "refresh",
		// ref event info
		RefPageId:   "order_list",
		RefPvId:     stringutil.Rand(10),
		RefLayoutId: "",
		RefPageKey:  "",
		RefModuleId: "list",
		// event data
		ExtendInfo: map[string]string{"foo": "bar"},
	}
}

func testSimpleEventLog() *pb.EventLog {
	return &pb.EventLog{
		EventId:   stringutil.Rand(10),
		EventTime: time.Now().UnixNano() / int64(1e6),
		SessionId: "s-" + stringutil.Rand(5),
		Event:     "PV",
		// event info
		PageId:     "order",
		PvId:       "unique-page-visit-id",
		LayoutId:   "",
		PageKey:    "{order_Id}",
		ModuleId:   "",
		ActionType: "refresh",
		// ref event info
		RefPageId:   "order_list",
		RefPvId:     "{order_list_pvid}",
		RefLayoutId: "",
		RefPageKey:  "",
		RefModuleId: "list",
		// event data
		ExtendInfo: map[string]string{"foo": "bar"},
	}
}

func testEventLogCommon() *pb.EventLogCommon {
	return &pb.EventLogCommon{
		Udid:             "udid-" + stringutil.Rand(10),
		Tkid:             "tkid-" + stringutil.Rand(10),
		Mid:              "mid-" + stringutil.Rand(10),
		Platform:         "android",
		AppVersion:       "2.0",
		AppChannel:       "xiaomi",
		AppType:          "myapp1",
		Os:               "android",
		OsVersion:        "8.0",
		DeviceModel:      "Redmi",
		DeviceVendor:     "XiaoMi",
		DeviceBrand:      "XiaoMi",
		ScreenSize:       "800x600",
		ScreenWidth:      800,
		ScreenHeight:     600,
		ScreenResolution: "2",
		Imei:             "imei-" + stringutil.Rand(10),
		AndroidId:        "androiid-" + stringutil.Rand(10),
		Idfa:             "",
		Oaid:             "oaid-" + stringutil.Rand(10),
		Carrier:          1,
		Network:          1,
		Lon:              "lan-2",
		Lat:              "lat-2",
		Mac:              "mac-2",
	}
}
