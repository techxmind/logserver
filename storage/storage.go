package storage

import (
	"encoding/json"
	"io"
)

type Storager interface {
	Write(*Message) error
	Close() error
}

type Marshaler interface {
	Marshal() ([]byte, error)
}

type StringMarshaler string

func (s StringMarshaler) Marshal() ([]byte, error) {
	return []byte(s), nil
}

type jsonMarshaler struct {
	v interface{}
}

func JSONMarshaler(v interface{}) Marshaler {
	return &jsonMarshaler{
		v: v,
	}
}

func (m *jsonMarshaler) Marshal() ([]byte, error) {
	return json.Marshal(m.v)
}

type Message struct {
	Topic string
	Key   string
	Value Marshaler
}

func NewMessage(topic, key string, value Marshaler) *Message {
	return &Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}
}

type Storage struct {
	w io.Writer
}

func New(w io.Writer) Storager {
	return &Storage{
		w: w,
	}
}

func (s *Storage) Write(msg *Message) error {
	if _, err := s.w.Write([]byte(msg.Topic + ":")); err != nil {
		return err
	}

	if bs, err := msg.Value.Marshal(); err != nil {
		return err
	} else {
		if _, err := s.w.Write(bs); err != nil {
			return err
		}
	}

	if _, err := s.w.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() error {
	return nil
}
