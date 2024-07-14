package translate

import (
	"github.com/justjack1521/mevium/pkg/genproto/protoidentity"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
)

type LobbySummaryTranslator Translator[lobby.Summary, *protomulti.ProtoLobbySummary]
type LobbyPlayerSlotSummaryTranslator Translator[lobby.PlayerSlotSummary, *protomulti.ProtoLobbyPlayerSlot]
type LobbyPlayerSummaryTranslator Translator[lobby.PlayerSummary, *protomulti.ProtoLobbyPlayer]
type LobbyIdentityTranslator Translator[lobby.PlayerIdentity, *protoidentity.ProtoPlayerIdentity]
type LobbyLoadoutTranslator Translator[lobby.PlayerLoadout, *protoidentity.ProtoPlayerLoadoutIdentity]

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
	identity LobbyIdentityTranslator
	loadout  LobbyLoadoutTranslator
}

func NewLobbyPlayerSummaryTranslator() LobbyPlayerSummaryTranslator {
	return lobbyPlayerSummaryTranslator{
		identity: NewLobbyPlayerIdentityTranslator(),
		loadout:  NewLobbyPlayerLoadoutTranslator(),
	}
}

func (t lobbyPlayerSummaryTranslator) Marshall(data lobby.PlayerSummary) (out *protomulti.ProtoLobbyPlayer, err error) {

	identity, err := t.identity.Marshall(data.Identity)
	if err != nil {
		return nil, err
	}

	loadout, err := t.loadout.Marshall(data.Loadout)
	if err != nil {
		return nil, err
	}

	return &protomulti.ProtoLobbyPlayer{
		Identity: identity,
		Loadout:  loadout,
	}, nil

}

func (t lobbyPlayerSummaryTranslator) Unmarshall(data *protomulti.ProtoLobbyPlayer) (out lobby.PlayerSummary, err error) {
	identity, err := t.identity.Unmarshall(data.Identity)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	loadout, err := t.loadout.Unmarshall(data.Loadout)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	return lobby.PlayerSummary{
		Identity: identity,
		Loadout:  loadout,
	}, nil
}

type lobbyPlayerIdentityTranslator struct {
}

func NewLobbyPlayerIdentityTranslator() LobbyIdentityTranslator {
	return lobbyPlayerIdentityTranslator{}
}

func (t lobbyPlayerIdentityTranslator) Marshall(data lobby.PlayerIdentity) (out *protoidentity.ProtoPlayerIdentity, err error) {
	return &protoidentity.ProtoPlayerIdentity{
		PlayerId:      data.PlayerID.String(),
		PlayerName:    data.PlayerName,
		PlayerLevel:   int32(data.PlayerLevel),
		PlayerComment: data.PlayerComment,
	}, nil
}

func (t lobbyPlayerIdentityTranslator) Unmarshall(data *protoidentity.ProtoPlayerIdentity) (out lobby.PlayerIdentity, err error) {
	return lobby.PlayerIdentity{
		PlayerID:      uuid.FromStringOrNil(data.PlayerId),
		PlayerName:    data.PlayerName,
		PlayerComment: data.PlayerComment,
		PlayerLevel:   int(data.PlayerLevel),
	}, nil
}

type lobbyPlayerLoadoutTranslator struct {
}

func NewLobbyPlayerLoadoutTranslator() LobbyLoadoutTranslator {
	return lobbyPlayerLoadoutTranslator{}
}

func (t lobbyPlayerLoadoutTranslator) Marshall(data lobby.PlayerLoadout) (out *protoidentity.ProtoPlayerLoadoutIdentity, err error) {
	var result = &protoidentity.ProtoPlayerLoadoutIdentity{
		JobCard: &protoidentity.ProtoPlayerJobIdentity{
			JobCardId:      data.JobCard.JobCardID.String(),
			SubJobIndex:    int32(data.JobCard.SubJobIndex),
			CrownLevel:     int32(data.JobCard.CrownLevel),
			OverBoostLevel: int32(data.JobCard.OverBoostLevel),
		},
		Weapon: &protoidentity.ProtoPlayerWeaponIdentity{
			WeaponId:        data.Weapon.WeaponID.String(),
			SubWeaponUnlock: int32(data.Weapon.SubWeaponUnlock),
		},
		AbilityCards: make([]*protoidentity.ProtoAbilityCardIdentity, len(data.AbilityCards)),
	}

	for index, value := range data.AbilityCards {
		result.AbilityCards[index] = &protoidentity.ProtoAbilityCardIdentity{
			AbilityCardId:    value.AbilityCardID.String(),
			AbilityCardLevel: int32(value.AbilityCardLevel),
			AbilityLevel:     int32(value.AbilityLevel),
			ExtraSkillUnlock: int32(value.ExtraSkillUnlock),
			OverBoostLevel:   int32(value.OverBoostLevel),
			SlotIndex:        int32(value.SlotIndex),
		}
	}

	return result, nil

}

func (t lobbyPlayerLoadoutTranslator) Unmarshall(data *protoidentity.ProtoPlayerLoadoutIdentity) (out lobby.PlayerLoadout, err error) {

	var result = lobby.PlayerLoadout{
		DeckIndex: 0,
		JobCard: lobby.PlayerJobCardSummary{
			JobCardID:      uuid.FromStringOrNil(data.JobCard.JobCardId),
			SubJobIndex:    int(data.JobCard.SubJobIndex),
			CrownLevel:     int(data.JobCard.CrownLevel),
			OverBoostLevel: int(data.JobCard.SubJobIndex),
		},
		Weapon: lobby.PlayerWeaponSummary{
			WeaponID:        uuid.FromStringOrNil(data.Weapon.WeaponId),
			SubWeaponUnlock: int(data.Weapon.SubWeaponUnlock),
		},
		AbilityCards: make([]lobby.PlayerAbilityCardSummary, len(data.AbilityCards)),
	}

	for index, value := range data.AbilityCards {
		result.AbilityCards[index] = lobby.PlayerAbilityCardSummary{
			AbilityCardID:    uuid.FromStringOrNil(value.AbilityCardId),
			SlotIndex:        int(value.SlotIndex),
			AbilityCardLevel: int(value.AbilityCardLevel),
			AbilityLevel:     int(value.AbilityLevel),
			ExtraSkillUnlock: int(value.ExtraSkillUnlock),
			OverBoostLevel:   int(value.OverBoostLevel),
		}
	}

	return result, nil

}
