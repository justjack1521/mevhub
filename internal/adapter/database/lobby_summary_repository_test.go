package database_test

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"mevhub/internal/adapter/database"
	"testing"
)

func TestLobbySummaryDatabaseRepository_Query(t *testing.T) {

	var repository = database.NewLobbySummaryDatabaseRepository(NewDatabaseConnection())

	lobby, err := repository.QueryByID(context.Background(), uuid.FromStringOrNil("a38b314b-b536-49d8-bde8-0ebc89498938"))

	assert.NoError(t, err)
	assert.NotNil(t, lobby)

}
