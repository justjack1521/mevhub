package database_test

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"mevhub/internal/adapter/database"
	"testing"
)

func TestGameQuestDatabaseRepository_QueryByID(t *testing.T) {

	var id = uuid.FromStringOrNil("d10c9ef8-91cd-4aa1-b6ac-893f92e1e63a")

	var db = NewDatabaseConnection()

	var repo = database.NewGameQuestDatabaseRepository(db)

	result, err := repo.QueryByID(id)

	assert.NoError(t, err)
	assert.False(t, result.Zero())
	assert.False(t, result.Tier.Zero())
	assert.False(t, result.Tier.GameMode.Zero())
	assert.True(t, len(result.Categories) > 0)
	for _, category := range result.Categories {
		assert.False(t, category.Zero())
	}

}
