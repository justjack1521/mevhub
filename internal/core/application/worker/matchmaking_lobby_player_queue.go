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
	matchmakingLobbyPlayerQueueWorkerFindInterval      = time.Second * 10
	matchmakingLobbyPlayerQueueWorkerLobbyReapInterval = time.Second * 20
	matchmakingLobbyPlayerQueueWorkerQuestReapInterval = time.Minute * 1
)

type LobbyPlayerMatchmakingQueueWorker struct {
	ctx        context.Context
	mode       game.ModeIdentifier
	repository port.MatchLobbyPlayerQueueRepository
	dispatcher port.PlayerMatchmakingDispatcher
}

func NewLobbyPlayerMatchmakingQueueWorker(ctx context.Context, mode game.ModeIdentifier, queues port.MatchLobbyPlayerQueueRepository, dispatcher port.PlayerMatchmakingDispatcher) *LobbyPlayerMatchmakingQueueWorker {
	return &LobbyPlayerMatchmakingQueueWorker{ctx: ctx, mode: mode, repository: queues, dispatcher: dispatcher}
}

func (w *LobbyPlayerMatchmakingQueueWorker) Run() {

	var findTicker = time.NewTicker(matchmakingLobbyPlayerQueueWorkerFindInterval)
	defer findTicker.Stop()

	var lobbyReapTicker = time.NewTicker(matchmakingLobbyPlayerQueueWorkerLobbyReapInterval)
	defer lobbyReapTicker.Stop()

	var questReapTicket = time.NewTicker(matchmakingLobbyPlayerQueueWorkerQuestReapInterval)
	defer questReapTicket.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-findTicker.C:
			actives, err := w.repository.GetActiveQuests(w.ctx, w.mode)
			if err != nil {
				continue
			}
			for _, active := range actives {
				if err := w.processFindMatch(active); err != nil {
					fmt.Println(err)
				}
			}
		case <-lobbyReapTicker.C:
			actives, err := w.repository.GetActiveQuests(w.ctx, w.mode)
			if err != nil {
				continue
			}
			for _, active := range actives {
				if err := w.repository.RemoveExpiredLobbies(w.ctx, w.mode, active); err != nil {
					fmt.Println(err)
				}
			}
		case <-questReapTicket.C:
			actives, err := w.repository.GetActiveQuests(w.ctx, w.mode)
			if err != nil {
				continue
			}
			for _, active := range actives {
				count, err := w.repository.GetCountQueuedLobbies(w.ctx, w.mode, active)
				if err != nil || count > 0 {
					continue
				}
				if err := w.repository.RemoveInactiveQuest(w.ctx, w.mode, active); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func (w *LobbyPlayerMatchmakingQueueWorker) processFindMatch(quest uuid.UUID) error {

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
