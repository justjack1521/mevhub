package translate

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/domain/game/action"
)

type GamePlayerRemoveChangeMarshaller Marshaller[*action.PlayerRemoveChange, *protomulti.GamePlayerRemoveNotification]
type GamePlayerDisconnectChangeMarshaller Marshaller[*action.PlayerDisconnectChange, *protomulti.GamePlayerDisconnectNotification]
type GamePlayerReconnectChangeMarshaller Marshaller[*action.PlayerReconnectChange, *protomulti.GamePlayerReconnectNotification]
type GamePlayerReadyChangeMarshaller Marshaller[*action.PlayerReadyChange, *protomulti.GamePlayerReadyNotification]
type GamePlayerEnqueueActionChangeMarshaller Marshaller[*action.PlayerEnqueueActionChange, *protomulti.GameEnqueueActionNotification]
type GamePlayerDequeueActionChangeMarshaller Marshaller[*action.PlayerDequeueActionChange, *protomulti.GameDequeueActionNotification]
type GamePlayerLockActionChangeMarshaller Marshaller[*action.PlayerLockActionChange, *protomulti.GameLockActionNotification]
type GameStateSyncChangeMarshaller Marshaller[*action.GameStateSyncChange, *protomulti.GameSyncNotification]

type gamePlayerRemoveChangeMarshaller struct{}

func NewGamePlayerRemoveChangeMarshaller() GamePlayerRemoveChangeMarshaller {
	return gamePlayerRemoveChangeMarshaller{}
}

func (g gamePlayerRemoveChangeMarshaller) Marshall(data *action.PlayerRemoveChange) (*protomulti.GamePlayerRemoveNotification, error) {
	return &protomulti.GamePlayerRemoveNotification{
		GameId:      data.InstanceID.String(),
		PartyIndex:  int32(data.PartyIndex),
		PlayerIndex: int32(data.PartySlot),
	}, nil
}

type gamePlayerDisconnectChangeMarshaller struct{}

func NewGamePlayerDisconnectChangeMarshaller() GamePlayerDisconnectChangeMarshaller {
	return gamePlayerDisconnectChangeMarshaller{}
}

func (g gamePlayerDisconnectChangeMarshaller) Marshall(data *action.PlayerDisconnectChange) (*protomulti.GamePlayerDisconnectNotification, error) {
	return &protomulti.GamePlayerDisconnectNotification{
		GameId:      data.InstanceID.String(),
		PartyIndex:  int32(data.PartyIndex),
		PlayerIndex: int32(data.PartySlot),
	}, nil
}

type gamePlayerReconnectChangeMarshaller struct{}

func NewGamePlayerReconnectChangeMarshaller() GamePlayerReconnectChangeMarshaller {
	return gamePlayerReconnectChangeMarshaller{}
}

func (g gamePlayerReconnectChangeMarshaller) Marshall(data *action.PlayerReconnectChange) (*protomulti.GamePlayerReconnectNotification, error) {
	return &protomulti.GamePlayerReconnectNotification{
		GameId:      data.InstanceID.String(),
		PartyIndex:  int32(data.PartyIndex),
		PlayerIndex: int32(data.PartySlot),
	}, nil
}

type gameStateSyncChangeMarshaller struct{}

func NewGameStateSyncChangeMarshaller() GameStateSyncChangeMarshaller {
	return gameStateSyncChangeMarshaller{}
}

