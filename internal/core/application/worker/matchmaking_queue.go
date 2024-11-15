package worker

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
	"time"
)

const (
	matchmakingQueueWorkerInterval = time.Second * 10
)

type MatchmakingQueueWorker struct {
	ctx        context.Context
	mode       game.ModeIdentifier
	repository port.MatchPlayerQueueRepository
	dispatcher port.MatchmakingDispatcher
}

func (w *MatchmakingQueueWorker) Run() {

	var ticker = time.NewTicker(matchmakingQueueWorkerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			actives, err := w.repository.GetActiveQuests(w.ctx, w.mode)
			if err != nil {
				continue
			}
			for _, active := range actives {
				w.process(active)
			}

		}
	}
}

func (w *MatchmakingQueueWorker) process(quest uuid.UUID) error {

	players, err := w.repository.GetQueuedPlayers(w.ctx, w.mode, quest)
	if err != nil {
		return err
	}

	for _, player := range players {
		found, err := w.repository.FindMatch(w.ctx, w.mode, player, 5)
		if err != nil {
			continue
		}
		if err := w.dispatcher.DispatchMatch(w.ctx, w.mode, quest, []match.PlayerQueueEntry{player, found}); err != nil {
			continue
		}
		if err := w.repository.RemovePlayerFromQueue(w.ctx, w.mode, quest, player.PlayerID); err != nil {
			continue
		}
	}

	return nil

}
