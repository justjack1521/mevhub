package game

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Summary struct {
	InstanceID   uuid.UUID
	QuestID      uuid.UUID
	Seed         int64
	State        InstanceState
	RegisteredAt time.Time
}

type PlayerSummary struct {
	PlayerID      uuid.UUID
	PlayerName    string
	PlayerComment string
	PlayerLevel   int
	DeckIndex     int
	JobCard       PlayerJobCardSummary
	Weapon        PlayerWeaponSummary
	AbilityCards  []PlayerAbilityCardSummary
}

type PlayerJobCardSummary struct {
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

type PlayerWeaponSummary struct {
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

type PlayerAbilityCardSummary struct {
	AbilityCardID    uuid.UUID
	SlotIndex        int
	AbilityCardLevel int
	AbilityLevel     int
	ExtraSkillUnlock int
	OverBoostLevel   int
	AutoAbilities    map[uuid.UUID]int
}
