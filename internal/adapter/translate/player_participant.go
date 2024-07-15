package translate

import (
	"github.com/justjack1521/mevium/pkg/genproto/protoidentity"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/game"
)

type GamePlayerParticipantTranslator Translator[game.PlayerParticipant, *protomulti.ProtoGameParticipant]
type GamePlayerLoadoutTranslator Translator[game.PlayerLoadout, *protoidentity.ProtoPlayerLoadout]
type GameJobCardLoadoutTranslator Translator[game.PlayerJobCardLoadout, *protoidentity.ProtoPlayerJobLoadout]
type GameWeaponLoadoutTranslator Translator[game.PlayerWeaponLoadout, *protoidentity.ProtoPlayerWeaponLoadout]
type GameAbilityCardLoadoutTranslator Translator[game.PlayerAbilityCardLoadout, *protoidentity.ProtoPlayerAbilityCardLoadout]

type gamePlayerParticipantTranslator struct {
	loadout GamePlayerLoadoutTranslator
}

func NewGameParticipantTranslator() GamePlayerParticipantTranslator {
	return gamePlayerParticipantTranslator{
		loadout: NewGamePlayerLoadoutTranslator(),
	}
}

func (f gamePlayerParticipantTranslator) Marshall(data game.PlayerParticipant) (out *protomulti.ProtoGameParticipant, err error) {
	loadout, err := f.loadout.Marshall(data.Loadout)
	return &protomulti.ProtoGameParticipant{
		PartySlot:  int32(data.PlayerSlot),
		BotControl: data.BotControl,
		Loadout:    loadout,
	}, nil
}

func (f gamePlayerParticipantTranslator) Unmarshall(data *protomulti.ProtoGameParticipant) (out game.PlayerParticipant, err error) {
	loadout, err := f.loadout.Unmarshall(data.Loadout)
	if err != nil {
		return game.PlayerParticipant{}, err
	}
	return game.PlayerParticipant{
		PlayerSlot: int(data.PartySlot),
		BotControl: data.BotControl,
		Loadout:    loadout,
	}, nil
}

type gamePlayerLoadoutTranslator struct {
	jobCardTranslator     GameJobCardLoadoutTranslator
	weaponTranslator      GameWeaponLoadoutTranslator
	abilityCardTranslator GameAbilityCardLoadoutTranslator
}

func NewGamePlayerLoadoutTranslator() GamePlayerLoadoutTranslator {
	return gamePlayerLoadoutTranslator{
		jobCardTranslator:     gameJobCardLoadoutTranslator{},
		weaponTranslator:      gameWeaponLoadoutTranslator{},
		abilityCardTranslator: gameAbilityCardLoadoutTranslator{},
	}
}

func (t gamePlayerLoadoutTranslator) Marshall(data game.PlayerLoadout) (out *protoidentity.ProtoPlayerLoadout, err error) {

	job, err := t.jobCardTranslator.Marshall(data.JobCard)
	if err != nil {
		return nil, err
	}

	weapon, err := t.weaponTranslator.Marshall(data.Weapon)
	if err != nil {
		return nil, err
	}

	var cards = make([]*protoidentity.ProtoPlayerAbilityCardLoadout, len(data.AbilityCards))

	for index, value := range data.AbilityCards {
		card, err := t.abilityCardTranslator.Marshall(value)
		if err != nil {
			return nil, err
		}
		cards[index] = card
	}

	return &protoidentity.ProtoPlayerLoadout{
		PlayerId:     data.PlayerID.String(),
		PlayerName:   data.PlayerName,
		DeckIndex:    int32(data.DeckIndex),
		Job:          job,
		Weapon:       weapon,
		AbilityCards: cards,
	}, nil

}