func (g gameStateSyncChangeMarshaller) Marshall(data *action.GameStateSyncChange) (*protomulti.GameSyncNotification, error) {

	var parties = make([]*protomulti.ProtoGameSyncParty, len(data.Parties))
	for i, party := range data.Parties {
		var players = make([]*protomulti.ProtoGameSyncPlayer, len(party.Players))
		for j, player := range party.Players {
			var actions = make([]*protomulti.ProtoGameAction, len(player.Actions))
			for k, a := range player.Actions {
				actions[k] = &protomulti.ProtoGameAction{
					Action:    protomulti.GamePlayerActionType(a.ActionType),
					Target:    int32(a.Target),
					SlotIndex: int32(a.SlotIndex),
					ElementId: a.ElementID.String(),
				}
			}
			players[j] = &protomulti.ProtoGameSyncPlayer{
				PlayerId:     player.PlayerID.String(),
				PlayerIndex:  int32(player.PlayerIndex),
				Ready:        player.Ready,
				Locked:       player.Locked,
				LockIndex:    int32(player.LockIndex),
				Disconnected: player.Disconnected,
				Actions:      actions,
			}
		}
		parties[i] = &protomulti.ProtoGameSyncParty{
			PartyIndex: int32(party.PartyIndex),
			PartyId:    party.PartyID.String(),
			Players:    players,
		}
	}

	var enemies = make([]*protomulti.ProtoGameEnemyHP, len(data.Enemies))
	for i, e := range data.Enemies {
		enemies[i] = &protomulti.ProtoGameEnemyHP{
			EnemyIndex: int32(e.EnemyIndex),
			Hp:         int32(e.HP),
		}
	}

	return &protomulti.GameSyncNotification{
		GameId:          data.InstanceID.String(),
		Parties:         parties,
		Phase:           protomulti.GameSyncPhase(data.Phase),
		TurnRemainingMs: data.TurnRemainingMs,
		Enemies:         enemies,
	}, nil
}

type gamePlayerReadyChangeMarshaller struct{}

func NewGamePlayerReadyChangeMarshaller() GamePlayerReadyChangeMarshaller {
	return gamePlayerReadyChangeMarshaller{}
}

func (g gamePlayerReadyChangeMarshaller) Marshall(data *action.PlayerReadyChange) (*protomulti.GamePlayerReadyNotification, error) {
	return &protomulti.GamePlayerReadyNotification{
		GameId:      data.InstanceID.String(),
		PartyIndex:  int32(data.PartyIndex),
		PlayerIndex: int32(data.PartySlot),
	}, nil
}

func NewGamePlayerEnqueueActionChangeMarshaller() GamePlayerEnqueueActionChangeMarshaller {
	return gamePlayerEnqueueActionChangeMarshaller{}
}

type gamePlayerEnqueueActionChangeMarshaller struct{}

func (g gamePlayerEnqueueActionChangeMarshaller) Marshall(data *action.PlayerEnqueueActionChange) (*protomulti.GameEnqueueActionNotification, error) {
	return &protomulti.GameEnqueueActionNotification{
		GameId:      data.InstanceID.String(),
		PartyIndex:  int32(data.PartyIndex),
		PlayerIndex: int32(data.PartySlot),
		Action:      protomulti.GamePlayerActionType(data.ActionType),
		SlotIndex:   int32(data.SlotIndex),
		Target:      int32(data.Target),
		ElementId:   data.ElementID.String(),
	}, nil
}

type gamePlayerDequeueActionChangeMarshaller struct{}

func NewGamePlayerDequeueActionChangeMarshaller() GamePlayerDequeueActionChangeMarshaller {
	return gamePlayerDequeueActionChangeMarshaller{}
}

func (g gamePlayerDequeueActionChangeMarshaller) Marshall(data *action.PlayerDequeueActionChange) (*protomulti.GameDequeueActionNotification, error) {
	return &protomulti.GameDequeueActionNotification{
		GameId:      data.InstanceID.String(),
		PartyIndex:  int32(data.PartyIndex),
		PlayerIndex: int32(data.PartySlot),
	}, nil
}

type gamePlayerLockActionChangeMarshaller struct{}

func NewGamePlayerLockActionChangeMarshaller() GamePlayerLockActionChangeMarshaller {
	return gamePlayerLockActionChangeMarshaller{}
}

func (g gamePlayerLockActionChangeMarshaller) Marshall(data *action.PlayerLockActionChange) (*protomulti.GameLockActionNotification, error) {
	return &protomulti.GameLockActionNotification{
		GameId:          data.InstanceID.String(),
		PartyIndex:      int32(data.PartyIndex),
		PlayerIndex:     int32(data.PartySlot),
		ActionLockIndex: int32(data.ActionLockIndex),
	}, nil
}
