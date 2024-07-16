package query

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/game"
)

type GameSummaryQuery struct {
	InstanceID uuid.UUID
}

func (g GameSummaryQuery) CommandName() string {
	return "query.game.summary"
}

func NewGameSummaryQuery(id uuid.UUID) GameSummaryQuery {
	return GameSummaryQuery{InstanceID: id}
}

type GameSummaryQueryHandler struct {
	InstanceRepository          game.InstanceReadRepository
	PlayerParticipantRepository game.PlayerParticipantReadRepository
}

func NewGameSummaryQueryHandler(instances game.InstanceReadRepository, players game.PlayerParticipantReadRepository) *GameSummaryQueryHandler {
	return &GameSummaryQueryHandler{InstanceRepository: instances, PlayerParticipantRepository: players}
}

func (h *GameSummaryQueryHandler) Handle(ctx Context, cmd GameSummaryQuery) (game.Summary, error) {

	instance, err := h.InstanceRepository.Get(ctx, cmd.InstanceID)
	if err != nil {
		return game.Summary{}, err
	}

	participants, err := h.PlayerParticipantRepository.QueryAll(ctx, cmd.InstanceID)
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
