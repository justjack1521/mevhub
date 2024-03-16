package game

import (
	"context"
	"fmt"
)

var (
	ErrInstanceOptionsNotFound = func(err error) error {
		return fmt.Errorf("instance options not found: %w", err)
	}
	ErrInstanceOptionsNotFoundWithIdentifier = func(identifier ModeIdentifier) error {
		return fmt.Errorf("instance options not found with identifier: %s", identifier)
	}
)

type InstanceOptionsRepository interface {
	QueryByIdentifier(ctx context.Context, identifier ModeIdentifier) (InstanceOptions, error)
}
