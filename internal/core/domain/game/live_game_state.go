package game

import (
	"time"
)

const (
	PendingStateMaxWaitDuration = time.Minute * 1
)

type State interface {
	Update(game *LiveGameInstance, t time.Time)
}
