package database

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"mevhub/internal/adapter/database/dto"
	"mevhub/internal/core/domain/lobby"
)

type LobbySummaryDatabaseRepository struct {
	database *gorm.DB
}

func NewLobbySummaryDatabaseRepository(db *gorm.DB) *LobbySummaryDatabaseRepository {
	return &LobbySummaryDatabaseRepository{database: db}
}

func (r *LobbySummaryDatabaseRepository) QueryByID(ctx context.Context, id uuid.UUID) (lobby.Summary, error) {

	var cond = &dto.LobbySummaryGorm{LobbyID: id}
	var res = &dto.LobbySummaryGorm{}

	if err := r.database.Model(cond).First(res, cond).Error; err != nil {
		return lobby.Summary{}, err
	}

	return res.ToEntity(), nil

}

func (r *LobbySummaryDatabaseRepository) QueryByPartyID(ctx context.Context, party string) (lobby.Summary, error) {
	var cond = &dto.LobbySummaryGorm{PartyID: party}
	var res = &dto.LobbySummaryGorm{}

	if err := r.database.Model(cond).First(res, cond).Error; err != nil {
		return lobby.Summary{}, err
	}

	return res.ToEntity(), nil

}

func (r *LobbySummaryDatabaseRepository) Create(ctx context.Context, summary lobby.Summary) error {
	if err := r.database.Model(dto.LobbySummaryGorm{}).Create(r.convertLobbySummary(summary)).Error; err != nil {
		return err
	}
	return nil
}

func (r *LobbySummaryDatabaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.database.Model(dto.LobbySummaryGorm{}).Delete(&dto.LobbySummaryGorm{LobbyID: id}).Error; err != nil {
		return err
	}
	return nil
}

func (r *LobbySummaryDatabaseRepository) convertLobbySummary(summary lobby.Summary) *dto.LobbySummaryGorm {
	return &dto.LobbySummaryGorm{
		LobbyID:            summary.InstanceID,
		QuestID:            summary.QuestID,
		PartyID:            summary.PartyID,
		LobbyComment:       summary.LobbyComment,
		MinimumPlayerLevel: summary.MinimumPlayerLevel,
	}
}
