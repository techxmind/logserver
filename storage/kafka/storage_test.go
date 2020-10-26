package kafka

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/techxmind/logserver/storage"
)

func TestAll(t *testing.T) {
	ast := assert.New(t)

	cfg := mocks.NewTestConfig()
	cfg.Producer.Return.Successes = true
	mp := mocks.NewAsyncProducer(t, cfg)

	s := newInstance(mp)

	mp.ExpectInputAndSucceed()
	mp.ExpectInputAndFail(sarama.ErrOutOfBrokers)

	s.Write(&storage.Message{
		Topic: "test",
		Key:   "key",
		Value: storage.StringMarshaler("value"),
	})

	s.Write(&storage.Message{
		Topic: "test",
		Key:   "key",
		Value: storage.StringMarshaler("value"),
	})

	msg1 := <-mp.Successes()

	ast.Equal("test", msg1.Topic)
	key, _ := msg1.Key.Encode()
	ast.Equal("key", string(key))
	value, _ := msg1.Value.Encode()
	ast.Equal("value", string(value))

	s.Close()
}
