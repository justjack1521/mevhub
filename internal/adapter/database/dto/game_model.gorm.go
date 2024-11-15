package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

type GameCategoryGorm struct {
	QuestID  uuid.UUID `gorm:"column:quest"`
	Category uuid.UUID `gorm:"column:element"`
}

func (GameCategoryGorm) TableName() string {
	return "multi.game_quest_category"
}

func (x *GameCategoryGorm) ToEntity() game.Category {
	return game.Category{SysID: x.Category}
}

type GameQuestGorm struct {
	SysID      uuid.UUID          `gorm:"primaryKey;column:sys_id"`
	TierID     uuid.UUID          `gorm:"column:tier"`
	Tier       *GameModeTierGorm  `gorm:"foreignKey:TierID"`
	Name       string             `gorm:"column:name"`
	Categories []GameCategoryGorm `gorm:"foreignKey:QuestID"`
}

func (GameQuestGorm) TableName() string {
	return "multi.game_quest"
}

func (x *GameQuestGorm) ToEntity() game.Quest {
	var result = game.Quest{
		SysID:      x.SysID,
		Tier:       x.Tier.ToEntity(),
		Categories: make([]game.Category, len(x.Categories)),
	}
	if x.Categories == nil {
		return result
	}
	for index, value := range x.Categories {
		result.Categories[index] = value.ToEntity()
	}
	return result
}

type GameModeTierGorm struct {
	SysID          uuid.UUID     `gorm:"primaryKey;column:sys_id"`
	GameModeID     uuid.UUID     `gorm:"column:game_mode"`
	Mode           *GameModeGorm `gorm:"foreignKey:GameModeID"`
	StarLevel      int           `gorm:"column:star_level"`
	StaminaCost    int           `gorm:"column:stamina_cost"`
	TimeLimit      int           `gorm:"column:time_limit"`
	SeedMultiplier int           `gorm:"column:seed_multiplier"`
}

func (GameModeTierGorm) TableName() string {
	return "multi.game_mode_tier"
}

func (x *GameModeTierGorm) ToEntity() game.Tier {
	return game.Tier{
		SysID:          x.SysID,
		GameMode:       x.Mode.ToEntity(),
		StarLevel:      x.StarLevel,
		StaminaCost:    x.StaminaCost,
		TimeLimit:      time.Duration(x.TimeLimit) * time.Minute,
		SeedMultiplier: x.SeedMultiplier,
	}
}

type GameModeGorm struct {
	SysID      uuid.UUID `gorm:"primaryKey;column:sys_id"`
	Identifier string    `gorm:"column:identifier"`
	MaxPlayers int       `gorm:"column:max_players"`
}

func (GameModeGorm) TableName() string {
	return "multi.game_mode"
}

func (x *GameModeGorm) ToEntity() game.Mode {
	return game.Mode{
		SysID:             x.SysID,
		OptionsIdentifier: game.ModeIdentifier(x.Identifier),
		MaxPlayers:        x.MaxPlayers,
	}
}
