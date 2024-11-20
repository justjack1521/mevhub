package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/memory/dto"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
	"strconv"
	"strings"
)

const lobbyParticipantKey = "lobby_participant"
const lobbyKeySeparator = ":"

type LobbyParticipantRepository struct {
	client *redis.Client
}

func NewLobbyParticipantRepository(client *redis.Client) *LobbyParticipantRepository {
	return &LobbyParticipantRepository{client: client}
}

func (r *LobbyParticipantRepository) QueryCountForLobby(ctx context.Context, id uuid.UUID) (int, error) {
	var cursor uint64
	var total int
	var key = r.GenerateLobbyKey(id)
	for {
		keys, next, err := r.client.Scan(ctx, cursor, key, 10).Result()
		if err != nil {
			return 0, err
		}
		cursor = next
		total += len(keys)
		if cursor == 0 {
			break
		}
	}

	return total, nil
}

func (r *LobbyParticipantRepository) QueryParticipantForLobby(ctx context.Context, id uuid.UUID, slot int) (*lobby.Participant, error) {
	var cmd = r.client.HGetAll(ctx, r.GenerateParticipantKey(id, slot))
	if cmd.Err() != nil {
		return nil, port.ErrFailedQueryParticipantForLobby(id, port.ErrFailedQueryParticipant(slot, cmd.Err()))
	}
	var result = &dto.LobbyParticipantRedis{}
	if err := cmd.Scan(result); err != nil {
		return nil, port.ErrFailedQueryParticipantForLobby(id, port.ErrFailedQueryParticipant(slot, err))
	}
	return result.ToEntity(), nil
}

func (r *LobbyParticipantRepository) QueryParticipantExists(ctx context.Context, id uuid.UUID, slot int) (bool, error) {
	var key = r.GenerateParticipantKey(id, slot)
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, port.ErrFailedQueryParticipantForLobby(id, port.ErrFailedQueryParticipantExists(slot, err))
	}
	return result > 0, nil
}

func (r *LobbyParticipantRepository) Delete(ctx context.Context, participant *lobby.Participant) error {
	var key = r.GenerateParticipantKey(participant.LobbyID, participant.PlayerSlot)
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *LobbyParticipantRepository) DeleteAllForLobby(ctx context.Context, id uuid.UUID) error {

	keys, err := r.client.Keys(ctx, r.GenerateLobbyKey(id)).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return err
	}

	return nil

}

func (r *LobbyParticipantRepository) QueryAllForLobby(ctx context.Context, id uuid.UUID) ([]*lobby.Participant, error) {
	keys, err := r.client.Keys(ctx, r.GenerateLobbyKey(id)).Result()
	if err != nil {
		return nil, port.ErrFailedQueryAllParticipantsForLobby(id, err)
	}
	var participants = make([]*lobby.Participant, len(keys))
	for index, key := range keys {
		var cmd = r.client.HGetAll(ctx, key)
		if cmd.Err() != nil {
			return nil, port.ErrFailedQueryAllParticipantsForLobby(id, port.ErrFailedQueryParticipant(index, err))
		}
		var result = &dto.LobbyParticipantRedis{}
		if err := cmd.Scan(result); err != nil {
			return nil, port.ErrFailedQueryAllParticipantsForLobby(id, port.ErrFailedQueryParticipant(index, err))
		}
		participants[index] = result.ToEntity()
	}
	return participants, nil
}

func (r *LobbyParticipantRepository) Create(ctx context.Context, participant *lobby.Participant) error {

	result, err := r.ParticipantToTransfer(participant)
	if err != nil {
		return err
	}

	var key = r.GenerateParticipantKey(participant.LobbyID, participant.PlayerSlot)

	if err := r.client.HSet(ctx, key, result.ToMapStringInterface()).Err(); err != nil {
		return port.ErrFailedCreateParticipantForLobby(participant.LobbyID, port.ErrFailedCreateParticipant(participant.PlayerSlot, err))
	}

	r.client.Expire(ctx, key, lobby.KeepAliveTime)

	return nil
}

func (r *LobbyParticipantRepository) Update(ctx context.Context, participant *lobby.Participant) error {

	result, err := r.ParticipantToTransfer(participant)
	if err != nil {
		return err
	}

	var key = r.GenerateParticipantKey(participant.LobbyID, participant.PlayerSlot)

	if err := r.client.HSet(ctx, key, result.ToMapStringInterface()).Err(); err != nil {
		return port.ErrFailedUpdateParticipantForLobby(participant.LobbyID, port.ErrFailedUpdateParticipant(participant.PlayerSlot, err))
	}

	return nil

}

func (r *LobbyParticipantRepository) ParticipantToTransfer(participant *lobby.Participant) (dto.LobbyParticipantRedis, error) {
	if participant == nil {
		return dto.LobbyParticipantRedis{}, port.ErrParticipantNil
	}
	var result = dto.LobbyParticipantRedis{
		UserID:          participant.UserID.String(),
		PlayerID:        participant.PlayerID.String(),
		LobbyID:         participant.LobbyID.String(),
		RoleRestriction: participant.RoleRestriction.String(),
		Locked:          participant.Locked,
		InviteOnly:      participant.InviteOnly,
		Role:            participant.Role.String(),
		PlayerSlot:      participant.PlayerSlot,
		DeckIndex:       participant.DeckIndex,
		UseStamina:      participant.UseStamina,
		FromInvite:      participant.FromInvite,
		Ready:           participant.Ready,
	}
	return result, nil
}

func (r *LobbyParticipantRepository) GenerateLobbyKey(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, lobbyParticipantKey, id.String(), "*"}, lobbyKeySeparator)
}

func (r *LobbyParticipantRepository) GenerateParticipantKey(id uuid.UUID, slot int) string {
	return strings.Join([]string{serviceKey, lobbyParticipantKey, id.String(), strconv.Itoa(slot)}, lobbyKeySeparator)
}
