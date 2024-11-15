package lobby

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrFailedQueryAllParticipantsForLobby = func(id uuid.UUID, err error) error {
		return fmt.Errorf("failed to query all participants for lobby %s: %w", id, err)
	}
	ErrFailedQueryParticipantForLobby = func(id uuid.UUID, err error) error {
		return fmt.Errorf("failed to query participant for lobby %s: %w", id, err)
	}
	ErrFailedQueryParticipant = func(slot int, err error) error {
		return fmt.Errorf("failed to query participant %d: %w", slot, err)
	}
	ErrFailedQueryParticipantExists = func(slot int, err error) error {
		return fmt.Errorf("failed to query participant %d: %w", slot, err)
	}
)

type ParticipantReadRepository interface {
	QueryParticipantExists(ctx context.Context, id uuid.UUID, slot int) (bool, error)
	QueryAllForLobby(ctx context.Context, id uuid.UUID) ([]*Participant, error)
	QueryParticipantForLobby(ctx context.Context, id uuid.UUID, slot int) (*Participant, error)
}

var (
	ErrFailedCreateParticipantForLobby = func(id uuid.UUID, err error) error {
		return fmt.Errorf("failed to create participant for lobby %s: %w", id, err)
	}
	ErrFailedCreateParticipant = func(slot int, err error) error {
		return fmt.Errorf("failed to create participant %d: %w", slot, err)
	}
	ErrFailedUpdateParticipantForLobby = func(id uuid.UUID, err error) error {
		return fmt.Errorf("failed to update participant for lobby %s: %w", id, err)
	}
	ErrFailedUpdateParticipant = func(slot int, err error) error {
		return fmt.Errorf("failed to update participant %d: %w", slot, err)
	}
	ErrParticipantNil = errors.New("participant is nil")
)

type ParticipantWriteRepository interface {
	Create(ctx context.Context, participant *Participant) error
	Update(ctx context.Context, participant *Participant) error
	Delete(ctx context.Context, participant *Participant) error
	DeleteAllForLobby(ctx context.Context, id uuid.UUID) error
}

type ParticipantRepository interface {
	ParticipantReadRepository
	ParticipantWriteRepository
}
