package kafka

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"

	"github.com/techxmind/logserver/config"
	"github.com/techxmind/logserver/logger"
	"github.com/techxmind/logserver/metrics"
	"github.com/techxmind/logserver/storage"
)

var (
	_ = fmt.Println
)

type kafkaStorage struct {
	quit     chan bool
	producer sarama.AsyncProducer
}

func New(cfg *config.KafkaConfig) (storage.Storager, error) {
	var addrs = make([]string, 0, 1)

	for _, addr := range strings.Split(cfg.Addrs, ",") {
		addr = strings.TrimSpace(addr)
		if addr != "" {
			addrs = append(addrs, addr)
		}
	}

	if len(addrs) == 0 {
		return nil, errors.New("Kafka storage address configuration is missing")
	}

	producer, err := sarama.NewAsyncProducer(addrs, cfg.ProducerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "New kafkaStorage")
	}

	return newInstance(producer), nil
}

func newInstance(producer sarama.AsyncProducer) *kafkaStorage {
	s := &kafkaStorage{
		quit:     make(chan bool),
		producer: producer,
	}

	go s.daemon()

	return s
}

func (s *kafkaStorage) daemon() {
MainLoop:
	for {
		select {
		case err := <-s.producer.Errors():
			metrics.CounterAdd("kafka_error", 1)
			logger.Errorf("Producer", "err", err)
		case <-s.quit:
			s.quit <- true
			break MainLoop
		}
	}
}

func (s *kafkaStorage) Write(msg *storage.Message) error {
	bs, err := msg.Value.Marshal()
	if err != nil {
		return err
	}
	pmsg := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Key:   sarama.StringEncoder(msg.Key),
		Value: sarama.ByteEncoder(bs),
	}

	s.producer.Input() <- pmsg

	return nil
}

func (s *kafkaStorage) Close() (err error) {

	s.quit <- true
	<-s.quit

	if s.producer != nil {
		err = s.producer.Close()
	}

	return
}
