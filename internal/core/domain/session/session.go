package session

import (
	"errors"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrUserIDNil   = errors.New("session user id is nil")
	ErrPlayerIDNil = errors.New("player id is nil")
)

type Instance struct {
	UserID              uuid.UUID
	PlayerID            uuid.UUID
	LobbyID             uuid.UUID
	GameID              uuid.UUID
	PartySlot           int
	DeckIndex           int
	DisconnectSessionID uuid.UUID
	// LastConnEventAt is the gateway emit timestamp of the most recently applied
	// connect/disconnect event. Connect/disconnect notifications that arrive out
	// of order (older timestamp) are ignored as stale.
	LastConnEventAt int64
	// DisconnectedAt is the unix-millisecond time the in-game disconnect grace
	// period started; 0 means not disconnected. Persisted so grace state
	// survives a process restart instead of living only in server memory.
	DisconnectedAt int64
	// CurrentSessionID is the gateway session that most recently connected
	// for this user. A connect from a different gateway session while this
	// session still references a lobby is an app relaunch, not a resume.
	CurrentSessionID uuid.UUID
}

func NewInstance(user uuid.UUID, player uuid.UUID) (*Instance, error) {
	if user == uuid.Nil {
		return nil, ErrUserIDNil
	}
	if player == uuid.Nil {
		return nil, ErrPlayerIDNil
	}
	return &Instance{UserID: user, PlayerID: player}, nil
}

// ApplyConnEvent reports whether a connect/disconnect event predates the most
// recently applied one for this session (out-of-order delivery, to be
// ignored). When the event is current it advances the watermark. A zero
// timestamp carries no ordering information and is never treated as stale.
func (x *Instance) ApplyConnEvent(timestamp int64) (stale bool) {
	if timestamp != 0 && timestamp <= x.LastConnEventAt {
		return true
	}
	if timestamp > x.LastConnEventAt {
		x.LastConnEventAt = timestamp
	}
	return false
}

func (x *Instance) CanJoinLobby() bool {
	return uuid.Equal(x.LobbyID, uuid.Nil)
}

func (x *Instance) CanLeaveLobby() bool {
	return uuid.Equal(x.LobbyID, uuid.Nil) == false
}
