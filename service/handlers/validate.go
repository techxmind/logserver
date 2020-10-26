package handlers

import (
	"regexp"

	"github.com/techxmind/logserver/errors"
	pb "github.com/techxmind/logserver/interface-defs"
)

var (
	_rxAlphanumeric = regexp.MustCompile(`^\w*$`)
)

func validate(log *pb.EventLog) (err error) {
	if err = checkRequired(log); err != nil {
		return
	}

	return
}

func checkRequired(log *pb.EventLog) error {
	// TPL.CHECK_REQUIRED.START EventId,EventTime,Event
	if log.EventId == "" {
		return errors.Wrap(errors.ErrFieldRequired, "EventId")
	}

	if log.EventTime == 0 {
		return errors.Wrap(errors.ErrFieldRequired, "EventTime")
	}

	if log.Event == "" {
		return errors.Wrap(errors.ErrFieldRequired, "Event")
	}

	// TPL.CHECK_REQUIRED.END

	if !_rxAlphanumeric.MatchString(log.AppType) {
		return errors.Wrap(errors.ErrFieldInvalid, "'AppType' continas chars other than alpha, numeric")
	}

	if !_rxAlphanumeric.MatchString(log.Event) {
		return errors.Wrap(errors.ErrFieldInvalid, "'AppType' continas chars other than alpha, numeric")
	}

	return nil
}
