package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	_ = assert.New(t)

	g := NewGroup(New(os.Stdout), New(os.Stdout))

	g.Write(&Message{
		Topic: "test",
		Value: StringMarshaler("message"),
	})
}
