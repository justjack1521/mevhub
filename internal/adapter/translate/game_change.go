package translate

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/domain/game/action"
)

type GamePlayerRemoveChangeMarshaller Marshaller[action.PlayerRemoveChange, *protomulti.GamePlayerRemoveNotification]
type GamePlayerReadyChangeMarshaller Marshaller[action.PlayerReadyChange, *protomulti.GamePlayerReadyNotification]
type GamePlayerEnqueueActionChangeMarshaller Marshaller[action.PlayerEnqueueActionChange, *protomulti.GameEnqueueActionNotification]
type GamePlayerDequeueActionChangeMarshaller Marshaller[action.PlayerDequeueActionChange, *protomulti.GameDequeueActionNotification]
type GamePlayerLockActionChangeMarshaller Marshaller[action.PlayerLockActionChange, *protomulti.GameLockActionNotification]

type gamePlayerRemoveChangeMarshaller struct{}

func NewGamePlayerRemoveChangeMarshaller() GamePlayerRemoveChangeMarshaller {
	return gamePlayerRemoveChangeMarshaller{}
}

func (g gamePlayerRemoveChangeMarshaller) Marshall(data action.PlayerRemoveChange) (*protomulti.GamePlayerRemoveNotification, error) {
	return &protomulti.GamePlayerRemoveNotification{
		GameId:      data.InstanceID.String(),
		PartyIndex:  int32(data.PartyIndex),
		PlayerIndex: int32(data.PartySlot),
	}, nil
}

type gamePlayerReadyChangeMarshaller struct{}

func NewGamePlayerReadyChangeMarshaller() GamePlayerReadyChangeMarshaller {
	return gamePlayerReadyChangeMarshaller{}
}

func (g gamePlayerReadyChangeMarshaller) Marshall(data action.PlayerReadyChange) (*protomulti.GamePlayerReadyNotification, error) {
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

func (g gamePlayerEnqueueActionChangeMarshaller) Marshall(data action.PlayerEnqueueActionChange) (*protomulti.GameEnqueueActionNotification, error) {
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

func (g gamePlayerDequeueActionChangeMarshaller) Marshall(data action.PlayerDequeueActionChange) (*protomulti.GameDequeueActionNotification, error) {
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

func (g gamePlayerLockActionChangeMarshaller) Marshall(data action.PlayerLockActionChange) (*protomulti.GameLockActionNotification, error) {
	return &protomulti.GameLockActionNotification{
		GameId:          data.InstanceID.String(),
		PartyIndex:      int32(data.PartyIndex),
		PlayerIndex:     int32(data.PartySlot),
		ActionLockIndex: int32(data.ActionLockIndex),
	}, nil
}
