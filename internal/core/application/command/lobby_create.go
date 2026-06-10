package command

import (
	"fmt"
	"math/rand"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
	"time"

	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
)

type LobbyCreateCommand struct {
	BasicCommand
	LobbyID   uuid.UUID
	QuestID   uuid.UUID
	PartyID   string
	DeckIndex int
	Comment   string
	Options   CreateLobbyOptions
}

type CreateLobbyOptions struct {
	MinimumPlayerLevel int
	Restrictions       []lobby.PlayerSlotRestriction
}

func (c LobbyCreateCommand) CommandName() string {
	return "lobby.create"
}

func NewLobbyCreateCommand(quest uuid.UUID, deck int, comment string, options CreateLobbyOptions) *LobbyCreateCommand {
	return &LobbyCreateCommand{
		LobbyID:   uuid.NewV4(),
		QuestID:   quest,
		PartyID:   fmt.Sprintf("%08d", rand.Intn(100000000)),
		DeckIndex: deck,
		Comment:   comment,
		Options:   options,
	}
}

type LobbyCreateCommandHandler struct {
	EventPublisher          *mevent.Publisher
	SessionRepository       port.SessionInstanceRepository
	InstanceRepository      port.LobbyInstanceWriteRepository
	QuestRepository         port.QuestRepository
	ParticipantFactory      lobby.ParticipantFactory
	ParticipantRepository   port.LobbyParticipantWriteRepository
	SummaryRepository       port.LobbySummaryWriteRepository
	SearchRepository        port.LobbySearchWriteRepository
	PlayerQueueRepository   port.MatchLobbyPlayerQueueWriteRepository
	PlayerSummaryRepository port.LobbyPlayerSummaryReadRepository
	ListenerRepository      lobby.NotificationListenerWriteRepository
	ChannelOpener           port.LobbyNotificationChannelOpener
}

func NewLobbyCreateCommandHandler(
	publisher *mevent.Publisher,
	sessions port.SessionInstanceRepository,
	instances port.LobbyInstanceWriteRepository,
	quests port.QuestRepository,
	participants port.LobbyParticipantWriteRepository,
	summaries port.LobbySummaryWriteRepository,
	search port.LobbySearchWriteRepository,
	playerQueue port.MatchLobbyPlayerQueueWriteRepository,
	playerSummaries port.LobbyPlayerSummaryReadRepository,
	listeners lobby.NotificationListenerWriteRepository,
	channelOpener port.LobbyNotificationChannelOpener,
) *LobbyCreateCommandHandler {
	return &LobbyCreateCommandHandler{
		EventPublisher:          publisher,
		SessionRepository:       sessions,
		InstanceRepository:      instances,
		QuestRepository:         quests,
		ParticipantFactory:      lobby.ParticipantFactory{},
		ParticipantRepository:   participants,
		SummaryRepository:       summaries,
		SearchRepository:        search,
		PlayerQueueRepository:   playerQueue,
		PlayerSummaryRepository: playerSummaries,
		ListenerRepository:      listeners,
		ChannelOpener:           channelOpener,
	}
}

func (h *LobbyCreateCommandHandler) Handle(ctx Context, cmd *LobbyCreateCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	if current.CanJoinLobby() == false {
		return lobby.ErrPlayerAlreadyInLobby(current.LobbyID)
	}

	quest, err := h.QuestRepository.QueryByID(cmd.QuestID)
	if err != nil {
		return err
	}

	var factory = lobby.NewInstanceFactory(ctx, ctx.UserID(), ctx.PlayerID())

	var opts = lobby.InstanceFactoryOptions{
		QuestID:            quest.SysID,
		PlayerSlots:        quest.Tier.GameMode.MaxPlayers,
		MinimumPlayerLevel: cmd.Options.MinimumPlayerLevel,
		SlotRestrictions:   make(map[int]lobby.PlayerSlotRestriction),
	}

	for _, value := range cmd.Options.Restrictions {
		opts.SlotRestrictions[value.Index] = value
	}

	instance, err := factory.Create(cmd.LobbyID, cmd.PartyID, opts)
	if err != nil {
		return err
	}

	if err := h.InstanceRepository.Create(ctx, instance); err != nil {
		return err
	}

	var summary = lobby.Summary{
		InstanceID:         instance.SysID,
		QuestID:            quest.SysID,
		PartyID:            cmd.PartyID,
		LobbyComment:       cmd.Comment,
		MinimumPlayerLevel: instance.MinimumPlayerLevel,
	}
	if err := h.SummaryRepository.Create(ctx, summary); err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod == game.FulfillMethodSearch {
		var categories = make([]uuid.UUID, len(quest.Categories))
		for i, category := range quest.Categories {
			if category.Zero() {
				continue
			}
			categories[i] = category.SysID
		}
		var search = lobby.SearchEntry{
			InstanceID:         instance.SysID,
			ModeIdentifier:     string(quest.Tier.GameMode.ModeIdentifier),
			Level:              quest.Tier.StarLevel,
			MinimumPlayerLevel: instance.MinimumPlayerLevel,
			Categories:         categories,
		}
		if err := h.SearchRepository.Create(ctx, search); err != nil {
			return err
		}
	}

	if quest.Tier.GameMode.FulfillMethod == game.FulfillMethodMatch {
		hostSummary, err := h.PlayerSummaryRepository.Query(ctx, ctx.PlayerID())
		if err != nil {
			return err
		}
		var entry = match.LobbyQueueEntry{
			LobbyID:      instance.SysID,
			QuestID:      instance.QuestID,
			AverageLevel: hostSummary.Loadout.CalculateDeckLevel(),
			JoinedAt:     time.Now().UTC(),
		}
		if err := h.PlayerQueueRepository.AddLobbyToQueue(ctx, quest.Tier.GameMode.ModeIdentifier, entry); err != nil {
			return err
		}
	}

	h.ChannelOpener.Open(ctx, instance.SysID)

	for i := 0; i < instance.PlayerSlotCount; i++ {

		var part = lobby.ParticipantJoinOptions{
			RoleID:     uuid.Nil,
			SlotIndex:  i,
			DeckIndex:  0,
			UseStamina: false,
		}

		var user = uuid.Nil
		var player = uuid.Nil

		if i == 0 {
			part.DeckIndex = cmd.DeckIndex
			part.UseStamina = true
			user = ctx.UserID()
			player = ctx.PlayerID()
		}

		participant, err := h.ParticipantFactory.Create(user, player, instance, opts.SlotRestrictions[i], part)
		if err != nil {
			return err
		}

		if err := h.ParticipantRepository.Create(ctx, participant); err != nil {
			return err
		}

		if uuid.Equal(player, uuid.Nil) == false {
			current.LobbyID = participant.LobbyID
			current.PartySlot = participant.PlayerSlot
			if err := h.SessionRepository.Update(ctx, current); err != nil {
				return err
			}

			if err := h.ListenerRepository.CreateListener(ctx, instance.SysID, ctx.UserID()); err != nil {
				return err
			}

			cmd.QueueEvent(lobby.NewParticipantCreatedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.DeckIndex, participant.PlayerSlot))
		}

	}

	return nil

}
