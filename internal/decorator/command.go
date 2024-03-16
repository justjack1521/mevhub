package decorator

import (
	"github.com/sirupsen/logrus"
	"mevhub/internal/app/command"
)

type Command interface {
	CommandName() string
}

type CommandHandler[C Command] interface {
	Handle(ctx *command.Context, cmd C) error
}

type CommandHandlerWithLogger[C Command] struct {
	logger *logrus.Logger
	base   CommandHandler[C]
}

func NewCommandHandlerWithLogger[C Command](logger *logrus.Logger, base CommandHandler[C]) CommandHandler[C] {
	return &CommandHandlerWithLogger[C]{
		logger: logger,
		base:   base,
	}
}

func (h *CommandHandlerWithLogger[C]) Handle(ctx *command.Context, cmd C) (err error) {
	var entry = h.logger.WithFields(logrus.Fields{
		"client.id":    ctx.ClientID,
		"command.name": cmd.CommandName(),
	})

	entry.Info("Executing Command")

	defer func() {
		if err == nil {
			entry.Info("Command Executed")
		} else {
			entry.WithError(err).Error("Command Failed")
		}
	}()

	return h.base.Handle(ctx, cmd)

}
