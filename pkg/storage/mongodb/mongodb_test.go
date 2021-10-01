package mongodb_test

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	"github.com/DanilLagunov/jokes-api/pkg/storage/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

const requestTimeout time.Duration = time.Second * 2

func TestGetJokes(t *testing.T) {
	db, err := mongodb.NewDatabase()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	expected := []models.Joke{
		{ID: "614d988ed2ab61b5dba36489", Title: "First joke", Body: "Normal", Score: 3},
		{ID: "614d98a7d2ab61b5dba3648b", Title: "Second joke", Body: "Incredible", Score: 35},
		{ID: "6150ed6dc471125ddd1a0912", Title: "Third joke", Body: "Funny", Score: 15},
	}

	result, err := db.GetJokes(ctx)
	require.NoError(t, err)

	assert.EqualValues(t, expected, result)
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
		Expected []models.Joke
		Valid    bool
	}{
		{
			Text: "first",
			Expected: []models.Joke{
				{ID: "614d988ed2ab61b5dba36489", Title: "First joke", Body: "Normal", Score: 3},
			},
			Valid: true,
		},
		{
			Text: "joke",
			Expected: []models.Joke{
				{ID: "614d988ed2ab61b5dba36489", Title: "First joke", Body: "Normal", Score: 3},
				{ID: "614d98a7d2ab61b5dba3648b", Title: "Second joke", Body: "Incredible", Score: 35},
				{ID: "6150ed6dc471125ddd1a0912", Title: "Third joke", Body: "Funny", Score: 15},
			},
			Valid: true,
		},
		{
			Text:     "foo",
			Expected: []models.Joke{},
			Valid:    false,
		},
	}

	for _, tc := range tests {
		result, err := db.GetJokeByText(ctx, tc.Text)
		if tc.Valid {
			require.NoError(t, err)
		}

		assert.EqualValues(t, tc.Expected, result)
	}
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
			ID:       "614d988ed2ab61b5dba36489",
			Expected: models.Joke{ID: "614d988ed2ab61b5dba36489", Title: "First joke", Body: "Normal", Score: 3},
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
		}

		assert.EqualValues(t, tc.Expected, result)
	}
}

func TestGetFunniestJokes(t *testing.T) {
	db, err := mongodb.NewDatabase()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	expected := []models.Joke{
		{ID: "614d98a7d2ab61b5dba3648b", Title: "Second joke", Body: "Incredible", Score: 35},
		{ID: "6150ed6dc471125ddd1a0912", Title: "Third joke", Body: "Funny", Score: 15},
		{ID: "614d988ed2ab61b5dba36489", Title: "First joke", Body: "Normal", Score: 3},
	}

	result, err := db.GetFunniestJokes(ctx)
	require.NoError(t, err)

	assert.EqualValues(t, expected, result)
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

	random, err := db.GetRandomJokes(ctx)
	require.NoError(t, err)

	if reflect.DeepEqual(random, origin) {
		t.Fatal("jokes not in random order")
	}
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

	result, err := db.AddJoke(ctx, expTitle, expBody)
	require.NoError(t, err)

	assert.EqualValues(t, expTitle, result.Title)
	assert.EqualValues(t, expBody, result.Body)

	_, err = db.JokesCollection.DeleteOne(ctx, bson.M{"id": result.ID})
	if err != nil {
		log.Fatal(err)
	}
}
