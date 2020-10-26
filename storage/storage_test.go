package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	_ = assert.New(t)

	s := New(os.Stdout)

	s.Write(&Message{
		Topic: "test",
		Value: StringMarshaler("message"),
	})
}
