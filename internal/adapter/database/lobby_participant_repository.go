package database

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"mevhub/internal/adapter/database/dto"
	"mevhub/internal/core/domain/lobby"
)

type LobbyParticipantDatabaseRepository struct {
	database *gorm.DB
}

func NewGameInstanceParticipantDatabaseRepository(db *gorm.DB) *LobbyParticipantDatabaseRepository {
	return &LobbyParticipantDatabaseRepository{database: db}
}

func (r *LobbyParticipantDatabaseRepository) QueryForClientID(id uuid.UUID) (*lobby.Participant, error) {
	var cond = &dto.LobbyParticipantGorm{UserID: id}
	var res = &dto.LobbyParticipantGorm{}
	if err := r.database.Model(cond).First(res, cond).Error; err != nil {
		return nil, err
	}
	return res.ToEntity(), nil
}

func (r *LobbyParticipantDatabaseRepository) QueryForPlayerID(id uuid.UUID) (*lobby.Participant, error) {
	var cond = &dto.LobbyParticipantGorm{PlayerID: id}
	var res = &dto.LobbyParticipantGorm{}
	if err := r.database.Model(cond).First(res, cond).Error; err != nil {
		return nil, err
	}
	return res.ToEntity(), nil
}

func (r *LobbyParticipantDatabaseRepository) QueryForLobby(id uuid.UUID, slot int) (*lobby.Participant, error) {
	var cond = &dto.LobbyParticipantGorm{GameInstanceID: id, PlayerSlot: slot}
	var res = &dto.LobbyParticipantGorm{}
	if err := r.database.Model(cond).First(res, cond).Error; err != nil {
		return nil, err
	}
	return res.ToEntity(), nil
}

func (r *LobbyParticipantDatabaseRepository) QueryAllForLobby(id uuid.UUID) ([]*lobby.Participant, error) {
	var cond = &dto.LobbyParticipantGorm{GameInstanceID: id}
	var res []*dto.LobbyParticipantGorm
	if err := r.database.Model(cond).Find(&res, cond).Error; err != nil {
		return nil, err
	}
	var dest = make([]*lobby.Participant, len(res))
	for index, value := range res {
		dest[index] = value.ToEntity()
	}
	return dest, nil
}

func (r *LobbyParticipantDatabaseRepository) Create(participant *lobby.Participant) error {
	var res = &dto.LobbyParticipantGorm{
		UserID:         participant.UserID,
		PlayerID:       participant.PlayerID,
		GameInstanceID: participant.LobbyID,
		PlayerSlot:     participant.PlayerSlot,
		DeckIndex:      participant.DeckIndex,
		UseStamina:     participant.UseStamina,
		FromInvite:     participant.FromInvite,
	}
	if err := r.database.Model(&dto.LobbyParticipantGorm{}).Create(res).Error; err != nil {
		return err
	}
	return nil
}
