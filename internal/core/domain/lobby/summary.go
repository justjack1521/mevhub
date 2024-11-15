package lobby

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Summary struct {
	InstanceID         uuid.UUID
	QuestID            uuid.UUID
	PartyID            string
	LobbyComment       string
	MinimumPlayerLevel int
	RegisteredAt       time.Time
	Players            []PlayerSlotSummary
}

type PlayerSlotSummary struct {
	PartySlot     int
	Ready         bool
	PlayerSummary PlayerSummary
}

type PlayerSummary struct {
	Identity PlayerIdentity
	Loadout  PlayerLoadout
}

type PlayerIdentity struct {
	PlayerID      uuid.UUID
	PlayerName    string
	PlayerComment string
	PlayerLevel   int
}

type PlayerLoadout struct {
	DeckIndex    int
	JobCard      PlayerJobCardSummary
	Weapon       PlayerWeaponSummary
	AbilityCards []PlayerAbilityCardSummary
}

type PlayerJobCardSummary struct {
	JobCardID      uuid.UUID
	SubJobIndex    int
	CrownLevel     int
	OverBoostLevel int
}

type PlayerWeaponSummary struct {
	WeaponID        uuid.UUID
	SubWeaponUnlock int
}

type PlayerAbilityCardSummary struct {
	AbilityCardID    uuid.UUID
	SlotIndex        int
	AbilityCardLevel int
	AbilityLevel     int
	ExtraSkillUnlock int
	OverBoostLevel   int
}
