package mongodb_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	"github.com/DanilLagunov/jokes-api/pkg/storage/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const requestTimeout time.Duration = time.Second * 2

var jokeID string

func TestMain(m *testing.M) {
	prepareDBForTests()
	os.Exit(m.Run())
}

func prepareDBForTests() {
	db, err := mongodb.NewDatabase()
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	db.AddJoke(ctx, "First joke", "Normal", 3)
	db.AddJoke(ctx, "Second joke", "Incredible", 35)
	joke, err := db.AddJoke(ctx, "Third joke", "Funny", 15)
	if err != nil {
		log.Fatal()
	}
	jokeID = joke.ID
}

func TestGetJokes(t *testing.T) {
	db, err := mongodb.NewDatabase()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	result, amount, err := db.GetJokes(ctx, 0, 3)
	require.NoError(t, err)

	assert.EqualValues(t, "First joke", result[0].Title)
	assert.EqualValues(t, 3, amount)

	result, amount, err = db.GetJokes(ctx, 1, 4)
	require.NoError(t, err)

	assert.EqualValues(t, "Second joke", result[0].Title)
	assert.EqualValues(t, 3, amount)

	result, _, err = db.GetJokes(ctx, 4, 4)
	require.NoError(t, err)

	assert.EqualValues(t, []models.Joke{}, result)
}

func TestGetJokeByText(t *testing.T) {
	db, err := mongodb.NewDatabase()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	tests := []struct {
		Text     string
		Expected []string
		Valid    bool
	}{
		{
			Text:     "first",
			Expected: []string{"First joke"},
			Valid:    true,
		},
		{
			Text:     "joke",
			Expected: []string{"First joke", "Second joke", "Third joke"},
			Valid:    true,
		},
		{
			Text:     "foo",
			Expected: []string{},
			Valid:    false,
		},
	}

	for _, tc := range tests {
		result, amount, err := db.GetJokesByText(ctx, 0, 3, tc.Text)
		if tc.Valid {
			require.NoError(t, err)
		}
		assert.EqualValues(t, amount, len(result))
		for i := 0; i < amount; i++ {
			assert.EqualValues(t, tc.Expected[i], result[i].Title)
		}
	}

	result, _, err := db.GetJokesByText(ctx, 4, 4, "first")
	require.NoError(t, err)

	assert.EqualValues(t, []models.Joke{}, result)
}

func TestGetJokeByID(t *testing.T) {
	db, err := mongodb.NewDatabase()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	tests := []struct {
		ID       string
		Expected models.Joke
		Valid    bool
	}{
		{
			ID:       jokeID,
			Expected: models.Joke{ID: jokeID, Title: "Third joke", Body: "Funny", Score: 15},
			Valid:    true,
		},
		{
			ID:       "123456",
			Expected: models.Joke{},
			Valid:    false,
		},
	}

	for _, tc := range tests {
		result, err := db.GetJokeByID(ctx, tc.ID)
		if tc.Valid {
			require.NoError(t, err)
			assert.EqualValues(t, tc.Expected, result)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetFunniestJokes(t *testing.T) {
	db, err := mongodb.NewDatabase()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	expected := []int{35, 15, 3}

	result, amount, err := db.GetFunniestJokes(ctx, 0, 3)
	require.NoError(t, err)

	for i := 0; i < amount; i++ {
		assert.EqualValues(t, expected[i], result[i].Score)
	}
	assert.EqualValues(t, 3, amount)

	result, _, err = db.GetFunniestJokes(ctx, 4, 4)
	require.NoError(t, err)

	assert.EqualValues(t, []models.Joke{}, result)
}

func TestGetRandomJokes(t *testing.T) {
	db, err := mongodb.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	origin := []models.Joke{
		{ID: "614d988ed2ab61b5dba36489", Title: "First joke", Body: "Normal", Score: 3},
		{ID: "614d98a7d2ab61b5dba3648b", Title: "Second joke", Body: "Incredible", Score: 35},
		{ID: "6150ed6dc471125ddd1a0912", Title: "Third joke", Body: "Funny", Score: 15},
	}

	random, size, err := db.GetRandomJokes(ctx, 3)
	require.NoError(t, err)

	pass := false
	for i := 0; i < 100; i++ {
		if reflect.DeepEqual(random, origin) == false {
			pass = true
			break
		}
	}
	if pass == false {
		t.Fatal("jokes not in random order")
	}
	assert.EqualValues(t, 3, size)
}

func TestAddJoke(t *testing.T) {
	db, err := mongodb.NewDatabase()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	expTitle := "Added joke"
	expBody := "New"

	result, err := db.AddJoke(ctx, expTitle, expBody, 0)
	require.NoError(t, err)

	assert.EqualValues(t, expTitle, result.Title)
	assert.EqualValues(t, expBody, result.Body)
}
