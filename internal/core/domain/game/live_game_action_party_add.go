package game

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type PartyAddAction struct {
	PartyID    uuid.UUID
	PartyIndex int
}

func NewPartyAddAction(partyID uuid.UUID, partyIndex int) *PartyAddAction {
	return &PartyAddAction{PartyID: partyID, PartyIndex: partyIndex}
}

func (a *PartyAddAction) Perform(game *LiveGameInstance) error {

	if game.PartyExists(a.PartyID) {
		return ErrPartyAlreadyInGame
	}

	if uuid.Equal(a.PartyID, uuid.Nil) {
		return ErrPartyIDNil
	}

	game.Parties[a.PartyID] = &LiveParty{
		PartyID:            a.PartyID,
		PartyIndex:         a.PartyIndex,
		Players:            make(map[uuid.UUID]*LivePlayer),
		MaxPlayerCount:     game.PartyOptions.MaxPlayerCount,
		PlayerTurnDuration: game.PartyOptions.PlayerTurnDuration,
		LastAction:         time.Now().UTC(),
	}

	game.SendChange(NewPartyAddChange(a.PartyID, a.PartyIndex))

	return nil

}
