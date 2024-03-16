package translate

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
)

type LobbySummaryTranslator Translator[lobby.Summary, *protomulti.ProtoLobbySummary]
type LobbyPlayerSlotSummaryTranslator Translator[lobby.PlayerSlotSummary, *protomulti.ProtoLobbyPlayerSlot]
type LobbyPlayerSummaryTranslator Translator[lobby.PlayerSummary, *protomulti.ProtoLobbyPlayer]
type LobbyJobCardTranslator Translator[lobby.PlayerJobCardSummary, *protomulti.ProtoLobbyJobCard]
type LobbyWeaponTranslator Translator[lobby.PlayerWeaponSummary, *protomulti.ProtoLobbyWeapon]
type LobbyAbilityCardTranslator Translator[lobby.PlayerAbilityCardSummary, *protomulti.ProtoLobbyAbilityCard]

type lobbySummaryTranslator struct {
	slot LobbyPlayerSlotSummaryTranslator
}

func NewLobbySummaryTranslator() LobbySummaryTranslator {
	return lobbySummaryTranslator{slot: NewLobbyPlayerSlotSummaryTranslator()}
}

func (t lobbySummaryTranslator) Marshall(data lobby.Summary) (out *protomulti.ProtoLobbySummary, err error) {
	var result = &protomulti.ProtoLobbySummary{
		InstanceId:         data.InstanceID.String(),
		QuestId:            data.QuestID.String(),
		Comment:            data.LobbyComment,
		MinimumPlayerLevel: int32(data.MinimumPlayerLevel),
		Players:            make([]*protomulti.ProtoLobbyPlayerSlot, len(data.Players)),
	}
	for index, value := range data.Players {
		player, err := t.slot.Marshall(value)
		if err != nil {
			return nil, err
		}
		result.Players[index] = player
	}
	return result, nil
}

func (t lobbySummaryTranslator) Unmarshall(data *protomulti.ProtoLobbySummary) (out lobby.Summary, err error) {
	var result = lobby.Summary{
		InstanceID:         uuid.FromStringOrNil(data.InstanceId),
		QuestID:            uuid.FromStringOrNil(data.QuestId),
		LobbyComment:       data.Comment,
		MinimumPlayerLevel: int(data.MinimumPlayerLevel),
		Players:            make([]lobby.PlayerSlotSummary, len(data.Players)),
	}
	for index, value := range data.Players {
		player, err := t.slot.Unmarshall(value)
		if err != nil {
			return lobby.Summary{}, err
		}
		result.Players[index] = player
	}
	return result, nil
}

type lobbyPlayerSlotSummaryTranslator struct {
	player LobbyPlayerSummaryTranslator
}

func (t lobbyPlayerSlotSummaryTranslator) Marshall(data lobby.PlayerSlotSummary) (out *protomulti.ProtoLobbyPlayerSlot, err error) {

	player, err := t.player.Marshall(data.PlayerSummary)
	if err != nil {
		return nil, err
	}
	var result = &protomulti.ProtoLobbyPlayerSlot{
		SlotIndex: int32(data.PartySlot),
		Ready:     data.Ready,
		Player:    player,
	}
	return result, nil
}

func (t lobbyPlayerSlotSummaryTranslator) Unmarshall(data *protomulti.ProtoLobbyPlayerSlot) (out lobby.PlayerSlotSummary, err error) {
	player, err := t.player.Unmarshall(data.Player)
	if err != nil {
		return lobby.PlayerSlotSummary{}, err
	}
	var result = lobby.PlayerSlotSummary{
		PartySlot:     int(data.SlotIndex),
		Ready:         data.Ready,
		PlayerSummary: player,
	}
	return result, nil
}

func NewLobbyPlayerSlotSummaryTranslator() LobbyPlayerSlotSummaryTranslator {
	return lobbyPlayerSlotSummaryTranslator{player: NewLobbyPlayerSummaryTranslator()}
}

type lobbyPlayerSummaryTranslator struct {
	job    LobbyJobCardTranslator
	weapon LobbyWeaponTranslator
	card   LobbyAbilityCardTranslator
}

func NewLobbyPlayerSummaryTranslator() LobbyPlayerSummaryTranslator {
	return lobbyPlayerSummaryTranslator{
		job:    lobbyJobCardTranslator{},
		weapon: lobbyWeaponTranslator{},
		card:   lobbyAbilityCardTranslator{},
	}
}

func (t lobbyPlayerSummaryTranslator) Marshall(data lobby.PlayerSummary) (out *protomulti.ProtoLobbyPlayer, err error) {

	job, err := t.job.Marshall(data.Loadout.JobCard)
	if err != nil {
		return nil, err
	}

	weapon, err := t.weapon.Marshall(data.Loadout.Weapon)
	if err != nil {
		return nil, err
	}

	var cards = make([]*protomulti.ProtoLobbyAbilityCard, len(data.Loadout.AbilityCards))
	for index, value := range data.Loadout.AbilityCards {
		card, err := t.card.Marshall(value)
		if err != nil {
			return nil, err
		}
		cards[index] = card
	}

	var result = &protomulti.ProtoLobbyPlayer{
		PlayerId:      data.Identity.PlayerID.String(),
		PlayerName:    data.Identity.PlayerName,
		PlayerComment: data.Identity.PlayerComment,
		PlayerLevel:   int32(data.Identity.PlayerLevel),
		DeckIndex:     int32(data.Loadout.DeckIndex),
		JobCard:       job,
		Weapon:        weapon,
		AbilityCards:  cards,
	}
	return result, nil
}

