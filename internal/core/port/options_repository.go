package port

import (
	"context"
	"fmt"
	"mevhub/internal/core/domain/game"
)

var (
	ErrInstanceOptionsNotFound = func(err error) error {
		return fmt.Errorf("instance options not found: %w", err)
	}
	ErrInstanceOptionsNotFoundWithIdentifier = func(identifier game.ModeIdentifier) error {
		return fmt.Errorf("instance options not found with identifier: %s", identifier)
	}
)

type InstanceOptionsRepository interface {
	QueryByIdentifier(ctx context.Context, identifier game.ModeIdentifier) (game.InstanceOptions, error)
}
