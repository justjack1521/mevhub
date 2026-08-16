package game

import (
	"time"
)

const (
	PendingStateMaxWaitDuration = time.Minute * 30
)

type State interface {
	Update(game *LiveGameInstance, t time.Time)
}
