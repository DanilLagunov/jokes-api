package cache_test

import (
	"testing"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/cache"
	"github.com/DanilLagunov/jokes-api/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	tests := []struct {
		Item     models.Joke
		Duration time.Duration
		Expected models.Joke
		Valid    bool
	}{
		{
			Item: models.Joke{
				ID:    "1",
				Title: "First",
				Body:  "first",
				Score: 0,
			},
			Duration: 2 * time.Second,
			Expected: models.Joke{},
			Valid:    false,
		},
		{
			Item: models.Joke{
				ID:    "2",
				Title: "Second",
				Body:  "second",
				Score: 0,
			},
			Duration: 0,
			Expected: models.Joke{},
			Valid:    false,
		},
		{
			Item: models.Joke{
				ID:    "3",
				Title: "Third",
				Body:  "third",
				Score: 0,
			},
			Duration: 2 * time.Minute,
			Expected: models.Joke{
				ID:    "3",
				Title: "Third",
				Body:  "third",
				Score: 0,
			},
			Valid: true,
		},
	}
	cache := cache.NewCache(2*time.Second, 10*time.Second)
	for _, tc := range tests {
		cache.Set(tc.Item.ID, tc.Item, tc.Duration)
		item, err := cache.Get(tc.Item.ID)
		assert.EqualValues(t, tc.Item, item)
		require.NoError(t, err)
	}
	time.Sleep(10 * time.Second)
	for _, tc := range tests {
		item, err := cache.Get(tc.Item.ID)
		assert.EqualValues(t, tc.Expected, item)
		if tc.Valid {
			require.NoError(t, err)
		}
	}
}
