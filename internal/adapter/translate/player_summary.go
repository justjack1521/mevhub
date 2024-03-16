package translate

import (
	"github.com/justjack1521/mevium/pkg/genproto/protoplayer"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/game"
)

type GamePlayerSummaryTranslator Translator[game.PlayerSummary, *protoplayer.ProtoPlayerInfo]
type GameJobCardSummaryTranslator Translator[game.PlayerJobCardSummary, *protoplayer.ProtoJobCard]
type GameWeaponSummaryTranslator Translator[game.PlayerWeaponSummary, *protoplayer.ProtoWeapon]
type GameAbilityCardSummaryTranslator Translator[game.PlayerAbilityCardSummary, *protoplayer.ProtoAbilityCard]

type gamePlayerSummaryTranslator struct {
	jobCardTranslator     GameJobCardSummaryTranslator
	weaponTranslator      GameWeaponSummaryTranslator
	abilityCardTranslator GameAbilityCardSummaryTranslator
}

func NewGamePlayerSummaryTranslator() GamePlayerSummaryTranslator {
	return gamePlayerSummaryTranslator{
		jobCardTranslator:     gameJobCardSummaryTranslator{},
		weaponTranslator:      gameWeaponSummaryTranslator{},
		abilityCardTranslator: gameAbilityCardSummaryTranslator{},
	}
}

func (t gamePlayerSummaryTranslator) Marshall(data game.PlayerSummary) (out *protoplayer.ProtoPlayerInfo, err error) {
	//TODO implement me
	panic("implement me")
}

func (t gamePlayerSummaryTranslator) Unmarshall(data *protoplayer.ProtoPlayerInfo) (out game.PlayerSummary, err error) {

	job, err := t.jobCardTranslator.Unmarshall(data.Loadout.JobCard)
	if err != nil {
		return game.PlayerSummary{}, err
	}

	weapon, err := t.weaponTranslator.Unmarshall(data.Loadout.Weapon)
	if err != nil {
		return game.PlayerSummary{}, err
	}

	var cards = make([]game.PlayerAbilityCardSummary, len(data.Loadout.AbilityCards))

	for index, value := range data.Loadout.AbilityCards {
		card, err := t.abilityCardTranslator.Unmarshall(value)
		if err != nil {
			return game.PlayerSummary{}, err
		}
		cards[index] = card
	}

	var result = game.PlayerSummary{
		PlayerID:      uuid.FromStringOrNil(data.PlayerId),
		PlayerName:    data.PlayerName,
		PlayerLevel:   int(data.PlayerLevel),
		PlayerComment: "",
		DeckIndex:     int(data.Loadout.DeckIndex),
		JobCard:       job,
		Weapon:        weapon,
		AbilityCards:  cards,
	}

	return result, nil
}

type gameJobCardSummaryTranslator struct {
}

func (t gameJobCardSummaryTranslator) Marshall(data game.PlayerJobCardSummary) (out *protoplayer.ProtoJobCard, err error) {
	//TODO implement me
	panic("implement me")
}

func (t gameJobCardSummaryTranslator) Unmarshall(data *protoplayer.ProtoJobCard) (out game.PlayerJobCardSummary, err error) {
	var result = game.PlayerJobCardSummary{
		JobCardID:         uuid.FromStringOrNil(data.BaseJobId),
		SubJobIndex:       int(data.SubJobIndex),
		HPStatMod:         int(data.HpStatMod),
		AttackStatMod:     int(data.AttackStatMod),
		BreakStatMod:      int(data.BreakStatMod),
		MagicStatMod:      int(data.MagicStatMod),
		SpeedStatMod:      int(data.SpeedStatMod),
		DefenseStatMod:    int(data.DefenseStatMod),
		CritChanceStatMod: int(data.CritChanceStatMod),
		UltimateBoost:     int(data.UltimateBoost),
		OverBoostLevel:    int(data.OverBoostLevel),
		AutoAbilities:     make(map[uuid.UUID]int),
		CrownLevel:        int(data.CrownLevel),
	}
	for key, ability := range data.AutoAbilities {
		result.AutoAbilities[uuid.FromStringOrNil(key)] = int(ability)
	}
	return result, nil
}

type gameWeaponSummaryTranslator struct {
}

func (t gameWeaponSummaryTranslator) Marshall(data game.PlayerWeaponSummary) (out *protoplayer.ProtoWeapon, err error) {
	//TODO implement me
	panic("implement me")
}

func (t gameWeaponSummaryTranslator) Unmarshall(data *protoplayer.ProtoWeapon) (out game.PlayerWeaponSummary, err error) {
	var result = game.PlayerWeaponSummary{
		WeaponID:          uuid.FromStringOrNil(data.BaseWeaponId),
		SubWeaponUnlock:   int(data.SubWeaponUnlock),
		HPStatMod:         int(data.HpStatMod),
		AttackStatMod:     int(data.AttackStatMod),
		BreakStatMod:      int(data.BreakStatMod),
		MagicStatMod:      int(data.MagicStatMod),
		SpeedStatMod:      int(data.SpeedStatMod),
		DefenseStatMod:    int(data.DefenseStatMod),
		CritChanceStatMod: int(data.CritChanceStatMod),
		UltimateBoost:     int(data.UltimateBoost),
		AutoAbilities:     make(map[uuid.UUID]int),
	}
	for key, ability := range data.AutoAbilities {
		result.AutoAbilities[uuid.FromStringOrNil(key)] = int(ability)
	}
	return result, nil
}

type gameAbilityCardSummaryTranslator struct {
}

func (t gameAbilityCardSummaryTranslator) Marshall(data game.PlayerAbilityCardSummary) (out *protoplayer.ProtoAbilityCard, err error) {
	//TODO implement me
	panic("implement me")
}

func (t gameAbilityCardSummaryTranslator) Unmarshall(data *protoplayer.ProtoAbilityCard) (out game.PlayerAbilityCardSummary, err error) {
	var result = game.PlayerAbilityCardSummary{
		AbilityCardID:    uuid.FromStringOrNil(data.AbilityCardId),
		SlotIndex:        int(data.SlotIndex),
		AbilityCardLevel: int(data.AbilityCardLevel),
		AbilityLevel:     int(data.AbilityLevel),
		ExtraSkillUnlock: int(data.ExtraSkillUnlock),
		OverBoostLevel:   int(data.OverBoostLevel),
		AutoAbilities:    make(map[uuid.UUID]int),
	}
	for key, ability := range data.AutoAbilities {
		result.AutoAbilities[uuid.FromStringOrNil(key)] = int(ability)
	}
	return result, nil
}
