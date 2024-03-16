package memory_test

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"mevhub/internal/adapter/memory"
	"mevhub/internal/domain/game"
	"mevhub/internal/domain/lobby"
	"testing"
)

func TestLobbyInstanceRepository_GenerateKeysFromInstance(t *testing.T) {

	var repository = memory.LobbySearchRedisRepository{}

	type test struct {
		name     string
		instance lobby.Instance
		expected []string
	}

	var tests = []test{
		{
			name: "Lobby Instance With no category",
			instance: lobby.Instance{
				ModeIdentifier: game.ModeIdentifierCoopDefault,
				Level:          3,
				Categories:     nil,
			},
			expected: []string{"multi_coop_default_lobby:3"},
		},
		{
			name: "Lobby instance With single category",
			instance: lobby.Instance{
				ModeIdentifier: game.ModeIdentifierCoopDefault,
				Level:          2,
				Categories:     []uuid.UUID{uuid.FromStringOrNil("9f27d88b-4629-4be5-8163-8e1726705605")},
			},
			expected: []string{"multi_coop_default_lobby:2:9f27d88b-4629-4be5-8163-8e1726705605"},
		},
		{
			name: "Lobby Instance with multiple categories",
			instance: lobby.Instance{
				ModeIdentifier: game.ModeIdentifierCoopDefault,
				Level:          1,
				Categories: []uuid.UUID{
					uuid.FromStringOrNil("a27f190c-5efe-4f79-a4ee-a457de7e6122"),
					uuid.FromStringOrNil("de50edef-c6d8-475a-a7fc-253a884ce8f9"),
				},
			},
			expected: []string{
				"multi_coop_default_lobby:1:a27f190c-5efe-4f79-a4ee-a457de7e6122",
				"multi_coop_default_lobby:1:de50edef-c6d8-475a-a7fc-253a884ce8f9",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var keys = repository.GenerateKeysFromInstance(tc.instance)
			assert.Equal(t, len(tc.expected), len(keys))
			for index, key := range keys {
				assert.Equal(t, tc.expected[index], key)
			}
		})
	}

}

func TestLobbyInstanceRepository_GenerateKeysFromQuery(t *testing.T) {

	var repository = memory.LobbySearchRedisRepository{}

	type test struct {
		name     string
		query    lobby.SearchQuery
		expected []string
	}

	var tests = []test{
		{
			name: "Query with no level and no categories",
			query: lobby.SearchQuery{
				ModeIdentifier: game.ModeIdentifierCoopDefault,
			},
			expected: []string{"multi_coop_default_lobby"},
		},
		{
			name: "Query with single level and no categories",
			query: lobby.SearchQuery{
				ModeIdentifier: game.ModeIdentifierCoopDefault,
				Levels:         []int{3},
			},
			expected: []string{"multi_coop_default_lobby:3"},
		},
		{
			name: "Query with single level and single category",
			query: lobby.SearchQuery{
				ModeIdentifier: game.ModeIdentifierCoopDefault,
				Levels:         []int{3},
				Categories:     []uuid.UUID{uuid.FromStringOrNil("b17fd078-4789-4582-b36b-ad13b02acd42")},
			},
			expected: []string{"multi_coop_default_lobby:3:b17fd078-4789-4582-b36b-ad13b02acd42"},
		},
		{
			name: "Query with multiple levels and single category",
			query: lobby.SearchQuery{
				ModeIdentifier: game.ModeIdentifierCoopDefault,
				Levels:         []int{3, 4},
				Categories: []uuid.UUID{
					uuid.FromStringOrNil("cdf345c2-7223-4dfb-9f71-fd94e71f3b29"),
				},
			},
			expected: []string{
				"multi_coop_default_lobby:3:cdf345c2-7223-4dfb-9f71-fd94e71f3b29",
				"multi_coop_default_lobby:4:cdf345c2-7223-4dfb-9f71-fd94e71f3b29",
			},
		},
		{
			name: "Query with single level and multiple categories",
			query: lobby.SearchQuery{
				ModeIdentifier: game.ModeIdentifierCoopDefault,
				Levels:         []int{3},
				Categories: []uuid.UUID{
					uuid.FromStringOrNil("dab886fe-44ff-4883-8be1-0c676c440aad"),
					uuid.FromStringOrNil("f3963ca5-d1b1-4330-9d54-c85678e0ad60"),
				},
			},
			expected: []string{
				"multi_coop_default_lobby:3:dab886fe-44ff-4883-8be1-0c676c440aad",
				"multi_coop_default_lobby:3:f3963ca5-d1b1-4330-9d54-c85678e0ad60",
			},
		},
		{
			name: "Query with multiple levels and multiple categories",
			query: lobby.SearchQuery{
				ModeIdentifier: game.ModeIdentifierCoopDefault,
				Levels:         []int{3, 4},
				Categories: []uuid.UUID{
					uuid.FromStringOrNil("11df4207-1013-4bec-88ee-dfbf54be447a"),
					uuid.FromStringOrNil("a9b9185d-a5f9-4a55-9ddf-de67080eea7f"),
				},
			},
			expected: []string{
				"multi_coop_default_lobby:3:11df4207-1013-4bec-88ee-dfbf54be447a",
				"multi_coop_default_lobby:3:a9b9185d-a5f9-4a55-9ddf-de67080eea7f",
				"multi_coop_default_lobby:4:11df4207-1013-4bec-88ee-dfbf54be447a",
				"multi_coop_default_lobby:4:a9b9185d-a5f9-4a55-9ddf-de67080eea7f",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var keys = repository.GenerateKeysFromQuery(tc.query)
			assert.Equal(t, len(tc.expected), len(keys))
			for index, key := range keys {
				assert.Equal(t, tc.expected[index], key)
			}
		})
	}

}
