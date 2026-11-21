package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

// AwayAfterSilentRounds is how many whole rounds in a row a connected player
// must send nothing at all before they are marked away. Missing one wait alone
// is normal: a slow device can overrun the ready wait on a long round, and a
// player can spend the whole turn deciding. A round with no queue edit, lock,
// ready or stamp is not, and two in a row is a client that has stopped sending
// or a player who has walked off.
const AwayAfterSilentRounds = 2

// noteActivity records that a player's client sent something. Every action a
// client sends calls it first, accepted or not: a refused request still shows
// the client is alive. An away player returns at once, so from this message on
// the turn waits count them again.
func noteActivity(instance *game.LiveGameInstance, playerID uuid.UUID) {
	party, err := instance.GetPartyForPlayer(playerID)
	if err != nil {
		return
	}
	player, err := party.GetPlayer(playerID)
	if err != nil {
		return
	}
	player.HeardThisRound = true
	player.SilentRounds = 0
	if player.Away == false {
		return
	}
	player.Away = false
	instance.SendChange(NewPlayerReturnChange(instance.InstanceID, player.PlayerID, party.PartyIndex, player.PartySlot))
}

// markSilentPlayersAway closes out a round at the enemy -> player boundary,
// counting it against everyone who sent nothing since its player turn began.
// The dead are counted too: they still play each round out and ready, and the
// wipe wait would otherwise stall on a silent corpse. The disconnected are not;
// their silence is already explained, and the disconnect grace deals with it.
func markSilentPlayersAway(instance *game.LiveGameInstance) {
	for _, party := range instance.Parties {
		for _, player := range party.Players {
			if player.Away || player.Disconnected || player.HeardThisRound {
				continue
			}
			player.SilentRounds++
			if player.SilentRounds < AwayAfterSilentRounds {
				continue
			}
			player.Away = true
			instance.SendChange(NewPlayerAwayChange(instance.InstanceID, player.PlayerID, party.PartyIndex, player.PartySlot))
		}
	}
}

type PlayerAwayChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerAwayChange(instanceID, playerID uuid.UUID, partyIndex, partySlot int) *PlayerAwayChange {
	return &PlayerAwayChange{InstanceID: instanceID, PlayerID: playerID, PartyIndex: partyIndex, PartySlot: partySlot}
}

func (c PlayerAwayChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerAway
}

type PlayerReturnChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerReturnChange(instanceID, playerID uuid.UUID, partyIndex, partySlot int) *PlayerReturnChange {
	return &PlayerReturnChange{InstanceID: instanceID, PlayerID: playerID, PartyIndex: partyIndex, PartySlot: partySlot}
}

func (c PlayerReturnChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerReturn
}
