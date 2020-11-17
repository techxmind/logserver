package consumer

import (
	"context"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	_ "go.uber.org/atomic"

	pb "github.com/techxmind/logserver/interface-defs"
	"github.com/techxmind/logserver/logger"
)

type Consumer struct {
	group  sarama.ConsumerGroup
	topics []string
	sink   Sink
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func New(pctx context.Context, cfg *Config) (*Consumer, error) {
	var (
		addrs   = cfg.GetAddrs()
		topics  = cfg.GetTopics()
		sinkCfg = cfg.Sink
	)

	if len(addrs) == 0 {
		return nil, errors.New("Kafka address configuration is missing")
	}
	if len(topics) == 0 {
		return nil, errors.New("Topics configuration is missing")
	}
	if sinkCfg == nil {
		sinkCfg = DefaultConfig.Sink
	}

	marshaler, err := GetMarshaler(sinkCfg.Marshaler, sinkCfg.MarshalerArgs)
	if err != nil {
		return nil, err
	}
	sinkTarget, err := GetSinkTarget(sinkCfg.Target, sinkCfg.TargetArgs)
	if err != nil {
		return nil, err
	}
	sink := NewSink(
		sinkTarget,
		marshaler,
		WithInputBufferSize(sinkCfg.InputBufferSize),
		WithOutputBufferSize(sinkCfg.OutputBufferSize),
	)

	ctx, cancel := context.WithCancel(pctx)
	consumer := &Consumer{
		ctx:    ctx,
		topics: topics,
		cancel: cancel,
		sink:   sink,
	}

	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = cfg.GetKafkaVersion()
	consumerConfig.Consumer.Return.Errors = true

	if cfg.Offset == "oldest" {
		consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	group, err := sarama.NewConsumerGroup(addrs, cfg.GroupID, consumerConfig)
	if err != nil {
		logger.Fatal("NewConsumerGroup", "err", err)
	}

	consumer.group = group

	logger.Debug(
		"NewConsumer",
		"sink.marshaler", sinkCfg.Marshaler,
		"sink.marshalerArgs", sinkCfg.MarshalerArgs,
		"sink.target", sinkCfg.Target,
		"sink.targetArgs", sinkCfg.TargetArgs,
		"sink.inputBufferSize", sinkCfg.InputBufferSize,
		"sink.outputBufferSize", sinkCfg.OutputBufferSize,
		"kafka.addrs", strings.Join(addrs, ","),
		"kafka.version", consumerConfig.Version,
		"topics", strings.Join(topics, ","),
	)

	return consumer, nil
}

func (c *Consumer) Start() {

	go c.worker()

	for {
		err := c.group.Consume(c.ctx, c.topics, c)
		if err != nil {
			logger.Fatal("consume error", "err", err)
		}
	}
}

func (c *Consumer) worker() {
OUTLOOP:
	for {
		select {
		case err := <-c.group.Errors():
			logger.Error("Consumer", "err", err)

		case err := <-c.sink.Errors():
			logger.Error("Sink", "err", err)
			c.Close()
			break OUTLOOP

		case iack := <-c.sink.Ack():
			logger.Debug("sink.Ack")
			if ack, ok := iack.(*Ack); ok {
				ack.MarkOffset()
			}

		case <-c.ctx.Done():
			logger.Debug("Quit worker")
			break OUTLOOP
		}
	}
}

type Ack struct {
	Session   sarama.ConsumerGroupSession
	Topic     string
	Partition int32
	Offset    int64

	chained map[string]map[int32]int64
}

// Chain merge ack coming later
// It's not thread-safe
func (ack *Ack) Chain(dst SinkAck) SinkAck {
	dstAck, ok := dst.(*Ack)
	if !ok {
		return ack
	}
	if ack.chained == nil {
		ack.chained = make(map[string]map[int32]int64)
	}
	if _, ok := ack.chained[dstAck.Topic]; !ok {
		ack.chained[dstAck.Topic] = make(map[int32]int64)
	}

	if dstAck.Offset > ack.chained[dstAck.Topic][dstAck.Partition] {
		ack.chained[dstAck.Topic][dstAck.Partition] = dstAck.Offset
	}

	return ack
}

func (ack *Ack) MarkOffset() {
	logger.Debug("MarkOffset single", "topic", ack.Topic, "partition", ack.Partition, "offset", ack.Offset)
	ack.Session.MarkOffset(ack.Topic, ack.Partition, ack.Offset, "")

	for topic, partitions := range ack.chained {
		for partition, offset := range partitions {
			logger.Debug("MarkOffset chained", "topic", topic, "partition", partition, "offset", offset)
			ack.Session.MarkOffset(topic, partition, offset, "")
		}
	}
}

// ConsumeClaim implements sarama.ConsumerGroupHandler
//
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.wg.Add(1)
	defer c.wg.Done()

OUTLOOP:
	for {
		select {
		case msg := <-claim.Messages():
			if msg == nil {
				break OUTLOOP
			}
			event := &pb.EventLog{}
			err := event.XXX_Unmarshal(msg.Value)
			if err != nil {
				logger.Error("Message unmarshal", "err", err)
				continue
			}
			c.sink.Input() <- &SinkMessage{
				Topic: msg.Topic,
				Event: event,
				Ack: &Ack{
					Session:   sess,
					Topic:     msg.Topic,
					Partition: msg.Partition,
					Offset:    msg.Offset,
				},
			}

		case <-c.ctx.Done():
			logger.Debug("Consumer cancel")
			break OUTLOOP
		}
	}

	return nil
}

// Setup implements sarama.ConsumerGroupHandler
//
func (c *Consumer) Setup(s sarama.ConsumerGroupSession) error {
	logger.Info("Setup", "memberID", s.MemberID())
	return nil
}

// Cleanupimplements sarama.ConsumerGroupHandler
//
func (c *Consumer) Cleanup(s sarama.ConsumerGroupSession) error {
	logger.Info("Cleanup", "memberID", s.MemberID())
	return nil
}

func (c *Consumer) Close() {
	logger.Info("Close")
	c.sink.Close()
	c.cancel()
	c.wg.Wait()
	c.group.Close()
}
