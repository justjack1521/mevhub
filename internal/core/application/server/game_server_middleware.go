package server

import (
	"log/slog"
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
	logger *slog.Logger
	next   ErrorHandler
}

func NewErrorLoggingDecorator(logger *slog.Logger, next ErrorHandler) *ErrorHandlerWithLogging {
	return &ErrorHandlerWithLogging{logger: logger, next: next}
}

func (d ErrorHandlerWithLogging) Handle(svr *GameServer, err error) {
	d.logger.With("error", err.Error()).Error("game server error")
	d.next.Handle(svr, err)
}