func (t gamePlayerLoadoutTranslator) Unmarshall(data *protoidentity.ProtoPlayerLoadout) (out game.PlayerLoadout, err error) {

	job, err := t.jobCardTranslator.Unmarshall(data.Job)
	if err != nil {
		return game.PlayerLoadout{}, err
	}

	weapon, err := t.weaponTranslator.Unmarshall(data.Weapon)
	if err != nil {
		return game.PlayerLoadout{}, err
	}

	var cards = make([]game.PlayerAbilityCardLoadout, len(data.AbilityCards))

	for index, value := range data.AbilityCards {
		card, err := t.abilityCardTranslator.Unmarshall(value)
		if err != nil {
			return game.PlayerLoadout{}, err
		}
		cards[index] = card
	}

	var result = game.PlayerLoadout{
		PlayerID:     uuid.FromStringOrNil(data.PlayerId),
		PlayerName:   data.PlayerName,
		DeckIndex:    int(data.DeckIndex),
		JobCard:      job,
		Weapon:       weapon,
		AbilityCards: cards,
	}

	return result, nil
}

type gameJobCardLoadoutTranslator struct {
}

func (t gameJobCardLoadoutTranslator) Marshall(data game.PlayerJobCardLoadout) (out *protoidentity.ProtoPlayerJobLoadout, err error) {
	var result = &protoidentity.ProtoPlayerJobLoadout{
		Identity: &protoidentity.ProtoPlayerJobIdentity{
			JobCardId:      data.JobCardID.String(),
			SubJobIndex:    int32(data.SubJobIndex),
			CrownLevel:     int32(data.CrownLevel),
			OverBoostLevel: int32(data.OverBoostLevel),
		},
		Stat: &protoidentity.ProtoPlayerJobStat{
			HpStatMod:         int32(data.HPStatMod),
			AttackStatMod:     int32(data.AttackStatMod),
			BreakStatMod:      int32(data.BreakStatMod),
			MagicStatMod:      int32(data.MagicStatMod),
			SpeedStatMod:      int32(data.SpeedStatMod),
			DefenseStatMod:    int32(data.DefenseStatMod),
			CritChanceStatMod: int32(data.CritChanceStatMod),
			UltimateBoost:     int32(data.UltimateBoost),
			AutoAbilities:     make(map[string]int32),
		},
	}
	for key, ability := range data.AutoAbilities {
		result.Stat.AutoAbilities[key.String()] = int32(ability)
	}
	return result, nil
}

func (t gameJobCardLoadoutTranslator) Unmarshall(data *protoidentity.ProtoPlayerJobLoadout) (out game.PlayerJobCardLoadout, err error) {
	var result = game.PlayerJobCardLoadout{
		JobCardID:         uuid.FromStringOrNil(data.Identity.JobCardId),
		SubJobIndex:       int(data.Identity.SubJobIndex),
		HPStatMod:         int(data.Stat.HpStatMod),
		AttackStatMod:     int(data.Stat.AttackStatMod),
		BreakStatMod:      int(data.Stat.BreakStatMod),
		MagicStatMod:      int(data.Stat.MagicStatMod),
		SpeedStatMod:      int(data.Stat.SpeedStatMod),
		DefenseStatMod:    int(data.Stat.DefenseStatMod),
		CritChanceStatMod: int(data.Stat.CritChanceStatMod),
		UltimateBoost:     int(data.Stat.UltimateBoost),
		OverBoostLevel:    int(data.Identity.OverBoostLevel),
		AutoAbilities:     make(map[uuid.UUID]int),
		CrownLevel:        int(data.Identity.CrownLevel),
	}
	for key, ability := range data.Stat.AutoAbilities {
		result.AutoAbilities[uuid.FromStringOrNil(key)] = int(ability)
	}
	return result, nil
}

type gameWeaponLoadoutTranslator struct {
}

