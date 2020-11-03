package consumer

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"

	_ "github.com/pkg/errors"

	pb "github.com/techxmind/logserver/interface-defs"
	"github.com/techxmind/logserver/logger"
)

// Sink defines where the EventLog data goes
type Sink interface {
	// Sink errors
	Errors() <-chan error

	// Channel to send *SinkMessage
	Input() chan<- *SinkMessage

	// Ack notify the last successfully processed message
	Ack() <-chan SinkAck

	// Must eventually be called to ensure
	// that any buffered data is written to the underlying io.Writer
	Close() error
}

type SinkMessage struct {
	Ack   SinkAck
	Topic string
	Event *pb.EventLog
}

type SinkAck interface {
	Chain(SinkAck) SinkAck
}

// Marshaler marshal message
//
type Marshaler interface {
	Marshal(*SinkMessage) ([]byte, error)
}

type MarshalerFunc func(*SinkMessage) ([]byte, error)

func (f MarshalerFunc) Marshal(msg *SinkMessage) ([]byte, error) {
	return f(msg)
}

// JSONMarshaler marshal message to json data
//
func JSONMarshaler(msg *SinkMessage) ([]byte, error) {
	j, err := json.Marshal(msg.Event)
	if err != nil {
		return nil, err
	}
	j = append(j, '\n')
	return j, err
}

// CSVMarshaler marshal message to csv text
//
type CSVMarshaler struct {
	// CSV headers, e.g. []string{"EventID", "EventTime"}
	// Header name can be both camelCase or dash_spearated, e.g.  "eventID", "EventID", "event_id"
	// Special names: EventTimeStr, LoggedTimeStr, that are converted by EventTime and LoggedTime
	headers []*standardName

	w  *csv.Writer
	bf *bytes.Buffer
}

func NewCSVMarshaler(originalHeaders []string) (*CSVMarshaler, error) {
	headers := make([]*standardName, 0, len(originalHeaders))
	for _, originalHeader := range originalHeaders {
		if standardHeader, err := newStandardName(originalHeader); err != nil {
			return nil, err
		} else {
			headers = append(headers, standardHeader)
		}
	}

	bf := bytes.NewBuffer(make([]byte, 0, 2000))
	return &CSVMarshaler{
		headers: headers,
		w:       csv.NewWriter(bf),
		bf:      bf,
	}, nil
}

func (d *CSVMarshaler) Marshal(msg *SinkMessage) ([]byte, error) {
	fields := make([]string, 0, len(d.headers))
	for _, name := range d.headers {
		fields = append(fields, name.MustGetValue(msg.Event))
	}

	d.bf.Reset()
	d.w.Write(fields)
	d.w.Flush()

	p := d.bf.Bytes()

	return p, nil
}

type options struct {
	inputBufferSize  int
	outputBufferSize int
}

type Option interface {
	apply(*options)
}

type inputBufferSizeOption int

func (o inputBufferSizeOption) apply(opts *options) {
	opts.inputBufferSize = int(o)
}

type outputBufferSizeOption int

func (o outputBufferSizeOption) apply(opts *options) {
	opts.outputBufferSize = int(o)
}

func WithInputBufferSize(size int) Option {
	if size <= 0 {
		size = 1
	}
	return inputBufferSizeOption(size)
}

func WithOutputBufferSize(size int) Option {
	if size <= 0 {
		size = 1
	}
	return outputBufferSizeOption(size)
}

// NewSink returns new Sink
func NewSink(w io.Writer, m Marshaler, optList ...Option) Sink {
	opts := &options{
		inputBufferSize:  100,
		outputBufferSize: 4096,
	}
	for _, opt := range optList {
		opt.apply(opts)
	}

	s := &defaultSink{
		w:      bufio.NewWriterSize(w, opts.outputBufferSize),
		m:      m,
		input:  make(chan *SinkMessage, opts.inputBufferSize),
		errors: make(chan error, 1),
		ack:    make(chan SinkAck, 1),
		quit:   make(chan struct{}),
	}

	go s.run()

	return s
}

type defaultSink struct {
	w       *bufio.Writer
	m       Marshaler
	input   chan *SinkMessage
	errors  chan error
	ack     chan SinkAck
	lastAck SinkAck
	quit    chan struct{}
	opts    *options
}

func (s *defaultSink) run() {
	for {
		select {
		case msg := <-s.input:
			p, err := s.m.Marshal(msg)
			if err != nil {
				s.errors <- err
				return
			}
			l := len(p)
			if l > s.w.Available() && s.w.Buffered() != 0 {
				s.w.Flush()
				s.ack <- s.lastAck
				s.lastAck = nil
			}

			if _, err := s.w.Write(p); err != nil {
				s.errors <- err
				return
			}
			if s.lastAck != nil {
				s.lastAck = s.lastAck.Chain(msg.Ack)
			} else {
				s.lastAck = msg.Ack
			}
			bufed := s.w.Buffered()
			// if buffer has been flushed
			if l > bufed {
				// then write all buffered data to keep record integrity
				if bufed > 0 {
					s.w.Flush()
				}
				s.ack <- s.lastAck
				s.lastAck = nil
			}
		case <-s.quit:
			if s.lastAck != nil {
				s.ack <- s.lastAck
				s.lastAck = nil
			}
			s.w.Flush()
			s.quit <- struct{}{}
			return
		}
	}
}

func (s *defaultSink) Ack() <-chan SinkAck {
	return s.ack
}

func (s *defaultSink) Errors() <-chan error {
	return s.errors
}

func (s *defaultSink) Input() chan<- *SinkMessage {

	return s.input
}

func (s *defaultSink) Close() error {
	s.quit <- struct{}{}
	<-s.quit

	logger.Debug("Sink close")

	return nil
}
