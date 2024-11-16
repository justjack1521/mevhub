package query

import (
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
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
	InstanceRepository          port.GameInstanceReadRepository
	PlayerParticipantRepository port.PlayerParticipantReadRepository
}

func NewGameSummaryQueryHandler(sessions session.InstanceReadRepository, instances port.GameInstanceReadRepository, players port.PlayerParticipantReadRepository) *GameSummaryQueryHandler {
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
		Seed:         instance.Seed,
		Participants: participants,
	}, nil

}
