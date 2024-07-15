package database

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"mevhub/internal/adapter/dto"
	"mevhub/internal/domain/game"
)

type GameInstanceDatabaseRepository struct {
	db *gorm.DB
}

func NewGameInstanceDatabaseRepository(db *gorm.DB) *GameInstanceDatabaseRepository {
	return &GameInstanceDatabaseRepository{db: db}
}

func (r *GameInstanceDatabaseRepository) Get(ctx context.Context, id uuid.UUID) (*game.Instance, error) {
	var cond = &dto.GameInstanceGorm{SysID: id}
	var res = &dto.GameInstanceGorm{}
	if err := r.db.Model(cond).First(res, cond).Error; err != nil {
		return nil, err
	}
	return res.ToEntity(), nil
}

func (r *GameInstanceDatabaseRepository) Create(ctx context.Context, instance *game.Instance) error {
	var res = &dto.GameInstanceGorm{
		SysID:        instance.SysID,
		Seed:         instance.Seed,
		State:        int(instance.State),
		RegisteredAt: instance.RegisteredAt,
	}
	if err := r.db.Model(&dto.GameInstanceGorm{}).Create(res).Error; err != nil {
		return err
	}
	return nil
}
