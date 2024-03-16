package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/domain/lobby"
)

type LobbyPlayerSummaryWriter struct {
	EventPublisher          *mevent.Publisher
	PlayerSummaryRepository lobby.PlayerSummaryRepository
}

func NewLobbyPlayerSummaryWriter(publisher *mevent.Publisher, summary lobby.PlayerSummaryRepository) *LobbyPlayerSummaryWriter {
	var subscriber = &LobbyPlayerSummaryWriter{EventPublisher: publisher, PlayerSummaryRepository: summary}
	publisher.Subscribe(subscriber, lobby.ParticipantCreatedEvent{})
	return subscriber
}

func (s *LobbyPlayerSummaryWriter) Notify(event mevent.Event) {
	actual, valid := event.(lobby.ParticipantCreatedEvent)
	if valid == false {
		return
	}
	if err := s.Handle(actual); err != nil {
		fmt.Println(err)
	}
}

func (s *LobbyPlayerSummaryWriter) Handle(event lobby.ParticipantCreatedEvent) error {

	//player, err := s.PlayerSummaryRepository.Query(event.Context(), event.PlayerID(), event.DeckIndex())
	//if err != nil {
	//	return err
	//}
	//
	//var summary = lobby.PlayerSummary{
	//	Identity: lobby.PlayerIdentity{
	//		PlayerID:      player.PlayerID,
	//		PlayerName:    player.PlayerName,
	//		PlayerComment: player.PlayerComment,
	//		PlayerLevel:   player.PlayerLevel,
	//	},
	//	Loadout: lobby.PlayerLoadout{
	//		DeckIndex: player.DeckIndex,
	//		JobCard: lobby.PlayerJobCardSummary{
	//			JobCardID:      player.JobCard.JobCardID,
	//			SubJobIndex:    player.JobCard.SubJobIndex,
	//			OverBoostLevel: player.JobCard.OverBoostLevel,
	//			CrownLevel:     player.JobCard.CrownLevel,
	//		},
	//		Weapon: lobby.PlayerWeaponSummary{
	//			WeaponID:        player.Weapon.WeaponID,
	//			SubWeaponUnlock: player.Weapon.SubWeaponUnlock,
	//		},
	//		AbilityCards: make([]lobby.PlayerAbilityCardSummary, len(player.AbilityCards)),
	//	},
	//}
	//
	//for index, value := range summary.Loadout.AbilityCards {
	//	summary.Loadout.AbilityCards[index] = lobby.PlayerAbilityCardSummary{
	//		AbilityCardID:    value.AbilityCardID,
	//		SlotIndex:        value.SlotIndex,
	//		AbilityCardLevel: value.AbilityCardLevel,
	//		AbilityLevel:     value.AbilityLevel,
	//		OverBoostLevel:   value.OverBoostLevel,
	//	}
	//}
	//
	//if err := s.LobbySummaryRepository.Create(event.Context(), summary); err != nil {
	//	return err
	//}

	return nil

}
