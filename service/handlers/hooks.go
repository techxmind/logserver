package handlers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/techxmind/logserver/logger"
)

func InterruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)

	// Place whatever shutdown handling you want here
	if _service != nil {
		_service.close()
	}
	logger.Sync()

	errc <- terminateError
}
