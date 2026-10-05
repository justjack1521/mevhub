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
	errs map[uuid.UUID]error
}

func (f *fakeSearchSummaries) Query(ctx context.Context, id uuid.UUID) (lobby.Summary, error) {
	if err, ok := f.errs[id]; ok {
		return lobby.Summary{}, err
	}
	// One open slot, so the lobby is returned.
	return lobby.Summary{InstanceID: id, Players: []lobby.PlayerSlotSummary{{PartySlot: 1}}}, nil
}

func (f *fakeSearchSummaries) QueryByPartyID(ctx context.Context, party string) (lobby.Summary, error) {
	return lobby.Summary{}, errors.New("not used")
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
