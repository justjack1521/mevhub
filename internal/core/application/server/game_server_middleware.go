package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type ErrorHandler interface {
	Handle(svr *GameServer, err error)
}

type ErrorHandlerDefault struct {
}

func (d ErrorHandlerDefault) Handle(svr *GameServer, err error) {
	svr.errorCount++
}

type ErrorHandlerWithLogging struct {
	logger *logrus.Logger
	next   ErrorHandler
}

func NewErrorLoggingDecorator(logger *logrus.Logger, next ErrorHandler) *ErrorHandlerWithLogging {
	return &ErrorHandlerWithLogging{logger: logger, next: next}
}

func (d ErrorHandlerWithLogging) Handle(svr *GameServer, err error) {
	d.logger.WithFields(logrus.Fields{
		"instance.id": "test",
	}).Error(fmt.Errorf("game server error: %w", err))
	d.next.Handle(svr, err)
}
