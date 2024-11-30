package server

import "log/slog"

type ErrorHandlerWithLogging struct {
	logger  *slog.Logger
	handler ErrorHandler
}

func NewErrorLoggingDecorator(logger *slog.Logger, handler ErrorHandler) *ErrorHandlerWithLogging {
	return &ErrorHandlerWithLogging{logger: logger, handler: handler}
}

func (d *ErrorHandlerWithLogging) Handle(svr *GameServer, err error) {
	d.handler.Handle(svr, err)
	d.logger.With("instance.id", svr.InstanceID.String(), "error", err.Error(), "total", svr.errorCount).Error("game server error")
}
