package game

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	StateTickPeriod = time.Millisecond * 250
)

var (
	ErrFailedPerformAction = func(id uuid.UUID, err error) error {
		return fmt.Errorf("live game %s failed to perform action: %w", id, err)
	}
)

type PartyInstanceOptions struct {
	MaxPlayerCount     int
	PlayerTurnDuration time.Duration
}

type LiveGameInstance struct {
	InstanceID    uuid.UUID
	ActionChannel chan Action
	ChangeChannel chan Change
	ErrorChannel  chan error
	Parties       map[uuid.UUID]*LiveParty
	State         State
	GameDuration  time.Duration
	Ended         bool
	EndedAt       time.Time
	MaxPartyCount int
	PartyOptions  PartyInstanceOptions
}

func NewLiveGameInstance(source *Instance) *LiveGameInstance {
	var game = &LiveGameInstance{
		InstanceID:    source.SysID,
		ActionChannel: make(chan Action),
		ChangeChannel: make(chan Change),
		ErrorChannel:  make(chan error),
		Parties:       make(map[uuid.UUID]*LiveParty),
		GameDuration:  source.Options.MaxRunTime,
		MaxPartyCount: source.Options.MaxPartyCount,
		PartyOptions: PartyInstanceOptions{
			MaxPlayerCount:     source.Options.MaxPlayerCount,
			PlayerTurnDuration: source.Options.PlayerTurnDuration,
		},
	}
	return game
}

func (game *LiveGameInstance) GetPlayerCount() int {
	var total = 0
	for _, party := range game.Parties {
		total += party.GetPlayerCount()
	}
	return total
}

func (game *LiveGameInstance) PartyExists(id uuid.UUID) bool {
	_, exists := game.Parties[id]
	return exists
}

func (game *LiveGameInstance) GetParty(id uuid.UUID) (*LiveParty, error) {
	party, exists := game.Parties[id]
	if exists == false {
		return nil, ErrPlayerNotInGame
	}
	return party, nil
}

func (game *LiveGameInstance) PlayerExists(id uuid.UUID) bool {
	for _, party := range game.Parties {
		if party.PlayerExists(id) {
			return true
		}
	}
	return false
}

func (game *LiveGameInstance) GetPlayer(id uuid.UUID) (*LivePlayer, error) {
	for _, party := range game.Parties {
		player, err := party.GetPlayer(id)
		if err == nil {
			return player, nil
		}
	}
	return nil, ErrPlayerNotInGame
}

func (game *LiveGameInstance) GetReadyPlayerCount() int {
	var total = 0
	for _, party := range game.Parties {
		total += party.GetReadyPlayerCount()
	}
	return total
}

func (game *LiveGameInstance) RemovePlayer(id uuid.UUID) error {
	for _, party := range game.Parties {
		if party.PlayerExists(id) == false {
			continue
		}
		return party.RemovePlayer(id)
	}
	return ErrPlayerNotInGame
}

func (game *LiveGameInstance) GetActionLockedPlayerCount() int {
	var total = 0
	for _, party := range game.Parties {
		total += party.GetActionLockedPlayerCount()
	}
	return total
}

func (game *LiveGameInstance) Tick() {

	ticker := time.NewTicker(StateTickPeriod)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			game.State.Update(game, t)
		}
	}

}

func (game *LiveGameInstance) WatchActions() {
	for {
		action := <-game.ActionChannel
		if err := action.Perform(game); err != nil {
			game.ErrorChannel <- ErrFailedPerformAction(game.InstanceID, err)
		}
	}
}

func (game *LiveGameInstance) SendChange(change Change) {
	game.ChangeChannel <- change
}
