package port

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
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

type LobbyParticipantReadRepository interface {
	QueryParticipantExists(ctx context.Context, id uuid.UUID, slot int) (bool, error)
	QueryAllForLobby(ctx context.Context, id uuid.UUID) ([]*lobby.Participant, error)
	QueryCountForLobby(ctx context.Context, id uuid.UUID) (int, error)
	QueryParticipantForLobby(ctx context.Context, id uuid.UUID, slot int) (*lobby.Participant, error)
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

type LobbyParticipantWriteRepository interface {
	Create(ctx context.Context, participant *lobby.Participant) error
	Update(ctx context.Context, participant *lobby.Participant) error
	Delete(ctx context.Context, participant *lobby.Participant) error
	DeleteAllForLobby(ctx context.Context, id uuid.UUID) error
}

type LobbyParticipantRepository interface {
	LobbyParticipantReadRepository
	LobbyParticipantWriteRepository
}
