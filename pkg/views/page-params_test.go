package views_test

import (
	"testing"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	"github.com/DanilLagunov/jokes-api/pkg/views"
	"github.com/stretchr/testify/assert"
)

func TestCreatePageParams(t *testing.T) {
	content := [100]models.Joke{}
	tests := []struct {
		SrcSkip  int
		SrcSeed  int
		Expected views.JokesPageParams
	}{
		{0, 0, views.JokesPageParams{0, 0, 0, 0, []models.Joke{}, 0, 0}},
		{20, 20, views.JokesPageParams{20, 20, 2, 5, content[:], 40, 0}},
		{99, 20, views.JokesPageParams{99, 20, 5, 5, content[:1], 119, 79}},
		{101, 20, views.JokesPageParams{101, 20, 0, 0, []models.Joke{}, 0, 0}},
	}

	for _, tc := range tests {
		pageParams := views.CreatePageParams(tc.SrcSkip, tc.SrcSeed, 100, content[:])
		assert.EqualValues(t, tc.Expected, pageParams, "params are not equal")
	}
}
