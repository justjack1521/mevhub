package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

const lobbyStampCooldownKey = "lobby_stamp_cooldown"

type LobbyStampCooldownRedisRepository struct {
	client *redis.Client
}

func NewLobbyStampCooldownRedisRepository(client *redis.Client) *LobbyStampCooldownRedisRepository {
	return &LobbyStampCooldownRedisRepository{client: client}
}

// Claim is a single SET NX with the cooldown as its expiry, so two stamps racing
// each other cannot both get through and the key clears itself.
func (r *LobbyStampCooldownRedisRepository) Claim(ctx context.Context, user uuid.UUID, cooldown time.Duration) (bool, error) {
	return r.client.SetNX(ctx, r.GenerateKey(user), 1, cooldown).Result()
}

func (r *LobbyStampCooldownRedisRepository) GenerateKey(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, lobbyStampCooldownKey, id.String()}, sessionKeySeparator)
}
