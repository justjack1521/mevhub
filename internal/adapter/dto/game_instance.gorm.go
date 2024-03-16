package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/game"
	"time"
)

type GameInstanceGorm struct {
	SysID        uuid.UUID `gorm:"primaryKey;column:sys_id"`
	Seed         int64     `gorm:"column:seed"`
	State        int       `gorm:"column:state"`
	RegisteredAt time.Time `gorm:"column:registered_at"`
}

func (GameInstanceGorm) TableName() string {
	return "multi.game_instance"
}

func (x *GameInstanceGorm) ToEntity() *game.Instance {
	return &game.Instance{
		SysID:        x.SysID,
		Seed:         x.Seed,
		State:        game.InstanceState(x.State),
		RegisteredAt: x.RegisteredAt,
	}
}
