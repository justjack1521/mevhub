package dto

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
)

type PlayerIdentityRedis struct {
	PlayerID      string `redis:"PlayerID"`
	PlayerName    string `redis:"PlayerName"`
	PlayerComment string `redis:"PlayerComment"`
	PlayerLevel   int    `redis:"PlayerLevel"`
}

func (x PlayerIdentityRedis) ToEntity() lobby.PlayerIdentity {
	return lobby.PlayerIdentity{
		PlayerID:      uuid.FromStringOrNil(x.PlayerID),
		PlayerName:    x.PlayerName,
		PlayerComment: x.PlayerComment,
		PlayerLevel:   x.PlayerLevel,
	}
}

func (x PlayerIdentityRedis) ToMapStringInterface() (map[string]interface{}, error) {
	result := map[string]interface{}{
		"PlayerID":      x.PlayerID,
		"PlayerName":    x.PlayerName,
		"PlayerComment": x.PlayerComment,
		"PlayerLevel":   x.PlayerLevel,
	}
	return result, nil
}

type PlayerLoadoutRedis struct {
	DeckIndex      int    `redis:"DeckIndex"`
	JobCardID      string `redis:"JobCardID"`
	SubJobIndex    int    `redis:"SubJobIndex"`
	CrownLevel     int    `redis:"CrownLevel"`
	OverBoostLevel int    `redis:"OverBoostLevel"`

	WeaponID        string `redis:"WeaponID"`
	SubWeaponUnlock int    `redis:"SubWeaponUnlock"`

	AbilityCardsBytes []byte                   `redis:"AbilityCards"`
	AbilityCards      []PlayerAbilityCardRedis `json:"AbilityCards"`
}

func (x PlayerLoadoutRedis) ToEntity() lobby.PlayerLoadout {

	var loadout = lobby.PlayerLoadout{
		DeckIndex: x.DeckIndex,
		JobCard: lobby.PlayerJobCardSummary{
			JobCardID:      uuid.FromStringOrNil(x.JobCardID),
			SubJobIndex:    x.SubJobIndex,
			CrownLevel:     x.CrownLevel,
			OverBoostLevel: x.OverBoostLevel,
		},
		Weapon: lobby.PlayerWeaponSummary{
			WeaponID:        uuid.FromStringOrNil(x.WeaponID),
			SubWeaponUnlock: x.SubWeaponUnlock,
		},
		AbilityCards: nil,
	}

	var cards []PlayerAbilityCardRedis

	if x.AbilityCardsBytes != nil {
		if err := json.Unmarshal(x.AbilityCardsBytes, &cards); err != nil {
			return loadout
		}
	}

	loadout.AbilityCards = make([]lobby.PlayerAbilityCardSummary, len(cards))

	for i, v := range cards {
		loadout.AbilityCards[i] = lobby.PlayerAbilityCardSummary{
			AbilityCardID:    uuid.FromStringOrNil(v.AbilityCardID),
			SlotIndex:        v.SlotIndex,
			AbilityCardLevel: v.AbilityCardLevel,
			AbilityLevel:     v.AbilityLevel,
			OverBoostLevel:   v.OverBoostLevel,
		}
	}

	return loadout

}

func (x PlayerLoadoutRedis) ToMapStringInterface() (map[string]interface{}, error) {
	result := map[string]interface{}{
		"DeckIndex":       x.DeckIndex,
		"JobCardID":       x.JobCardID,
		"SubJobIndex":     x.SubJobIndex,
		"CrownLevel":      x.CrownLevel,
		"OverBoostLevel":  x.OverBoostLevel,
		"WeaponID":        x.WeaponID,
		"SubWeaponUnlock": x.SubWeaponUnlock,
	}
	cards, err := json.Marshal(x.AbilityCards)
	if err != nil {
		return nil, err
	}
	result["AbilityCards"] = cards
	return result, nil
}

type PlayerAbilityCardRedis struct {
	AbilityCardID    string
	SlotIndex        int
	AbilityCardLevel int
	AbilityLevel     int
	OverBoostLevel   int
}

func (x PlayerAbilityCardRedis) MarshalBinary() ([]byte, error) {
	return json.Marshal(x)
}

func (x PlayerAbilityCardRedis) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	return nil
}
