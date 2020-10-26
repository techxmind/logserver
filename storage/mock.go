package storage

type MockStorage struct {
	returnIndex int
	successes   chan *Message
}

func NewMockStorage(bufferSize int) (*MockStorage, error) {
	return &MockStorage{
		successes: make(chan *Message, bufferSize),
	}, nil
}

func (s *MockStorage) Successes() <-chan *Message {
	return s.successes
}

func (s *MockStorage) Write(msg *Message) error {
	s.successes <- msg

	return nil
}

func (s *MockStorage) Close() error {
	return nil
}
