package game

import (
	"time"
)

const (
	pendingStateMaxWaitDuration = time.Minute * 1
)

type State interface {
	Update(game *LiveGameInstance, t time.Time)
}
