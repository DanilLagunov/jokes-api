package chancache_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/cache/chancache"
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
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cache := chancache.NewChannelCache(ctx, 2*time.Second, 3*time.Second)
	time.Sleep(2*time.Second + 80*time.Millisecond)
	for _, tc := range tests {
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				cache.Set(tc.Item.ID, tc.Item, tc.Duration)
				item, err := cache.Get(tc.Item.ID)
				assert.EqualValues(t, tc.Item, item)
				require.NoError(t, err)
			}()
		}
		wg.Wait()
	}
	time.Sleep(4 * time.Second)
	for _, tc := range tests {
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				item, err := cache.Get(tc.Item.ID)
				if tc.Valid {
					require.NoError(t, err)
					assert.EqualValues(t, tc.Item, item)
					return
				}
				assert.EqualValues(t, models.Joke{}, item)
				assert.EqualError(t, err, "key not found")
			}()
		}
		wg.Wait()
	}
	return
}
