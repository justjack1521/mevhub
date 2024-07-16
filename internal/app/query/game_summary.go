package query

import (
	"mevhub/internal/domain/game"
	"mevhub/internal/domain/session"
)

type GameSummaryQuery struct {
}

func (g GameSummaryQuery) CommandName() string {
	return "query.game.summary"
}

func NewGameSummaryQuery() GameSummaryQuery {
	return GameSummaryQuery{}
}

type GameSummaryQueryHandler struct {
	SessionRepository           session.InstanceReadRepository
	InstanceRepository          game.InstanceReadRepository
	PlayerParticipantRepository game.PlayerParticipantReadRepository
}

func NewGameSummaryQueryHandler(sessions session.InstanceReadRepository, instances game.InstanceReadRepository, players game.PlayerParticipantReadRepository) *GameSummaryQueryHandler {
	return &GameSummaryQueryHandler{SessionRepository: sessions, InstanceRepository: instances, PlayerParticipantRepository: players}
}

func (h *GameSummaryQueryHandler) Handle(ctx Context, cmd GameSummaryQuery) (game.Summary, error) {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return game.Summary{}, err
	}

	instance, err := h.InstanceRepository.Get(ctx, current.LobbyID)
	if err != nil {
		return game.Summary{}, err
	}

	participants, err := h.PlayerParticipantRepository.QueryAll(ctx, current.LobbyID)
	if err != nil {
		return game.Summary{}, err
	}

	return game.Summary{
		SysID:        instance.SysID,
		PartyID:      instance.PartyID,
		Seed:         instance.Seed,
		Participants: participants,
	}, nil

}
