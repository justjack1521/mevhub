package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

type PartyAddAction struct {
	PartyID    uuid.UUID
	PartyIndex int
}

func NewPartyAddAction(partyID uuid.UUID, partyIndex int) *PartyAddAction {
	return &PartyAddAction{PartyID: partyID, PartyIndex: partyIndex}
}

func (a *PartyAddAction) Perform(instance *game.LiveGameInstance) error {

	if instance.PartyExists(a.PartyID) {
		return game.ErrPartyAlreadyInGame
	}

	if uuid.Equal(a.PartyID, uuid.Nil) {
		return game.ErrPartyIDNil
	}

	instance.Parties[a.PartyID] = &game.LiveParty{
		PartyID:            a.PartyID,
		PartyIndex:         a.PartyIndex,
		Players:            make(map[uuid.UUID]*game.LivePlayer),
		MaxPlayerCount:     instance.PartyOptions.MaxPlayerCount,
		PlayerTurnDuration: instance.PartyOptions.PlayerTurnDuration,
		LastAction:         time.Now().UTC(),
	}

	instance.SendChange(NewPartyAddChange(a.PartyID, a.PartyIndex))

	return nil

}
