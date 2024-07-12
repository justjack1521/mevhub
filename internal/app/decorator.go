package app

import (
	"mevhub/internal/app/query"
	"mevhub/internal/decorator"
	"mevhub/internal/domain/session"
)

type QueryHandlerWithSession[CTX query.Context, Q decorator.Command, R any] struct {
	repository session.InstanceReadRepository
	base       decorator.QueryHandler[CTX, Q, R]
}

func NewQueryHandlerWithSession[CTX query.Context, Q decorator.Command, R any](repository session.InstanceReadRepository, base decorator.QueryHandler[CTX, Q, R]) decorator.QueryHandler[CTX, Q, R] {
	return &QueryHandlerWithSession[CTX, Q, R]{
		repository: repository,
		base:       base,
	}
}

func (h *QueryHandlerWithSession[CTX, Q, R]) Handle(ctx CTX, qry Q) (R, error) {
	instance, err := h.repository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return *new(R), err
	}
	ctx.SetSession(instance)
	return h.base.Handle(ctx, qry)
}
