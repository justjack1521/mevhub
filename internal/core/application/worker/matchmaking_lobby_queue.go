package worker

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
	"time"
)

const (
	matchmakingLobbyQueueWorkerFindInterval      = time.Second * 10
	matchmakingLobbyQueueWorkerLobbyReapInterval = time.Second * 20
)

type LobbyMatchmakingQueueWorker struct {
	ctx        context.Context
	mode       game.ModeIdentifier
	repository port.MatchLobbyQueueRepository
	dispatcher port.LobbyMatchmakingDispatcher
}

func NewLobbyMatchmakingQueueWorker(ctx context.Context, mode game.ModeIdentifier, queues port.MatchLobbyQueueRepository, dispatcher port.LobbyMatchmakingDispatcher) *LobbyMatchmakingQueueWorker {
	return &LobbyMatchmakingQueueWorker{ctx: ctx, mode: mode, repository: queues, dispatcher: dispatcher}
}

func (w *LobbyMatchmakingQueueWorker) Run() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			go w.Run()
		}
	}()

	var findTicker = time.NewTicker(matchmakingLobbyQueueWorkerFindInterval)
	defer findTicker.Stop()

	var lobbyReapTicker = time.NewTicker(matchmakingLobbyQueueWorkerLobbyReapInterval)
	defer lobbyReapTicker.Stop()

	for {
		select {
		case <-findTicker.C:
			if err := w.findMatches(); err != nil {
				fmt.Println(err)
			}
		case <-lobbyReapTicker.C:
			if err := w.reapLobbies(); err != nil {
				fmt.Println(err)
			}
		}
	}

}

func (w *LobbyMatchmakingQueueWorker) findMatches() error {
	actives, err := w.repository.GetActiveQuests(w.ctx, w.mode)
	if err != nil {
		return err
	}
	for _, active := range actives {
		fmt.Println(fmt.Sprintf("finding matches for quest: %s", active.String()))
		if err := w.findMatch(active); err != nil {
			return err
		}
	}
	return nil
}

func (w *LobbyMatchmakingQueueWorker) findMatch(quest uuid.UUID) error {

	fmt.Println("find match started")
	defer fmt.Println("find match ended")

	lobbies, err := w.repository.GetQueuedLobbies(w.ctx, w.mode, quest)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%d queued lobbies found", len(lobbies)))

	for _, queued := range lobbies {

		found, err := w.repository.FindMatch(w.ctx, w.mode, queued, 5)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if found.Zero() {
			fmt.Println(fmt.Sprintf("found lobby for %s is empty, skipping", queued.LobbyID))
			continue
		}

		fmt.Println(fmt.Sprintf("found lobby for %s: %s", queued.LobbyID, found.LobbyID))

		if err := w.dispatcher.Dispatch(w.ctx, w.mode, quest, []match.LobbyQueueEntry{queued, found}); err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(fmt.Sprintf("dispatch lobby match for %s and %s", queued.LobbyID, found.LobbyID))

		if err := w.repository.RemoveLobbyFromQueue(w.ctx, w.mode, quest, queued.LobbyID); err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(fmt.Sprintf("removed lobby %s from queue", queued.LobbyID))

		if err := w.repository.RemoveLobbyFromQueue(w.ctx, w.mode, quest, found.LobbyID); err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(fmt.Sprintf("removed lobby %s from queue", found.LobbyID))

	}

	return nil

}

func (w *LobbyMatchmakingQueueWorker) reapLobbies() error {
	actives, err := w.repository.GetActiveQuests(w.ctx, w.mode)
	if err != nil {
		return err
	}
	for _, active := range actives {
		count, err := w.repository.RemoveExpiredLobbies(w.ctx, w.mode, active)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Sprintf("%d expired lobbies reaped", count))
	}
	return nil
}
