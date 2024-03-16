package mem

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
)

type LobbyInstanceMemoryRepository struct {
	data    map[uuid.UUID]*lobby.Instance
	parties map[string]*lobby.Instance
}

func (r *LobbyInstanceMemoryRepository) QueryByPartyID(ctx context.Context, party string) (*lobby.Instance, error) {
	instance, exists := r.parties[party]
	if exists == false {
		return nil, lobby.ErrFailedQueryLobbyInstance(lobby.ErrLobbyInstanceNotFound)
	}
	if instance == nil {
		return nil, lobby.ErrFailedQueryLobbyInstance(lobby.ErrLobbyInstanceNil)
	}
	return instance, nil
}

func NewLobbyInstanceMemoryRepository() *LobbyInstanceMemoryRepository {
	return &LobbyInstanceMemoryRepository{data: make(map[uuid.UUID]*lobby.Instance), parties: make(map[string]*lobby.Instance)}
}

func (r *LobbyInstanceMemoryRepository) QueryByID(ctx context.Context, id uuid.UUID) (*lobby.Instance, error) {
	instance, exists := r.data[id]
	if exists == false {
		return nil, lobby.ErrFailedQueryLobbyInstance(lobby.ErrLobbyInstanceNotFound)
	}
	if instance == nil {
		return nil, lobby.ErrFailedQueryLobbyInstance(lobby.ErrLobbyInstanceNil)
	}
	return instance, nil
}

func (r *LobbyInstanceMemoryRepository) Create(ctx context.Context, instance *lobby.Instance) error {
	r.data[instance.SysID] = instance
	r.parties[instance.PartyID] = instance
	return nil
}
