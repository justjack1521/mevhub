package memory

import (
	"testing"

	uuid "github.com/satori/go.uuid"

	"mevhub/internal/core/domain/lobby"
)

// Every query a searcher can send must read at least one set the lobby was
// written to, whenever the lobby matches the filters.
func TestLobbySearchKeysMatchAcrossFilters(t *testing.T) {
	var repo = &LobbySearchRedisRepository{}
	var fire, ice, water = uuid.NewV4(), uuid.NewV4(), uuid.NewV4()

	var categorized = lobby.SearchEntry{ModeIdentifier: "COOP_DEFAULT", Level: 1, Categories: []uuid.UUID{fire}}
	var uncategorized = lobby.SearchEntry{ModeIdentifier: "COOP_DEFAULT", Level: 2}

	var cases = []struct {
		name  string
		entry lobby.SearchEntry
		query lobby.SearchQuery
		want  bool
	}{
		{"categorized, all filters", categorized, lobby.SearchQuery{Levels: []int{1, 2, 3}, Categories: []uuid.UUID{fire, ice, water}}, true},
		{"categorized, no categories", categorized, lobby.SearchQuery{Levels: []int{1}}, true},
		{"categorized, no levels", categorized, lobby.SearchQuery{Categories: []uuid.UUID{fire}}, true},
		{"categorized, no filters", categorized, lobby.SearchQuery{}, true},
		{"categorized, other element", categorized, lobby.SearchQuery{Levels: []int{1}, Categories: []uuid.UUID{ice}}, false},
		{"categorized, other level", categorized, lobby.SearchQuery{Levels: []int{2}, Categories: []uuid.UUID{fire}}, false},
		{"uncategorized, all filters", uncategorized, lobby.SearchQuery{Levels: []int{1, 2, 3}, Categories: []uuid.UUID{fire, ice, water}}, true},
		{"uncategorized, no categories", uncategorized, lobby.SearchQuery{Levels: []int{2}}, true},
		{"uncategorized, no levels", uncategorized, lobby.SearchQuery{Categories: []uuid.UUID{ice}}, true},
		{"uncategorized, no filters", uncategorized, lobby.SearchQuery{}, true},
		{"uncategorized, other level", uncategorized, lobby.SearchQuery{Levels: []int{1}}, false},
		{"other mode", categorized, lobby.SearchQuery{ModeIdentifier: "COMP_SOLO"}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.query.ModeIdentifier == "" {
				c.query.ModeIdentifier = c.entry.ModeIdentifier
			}
			var written = make(map[string]bool)
			for _, key := range repo.GenerateKeysFromInstance(c.entry) {
				written[key] = true
			}
			var found bool
			for _, key := range repo.GenerateKeysFromQuery(c.query) {
				found = found || written[key]
			}
			if found != c.want {
				t.Errorf("found = %v, want %v\nwritten: %v\nread: %v", found, c.want, repo.GenerateKeysFromInstance(c.entry), repo.GenerateKeysFromQuery(c.query))
			}
		})
	}
}
