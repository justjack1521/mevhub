package database

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"mevhub/internal/adapter/dto"
	"mevhub/internal/domain/game"
)

type GameQuestDatabaseRepository struct {
	database *gorm.DB
}

func NewGameQuestDatabaseRepository(db *gorm.DB) *GameQuestDatabaseRepository {
	return &GameQuestDatabaseRepository{database: db}
}

func (r *GameQuestDatabaseRepository) QueryByID(id uuid.UUID) (game.Quest, error) {
	var cond = &dto.GameQuestGorm{SysID: id}
	var res = &dto.GameQuestGorm{}

	if err := r.database.Model(cond).Preload("Tier.Mode").Preload("Categories").First(res, cond).Error; err != nil {
		return game.Quest{}, err
	}
	return res.ToEntity(), nil

}
