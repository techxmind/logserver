package handlers

import (
	"context"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"

	"github.com/techxmind/logserver/config"
	"github.com/techxmind/logserver/eventlog"
	pb "github.com/techxmind/logserver/interface-defs"
	"github.com/techxmind/logserver/logger"
	"github.com/techxmind/logserver/metrics"
	"github.com/techxmind/logserver/storage"
)

var (
	_ = fmt.Sprint

	_service *logserviceService
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.LogServiceServer {

	cfg := config.DefaultConfig

	storage, err := StorageGet(cfg.Storage)

	if err != nil {
		logger.Fatal("storage init failed", "err", err)
	}

	_service = &logserviceService{
		storage: storage,
		config:  cfg,
	}

	return _service
}

type logserviceService struct {
	storage storage.Storager
	config  *config.Config
}

func (s logserviceService) SubmitSingle(ctx context.Context, in *pb.EventLog) (*pb.Response, error) {
	var resp pb.Response

	if err := validate(in); err != nil {
		return nil, err
	}

	eventlog.FillSingle(ctx, in)

	s.writeEventLog(in)

	return &resp, nil
}

func (s logserviceService) SubmitMultiple(ctx context.Context, in *pb.EventLogs) (*pb.Response, error) {
	var resp pb.Response

	if len(in.Events) == 0 {
		return nil, errors.New("No event")
	}

	eventlog.Fill(ctx, in.Common, in.Events)

	for _, event := range in.Events {
		if err := validate(event); err != nil {
			continue
		}

		s.writeEventLog(event)
	}

	return &resp, nil
}

func (s logserviceService) writeEventLog(event *pb.EventLog) {
	msg := storage.NewMessage(
		eventlog.GetTopic(event, s.config.TopicRouter),
		event.EventId,
		event,
	)

	//val, err := json.MarshalIndent(in, "", "  ")

	// Ignore it when topic is empty
	if msg.Topic != "" {
		if s.config.Storage != nil && s.config.Storage.DataType == "json" {
			msg.Value = storage.JSONMarshaler(msg.Value)
		}
		if err := s.storage.Write(msg); err != nil {
			metrics.CounterAdd("record_write_err", 1)
			logger.Error("write storage", "err", err)
		}
	}

	metrics.CountEvent(msg.Topic, event.AppType, event.Event)
}

func (s logserviceService) close() {
	if s.storage != nil {
		s.storage.Close()
	}
}

func (s logserviceService) Ping(ctx context.Context, in *pb.Empty) (*pb.Response, error) {
	var resp pb.Response
	return &resp, nil
}
