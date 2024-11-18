package worker

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
	"time"
)

const (
	matchmakingQueueWorkerInterval = time.Second * 10
)

type LobbyMatchmakingQueueWorker struct {
	ctx        context.Context
	mode       game.ModeIdentifier
	repository port.MatchLobbyPlayerQueueRepository
	dispatcher port.PlayerMatchmakingDispatcher
}

func NewLobbyMatchmakingQueueWorker(ctx context.Context, mode game.ModeIdentifier, queue port.MatchLobbyPlayerQueueRepository, dispatcher port.PlayerMatchmakingDispatcher) *LobbyMatchmakingQueueWorker {
	return &LobbyMatchmakingQueueWorker{ctx: ctx, mode: mode, repository: queue, dispatcher: dispatcher}
}

func (w *LobbyMatchmakingQueueWorker) Run() {

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
				if err := w.process(active); err != nil {
					fmt.Println(err)
				}
			}

		}
	}
}

func (w *LobbyMatchmakingQueueWorker) process(quest uuid.UUID) error {

	lobbies, err := w.repository.GetQueuedLobbies(w.ctx, w.mode, quest)
	if err != nil {
		return err
	}

	for _, queued := range lobbies {
		found, err := w.repository.FindMatch(w.ctx, w.mode, queued, 5)
		if err != nil {
			continue
		}
		if found.Zero() {
			continue
		}
		remove, err := w.dispatcher.Dispatch(w.ctx, w.mode, quest, queued, found)
		if err != nil {
			continue
		}
		if err := w.repository.RemovePlayerFromQueue(w.ctx, w.mode, quest, found.UserID); err != nil {
			continue
		}
		if remove {
			if err := w.repository.RemoveLobbyFromQueue(w.ctx, w.mode, quest, queued.LobbyID); err != nil {
				continue
			}
		}
	}

	return nil

}