func (t lobbyPlayerSummaryTranslator) Unmarshall(data *protomulti.ProtoLobbyPlayer) (out lobby.PlayerSummary, err error) {

	job, err := t.job.Unmarshall(data.JobCard)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	weapon, err := t.weapon.Unmarshall(data.Weapon)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	var cards = make([]lobby.PlayerAbilityCardSummary, len(data.AbilityCards))
	for index, value := range data.AbilityCards {
		card, err := t.card.Unmarshall(value)
		if err != nil {
			return lobby.PlayerSummary{}, err
		}
		cards[index] = card
	}

	var result = lobby.PlayerSummary{
		Identity: lobby.PlayerIdentity{
			PlayerID:      uuid.FromStringOrNil(data.PlayerId),
			PlayerName:    data.PlayerName,
			PlayerComment: data.PlayerComment,
			PlayerLevel:   int(data.PlayerLevel),
		},
		Loadout: lobby.PlayerLoadout{
			DeckIndex:    int(data.DeckIndex),
			JobCard:      job,
			Weapon:       weapon,
			AbilityCards: cards,
		},
	}
	return result, nil
}

var (
	ErrFailedUnmarshalLobbyJobCard = func(err error) error {
		return fmt.Errorf("failed to unmarshall lobby job card: %w", err)
	}
)

type lobbyJobCardTranslator struct {
}

func (t lobbyJobCardTranslator) Marshall(data lobby.PlayerJobCardSummary) (out *protomulti.ProtoLobbyJobCard, err error) {
	var result = &protomulti.ProtoLobbyJobCard{
		JobCardId:      data.JobCardID.String(),
		SubJobIndex:    int32(data.SubJobIndex),
		OverBoostLevel: int32(data.OverBoostLevel),
		CrownLevel:     int32(data.CrownLevel),
	}
	return result, nil
}

func (t lobbyJobCardTranslator) Unmarshall(data *protomulti.ProtoLobbyJobCard) (out lobby.PlayerJobCardSummary, err error) {
	id, err := uuid.FromString(data.JobCardId)
	if err != nil {
		return lobby.PlayerJobCardSummary{}, ErrFailedUnmarshalLobbyJobCard(err)
	}
	var result = lobby.PlayerJobCardSummary{
		JobCardID:      id,
		SubJobIndex:    int(data.SubJobIndex),
		OverBoostLevel: int(data.OverBoostLevel),
	}
	return result, nil
}

var (
	ErrFailedUnmarshallLobbyWeapon = func(err error) error {
		return fmt.Errorf("failed to unmarshall lobby weapon: %w", err)
	}
)

type lobbyWeaponTranslator struct {
}

func (t lobbyWeaponTranslator) Marshall(data lobby.PlayerWeaponSummary) (out *protomulti.ProtoLobbyWeapon, err error) {
	var result = &protomulti.ProtoLobbyWeapon{
		WeaponId:        data.WeaponID.String(),
		SubWeaponUnlock: int32(data.SubWeaponUnlock),
	}
	return result, nil
}

func (t lobbyWeaponTranslator) Unmarshall(data *protomulti.ProtoLobbyWeapon) (out lobby.PlayerWeaponSummary, err error) {
	id, err := uuid.FromString(data.WeaponId)
	if err != nil {
		return lobby.PlayerWeaponSummary{}, ErrFailedUnmarshallLobbyWeapon(err)
	}
	var result = lobby.PlayerWeaponSummary{
		WeaponID:        id,
		SubWeaponUnlock: int(data.SubWeaponUnlock),
	}
	return result, nil
}

type lobbyAbilityCardTranslator struct {
}

func (t lobbyAbilityCardTranslator) Marshall(data lobby.PlayerAbilityCardSummary) (out *protomulti.ProtoLobbyAbilityCard, err error) {
	var result = &protomulti.ProtoLobbyAbilityCard{
		AbilityCardId:    data.AbilityCardID.String(),
		AbilityCardLevel: int32(data.AbilityCardLevel),
		AbilityLevel:     int32(data.AbilityLevel),
		OverBoostLevel:   int32(data.OverBoostLevel),
	}
	return result, nil
}

func (t lobbyAbilityCardTranslator) Unmarshall(data *protomulti.ProtoLobbyAbilityCard) (out lobby.PlayerAbilityCardSummary, err error) {
	var result = lobby.PlayerAbilityCardSummary{
		AbilityCardID:    uuid.FromStringOrNil(data.AbilityCardId),
		SlotIndex:        int(data.SlotIndex),
		AbilityCardLevel: int(data.AbilityCardLevel),
		AbilityLevel:     int(data.AbilityLevel),
		OverBoostLevel:   int(data.OverBoostLevel),
	}
	return result, nil
}
