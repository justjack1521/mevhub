package game

import (
	uuid "github.com/satori/go.uuid"
)

type PlayerParticipant struct {
	UserID     uuid.UUID
	PlayerID   uuid.UUID
	PlayerSlot int
	BotControl bool
	Loadout    PlayerLoadout
}

type PlayerLoadout struct {
	PlayerID     uuid.UUID
	PlayerName   string
	DeckIndex    int
	JobCard      PlayerJobCardLoadout
	Weapon       PlayerWeaponLoadout
	AbilityCards []PlayerAbilityCardLoadout
}

func (x PlayerLoadout) CalculateDeckLevel() int {
	var level = x.JobCard.CrownLevel
	for _, card := range x.AbilityCards {
		level += card.AbilityCardLevel
	}
	return level
}

type PlayerJobCardLoadout struct {
	JobCardID         uuid.UUID
	SubJobIndex       int
	HPStatMod         int
	AttackStatMod     int
	BreakStatMod      int
	MagicStatMod      int
	SpeedStatMod      int
	DefenseStatMod    int
	CritChanceStatMod int
	UltimateBoost     int
	OverBoostLevel    int
	CrownLevel        int
	AutoAbilities     map[uuid.UUID]int
}

type PlayerWeaponLoadout struct {
	WeaponID          uuid.UUID
	SubWeaponUnlock   int
	HPStatMod         int
	AttackStatMod     int
	BreakStatMod      int
	MagicStatMod      int
	SpeedStatMod      int
	DefenseStatMod    int
	CritChanceStatMod int
	UltimateBoost     int
	AutoAbilities     map[uuid.UUID]int
}

type PlayerAbilityCardLoadout struct {
	AbilityCardID    uuid.UUID
	SlotIndex        int
	AbilityCardLevel int
	AbilityLevel     int
	ExtraSkillUnlock int
	OverBoostLevel   int
	AutoAbilities    map[uuid.UUID]int
}
