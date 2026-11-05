package query_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mevhub/internal/core/application/query"
	"mevhub/internal/core/domain/lobby"
)

type searchTestContext struct {
	context.Context
}

func (searchTestContext) UserID() uuid.UUID   { return uuid.Nil }
func (searchTestContext) PlayerID() uuid.UUID { return uuid.Nil }

type fakeSearchIndex struct {
	ids    []uuid.UUID
	pruned []uuid.UUID
}

func (f *fakeSearchIndex) Query(ctx context.Context, qry lobby.SearchQuery) ([]lobby.SearchResult, error) {
	var results = make([]lobby.SearchResult, len(f.ids))
	for i, id := range f.ids {
		results[i] = lobby.SearchResult{LobbyID: id}
	}
	return results, nil
}

func (f *fakeSearchIndex) Prune(ctx context.Context, qry lobby.SearchQuery, id uuid.UUID) error {
	f.pruned = append(f.pruned, id)
	return nil
}

type fakeSearchSummaries struct {
	errs     map[uuid.UUID]error
	party    lobby.Summary
	partyErr error
}

func (f *fakeSearchSummaries) Query(ctx context.Context, id uuid.UUID) (lobby.Summary, error) {
	if err, ok := f.errs[id]; ok {
		return lobby.Summary{}, err
	}
	// One open slot, so the lobby is returned.
	return lobby.Summary{InstanceID: id, Players: []lobby.PlayerSlotSummary{{PartySlot: 1}}}, nil
}

func (f *fakeSearchSummaries) QueryByPartyID(ctx context.Context, party string) (lobby.Summary, error) {
	return f.party, f.partyErr
}

func TestSearchLobby_ExpiredLobby_PrunedAndSkipped(t *testing.T) {

	live, expired, failing := uuid.NewV4(), uuid.NewV4(), uuid.NewV4()

	index := &fakeSearchIndex{ids: []uuid.UUID{live, expired, failing}}
	summaries := &fakeSearchSummaries{errs: map[uuid.UUID]error{
		// Wrapped the way the summary service and Redis repository wrap it.
		expired: fmt.Errorf("failed to query lobby summary: %w", lobby.ErrFailedQueryLobbyInstance(fmt.Errorf("%w by key: x", lobby.ErrLobbyInstanceNotFound))),
		// A real failure must not cost the lobby its place in the index.
		failing: errors.New("redis unavailable"),
	}}

	handler := query.NewSearchLobbyQueryHandler(index, summaries)
	results, err := handler.Handle(searchTestContext{context.Background()}, query.NewSearchLobbyQuery(lobby.SearchQuery{}, ""))
	require.NoError(t, err)

	require.Len(t, results, 1)
	assert.Equal(t, live, results[0].InstanceID)
	assert.Equal(t, []uuid.UUID{expired}, index.pruned)
}

func TestSearchLobby_PartyID(t *testing.T) {

	var notFound = fmt.Errorf("failed to query lobby party x summary: %w", lobby.ErrFailedQueryLobbyInstance(fmt.Errorf("%w by key: x", lobby.ErrLobbyInstanceNotFound)))
	var full = lobby.Summary{InstanceID: uuid.NewV4(), Players: []lobby.PlayerSlotSummary{{PartySlot: 0, PlayerSummary: lobby.PlayerSummary{Identity: lobby.PlayerIdentity{PlayerID: uuid.NewV4()}}}}}

	var cases = []struct {
		name    string
		repo    *fakeSearchSummaries
		want    int
		wantErr bool
	}{
		{"unknown party is empty", &fakeSearchSummaries{partyErr: notFound}, 0, false},
		{"full lobby still returned", &fakeSearchSummaries{party: full}, 1, false},
		{"real failure is an error", &fakeSearchSummaries{partyErr: errors.New("redis unavailable")}, 0, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Filters are set to show they are ignored.
			var qry = lobby.SearchQuery{ModeIdentifier: "other", Levels: []int{9}, PlayerLevel: 1}
			handler := query.NewSearchLobbyQueryHandler(&fakeSearchIndex{}, c.repo)
			results, err := handler.Handle(searchTestContext{context.Background()}, query.NewSearchLobbyQuery(qry, "PARTY"))
			if c.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, results)
			assert.Len(t, results, c.want)
		})
	}
}