func (t gameWeaponLoadoutTranslator) Marshall(data game.PlayerWeaponLoadout) (out *protoidentity.ProtoPlayerWeaponLoadout, err error) {
	var result = &protoidentity.ProtoPlayerWeaponLoadout{
		Identity: &protoidentity.ProtoPlayerWeaponIdentity{
			WeaponId:        data.WeaponID.String(),
			SubWeaponUnlock: int32(data.SubWeaponUnlock),
		},
		Stat: &protoidentity.ProtoPlayerWeaponStat{
			HpStatMod:         int32(data.HPStatMod),
			AttackStatMod:     int32(data.AttackStatMod),
			BreakStatMod:      int32(data.BreakStatMod),
			MagicStatMod:      int32(data.MagicStatMod),
			SpeedStatMod:      int32(data.SpeedStatMod),
			DefenseStatMod:    int32(data.DefenseStatMod),
			CritChanceStatMod: int32(data.CritChanceStatMod),
			UltimateBoost:     int32(data.UltimateBoost),
			AutoAbilities:     make(map[string]int32),
		},
	}
	for key, ability := range data.AutoAbilities {
		result.Stat.AutoAbilities[key.String()] = int32(ability)
	}
	return result, nil
}

func (t gameWeaponLoadoutTranslator) Unmarshall(data *protoidentity.ProtoPlayerWeaponLoadout) (out game.PlayerWeaponLoadout, err error) {
	var result = game.PlayerWeaponLoadout{
		WeaponID:          uuid.FromStringOrNil(data.Identity.WeaponId),
		SubWeaponUnlock:   int(data.Identity.SubWeaponUnlock),
		HPStatMod:         int(data.Stat.HpStatMod),
		AttackStatMod:     int(data.Stat.AttackStatMod),
		BreakStatMod:      int(data.Stat.BreakStatMod),
		MagicStatMod:      int(data.Stat.MagicStatMod),
		SpeedStatMod:      int(data.Stat.SpeedStatMod),
		DefenseStatMod:    int(data.Stat.DefenseStatMod),
		CritChanceStatMod: int(data.Stat.CritChanceStatMod),
		UltimateBoost:     int(data.Stat.UltimateBoost),
		AutoAbilities:     make(map[uuid.UUID]int),
	}
	for key, ability := range data.Stat.AutoAbilities {
		result.AutoAbilities[uuid.FromStringOrNil(key)] = int(ability)
	}
	return result, nil
}

type gameAbilityCardLoadoutTranslator struct {
}

func (t gameAbilityCardLoadoutTranslator) Marshall(data game.PlayerAbilityCardLoadout) (out *protoidentity.ProtoPlayerAbilityCardLoadout, err error) {
	var result = &protoidentity.ProtoPlayerAbilityCardLoadout{
		Identity: &protoidentity.ProtoAbilityCardIdentity{
			AbilityCardId:    data.AbilityCardID.String(),
			AbilityCardLevel: int32(data.AbilityCardLevel),
			AbilityLevel:     int32(data.AbilityLevel),
			ExtraSkillUnlock: int32(data.ExtraSkillUnlock),
			OverBoostLevel:   int32(data.OverBoostLevel),
			SlotIndex:        int32(data.SlotIndex),
		},
		Stat: &protoidentity.ProtoAbilityCardStat{AutoAbilities: make(map[string]int32)},
	}
	for key, ability := range data.AutoAbilities {
		result.Stat.AutoAbilities[key.String()] = int32(ability)
	}
	return result, nil
}

func (t gameAbilityCardLoadoutTranslator) Unmarshall(data *protoidentity.ProtoPlayerAbilityCardLoadout) (out game.PlayerAbilityCardLoadout, err error) {
	var result = game.PlayerAbilityCardLoadout{
		AbilityCardID:    uuid.FromStringOrNil(data.Identity.AbilityCardId),
		SlotIndex:        int(data.Identity.SlotIndex),
		AbilityCardLevel: int(data.Identity.AbilityCardLevel),
		AbilityLevel:     int(data.Identity.AbilityLevel),
		ExtraSkillUnlock: int(data.Identity.ExtraSkillUnlock),
		OverBoostLevel:   int(data.Identity.OverBoostLevel),
		AutoAbilities:    make(map[uuid.UUID]int),
	}
	for key, ability := range data.Stat.AutoAbilities {
		result.AutoAbilities[uuid.FromStringOrNil(key)] = int(ability)
	}
	return result, nil
}
