package app

import (
	"mevhub/internal/app/command"
	"mevhub/internal/app/query"
	"mevhub/internal/decorator"
	"mevhub/internal/domain/session"
)

type QueryHandlerWithSession[Q decorator.Query, R any] struct {
	repository session.InstanceReadRepository
	base       decorator.QueryHandler[Q, R]
}

func NewQueryHandlerWithSession[Q decorator.Query, R any](repository session.InstanceReadRepository, base decorator.QueryHandler[Q, R]) decorator.QueryHandler[Q, R] {
	return &QueryHandlerWithSession[Q, R]{
		repository: repository,
		base:       base,
	}
}

func (h *QueryHandlerWithSession[Q, R]) Handle(ctx *query.Context, qry Q) (R, error) {
	instance, err := h.repository.QueryByID(ctx.Context, ctx.ClientID)
	if err != nil {
		return *new(R), err
	}
	ctx.PlayerID = instance.PlayerID
	return h.base.Handle(ctx, qry)
}

type CommandHandlerWithSession[C decorator.Command] struct {
	repository session.InstanceReadRepository
	base       decorator.CommandHandler[C]
}

func NewCommandHandlerWithSession[C decorator.Command](repository session.InstanceReadRepository, base decorator.CommandHandler[C]) decorator.CommandHandler[C] {
	return &CommandHandlerWithSession[C]{
		repository: repository,
		base:       base,
	}
}

func (h *CommandHandlerWithSession[C]) Handle(ctx *command.Context, cmd C) (err error) {
	instance, err := h.repository.QueryByID(ctx.Context, ctx.ClientID)
	if err != nil {
		return err
	}
	ctx.Session = instance
	return h.base.Handle(ctx, cmd)
}
