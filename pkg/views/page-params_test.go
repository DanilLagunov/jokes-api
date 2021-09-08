package views

import (
	"testing"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatePageParams(t *testing.T) {
	content := [100]models.Joke{}
	tests := []struct {
		SrcSkip  int
		SrcSeed  int
		Expected JokesPageParams
	}{
		{0, 0, JokesPageParams{0, 0, 0, 0, []models.Joke{}, 0, 0}},
		{20, 20, JokesPageParams{20, 20, 2, 5, content[20:40], 40, 0}},
		{100, 20, JokesPageParams{100, 20, 6, 5, content[100:], 120, 80}},
		{101, 20, JokesPageParams{101, 20, 0, 0, []models.Joke{}, 0, 0}},
	}

	for _, tc := range tests {
		pageParams := CreatePageParams(tc.SrcSkip, tc.SrcSeed, content[:])
		assert.EqualValues(t, tc.Expected, pageParams, "params are not equal")
	}
}
