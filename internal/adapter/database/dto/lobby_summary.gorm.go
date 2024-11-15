package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"time"
)

type LobbySummaryGorm struct {
	LobbyID            uuid.UUID `gorm:"primaryKey;column:lobby_id"`
	QuestID            uuid.UUID `gorm:"column:quest_id"`
	PartyID            string    `gorm:"column:party_id"`
	LobbyComment       string    `gorm:"column:lobby_comment"`
	MinimumPlayerLevel int       `gorm:"column:minimum_player_level"`
	RegisteredAt       time.Time `gorm:"column:registered_at"`
}

func (LobbySummaryGorm) TableName() string {
	return "lobby.lobby_summary"
}

func (x *LobbySummaryGorm) ToEntity() lobby.Summary {
	var summary = lobby.Summary{
		InstanceID:         x.LobbyID,
		QuestID:            x.QuestID,
		PartyID:            x.PartyID,
		LobbyComment:       x.LobbyComment,
		MinimumPlayerLevel: x.MinimumPlayerLevel,
		RegisteredAt:       x.RegisteredAt,
	}
	return summary
}

type LobbyPlayerSummaryGorm struct {
	PlayerID        uuid.UUID                     `gorm:"primaryKey;column:player_id"`
	LobbyID         uuid.UUID                     `gorm:"column:lobby_id"`
	PlayerName      string                        `gorm:"column:player_name"`
	PlayerLevel     int                           `gorm:"column:player_level"`
	PlayerComment   string                        `gorm:"column:player_comment"`
	DeckIndex       int                           `gorm:"column:deck_index"`
	JobCardID       uuid.UUID                     `gorm:"column:job_id"`
	SubJobIndex     int                           `gorm:"column:sub_job_index"`
	CrownLevel      int                           `gorm:"column:crown_level"`
	OverBoostLevel  int                           `gorm:"column:overboost"`
	WeaponID        uuid.UUID                     `gorm:"column:weapon_id"`
	SubWeaponUnlock int                           `gorm:"column:sub_weapon_unlock"`
	AbilityCards    []LobbyAbilityCardSummaryGorm `gorm:"foreignKey:PlayerID"`
}

func (LobbyPlayerSummaryGorm) TableName() string {
	return "lobby.player_summary"
}

func (x *LobbyPlayerSummaryGorm) ToEntity() lobby.PlayerSummary {
	var player = lobby.PlayerSummary{
		Identity: lobby.PlayerIdentity{
			PlayerID:      x.PlayerID,
			PlayerName:    x.PlayerName,
			PlayerComment: x.PlayerComment,
			PlayerLevel:   x.PlayerLevel,
		},
		Loadout: lobby.PlayerLoadout{
			DeckIndex: x.DeckIndex,
			JobCard: lobby.PlayerJobCardSummary{
				JobCardID:      x.JobCardID,
				SubJobIndex:    x.SubJobIndex,
				CrownLevel:     x.CrownLevel,
				OverBoostLevel: x.OverBoostLevel,
			},
			Weapon: lobby.PlayerWeaponSummary{
				WeaponID:        x.WeaponID,
				SubWeaponUnlock: x.SubWeaponUnlock,
			},
			AbilityCards: make([]lobby.PlayerAbilityCardSummary, len(x.AbilityCards)),
		},
	}
	for index, value := range x.AbilityCards {
		player.Loadout.AbilityCards[index] = value.ToEntity()
	}
	return player
}

type LobbyAbilityCardSummaryGorm struct {
	PlayerID         uuid.UUID `gorm:"column:player_id"`
	AbilityCardID    uuid.UUID `gorm:"column:ability_card_id"`
	SlotIndex        int       `gorm:"column:slot_index"`
	AbilityCardLevel int       `gorm:"column:ability_card_level"`
	AbilityLevel     int       `gorm:"column:ability_level"`
	OverBoostLevel   int       `gorm:"column:overboost"`
}

func (LobbyAbilityCardSummaryGorm) TableName() string {
	return "lobby.ability_card_summary"
}

func (x *LobbyAbilityCardSummaryGorm) ToEntity() lobby.PlayerAbilityCardSummary {
	return lobby.PlayerAbilityCardSummary{
		AbilityCardID:    x.AbilityCardID,
		SlotIndex:        x.SlotIndex,
		AbilityCardLevel: x.AbilityCardLevel,
		AbilityLevel:     x.AbilityLevel,
		OverBoostLevel:   x.OverBoostLevel,
	}
}
