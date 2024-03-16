package decorator

import (
	"github.com/sirupsen/logrus"
	"mevhub/internal/app/query"
)

type Query interface {
	QueryName() string
}

type QueryHandler[Q Query, R any] interface {
	Handle(ctx *query.Context, qry Q) (R, error)
}

type QueryHandlerWithLogger[Q Query, R any] struct {
	logger *logrus.Logger
	base   QueryHandler[Q, R]
}

func NewQueryHandlerWithLogger[Q Query, R any](logger *logrus.Logger, base QueryHandler[Q, R]) QueryHandler[Q, R] {
	return &QueryHandlerWithLogger[Q, R]{
		logger: logger,
		base:   base,
	}
}

func (h *QueryHandlerWithLogger[Q, R]) Handle(ctx *query.Context, qry Q) (result R, err error) {
	var entry = h.logger.WithFields(logrus.Fields{
		"client.id":  ctx.ClientID,
		"query.name": qry.QueryName(),
	})

	entry.Info("Executing Query")

	defer func() {
		if err == nil {
			entry.Info("Query Executed")
		} else {
			entry.WithError(err).Error("Query Failed")
		}
	}()

	return h.base.Handle(ctx, qry)

}
