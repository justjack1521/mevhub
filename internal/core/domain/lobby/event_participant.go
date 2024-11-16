package lobby

import "github.com/justjack1521/mevium/pkg/mevent"

type ParticipantEvent interface {
	mevent.ContextEvent
	LobbyEvent
	DeckIndex() int
	SlotIndex() int
}
