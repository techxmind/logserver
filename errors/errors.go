package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	Wrap  = errors.Wrap
	Wrapf = errors.Wrapf
)

type LError int16

const (
	ErrNoError       LError = 0
	ErrUnknown       LError = -1
	ErrFieldRequired LError = 1
	ErrFieldInvalid  LError = 2
)

func (err LError) Error() string {
	switch err {
	case ErrNoError:
		return "No error"
	case ErrUnknown:
		return "Unexpected server error."
	case ErrFieldRequired:
		return "Field value is required."
	case ErrFieldInvalid:
		return "Field value is not valid."
	}

	return fmt.Sprintf("Unknow error. Error code = %d", err)
}

func (err LError) ErrorCode() int16 {
	return int16(err)
}
