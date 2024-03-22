package memory_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"mevhub/internal/adapter/memory"
	"testing"
)

func TestNewRedisConnection(t *testing.T) {
	var ctx = context.Background()
	_, err := memory.NewRedisConnection(ctx)
	assert.NoError(t, err)
}
