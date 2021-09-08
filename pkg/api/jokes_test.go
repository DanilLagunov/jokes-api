package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	file_storage "github.com/DanilLagunov/jokes-api/pkg/storage/file-storage"
	"github.com/DanilLagunov/jokes-api/pkg/views"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJokes(t *testing.T) {
	storage := file_storage.NewFileStorage("../storage/file-storage/test_jokes.json")
	templates := [7]string{"../../templates/index.html",
		"../../templates/get-joke-by-id.html",
		"../../templates/get-jokes-by-text.html",
		"../../templates/random.html",
		"../../templates/funniest.html",
		"../../templates/header.html",
		"../../templates/footer.html"}
	template := views.NewTemptale(templates)
	h := NewHandler(storage, template)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes", nil)

	h.getJokes(recorder, req)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	// b, err := io.ReadAll(recorder.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// t.Log(string(b))
}

func TestGetFunniestJokes(t *testing.T) {
	storage := file_storage.NewFileStorage("../storage/file-storage/test_jokes.json")
	templates := [7]string{"../../templates/index.html",
		"../../templates/get-joke-by-id.html",
		"../../templates/get-jokes-by-text.html",
		"../../templates/random.html",
		"../../templates/funniest.html",
		"../../templates/header.html",
		"../../templates/footer.html"}
	template := views.NewTemptale(templates)
	h := NewHandler(storage, template)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/funniest", nil)

	h.getFunniestJokes(recorder, req)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestGetRandomJokes(t *testing.T) {
	storage := file_storage.NewFileStorage("../storage/file-storage/test_jokes.json")
	templates := [7]string{"../../templates/index.html",
		"../../templates/get-joke-by-id.html",
		"../../templates/get-jokes-by-text.html",
		"../../templates/random.html",
		"../../templates/funniest.html",
		"../../templates/header.html",
		"../../templates/footer.html"}
	template := views.NewTemptale(templates)
	h := NewHandler(storage, template)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/random", nil)

	h.getRandomJokes(recorder, req)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestGetJoke(t *testing.T) {
	storage := file_storage.NewFileStorage("../storage/file-storage/test_jokes.json")
	templates := [7]string{"../../templates/index.html",
		"../../templates/get-joke-by-id.html",
		"../../templates/get-jokes-by-text.html",
		"../../templates/random.html",
		"../../templates/funniest.html",
		"../../templates/header.html",
		"../../templates/footer.html"}
	template := views.NewTemptale(templates)
	h := NewHandler(storage, template)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/search", nil)
	q := req.URL.Query()
	q.Add("text", "legs")
	req.URL.RawQuery = q.Encode()

	h.getJoke(recorder, req)
	assert.EqualValues(t, http.StatusOK, recorder.Code)

	req = httptest.NewRequest(http.MethodGet, "/jokes/search", nil)
	q = req.URL.Query()
	q.Add("id", "5tz0ef")
	req.URL.RawQuery = q.Encode()

	h.getJoke(recorder, req)
	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestAddJoke(t *testing.T) {
	storage := file_storage.NewFileStorage("../storage/file-storage/test_jokes.json")
	templates := [7]string{"../../templates/index.html",
		"../../templates/get-joke-by-id.html",
		"../../templates/get-jokes-by-text.html",
		"../../templates/random.html",
		"../../templates/funniest.html",
		"../../templates/header.html",
		"../../templates/footer.html"}
	template := views.NewTemptale(templates)
	h := NewHandler(storage, template)

	data := url.Values{}
	data.Set("title", "Test Joke")
	data.Set("body", "Test joke body")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/random", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	h.getRandomJokes(recorder, req)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestGetPaginationParams(t *testing.T) {
	tests := []struct {
		SrcSkip      string
		SrcSeed      string
		ExpectedSkip int
		ExpectedSeed int
		Valid        bool
	}{
		{"20", "20", 20, 20, true},
		{"", "", 0, 20, true},
		{"-30", "20", 0, 0, false},
		{"-dasdad", "20", 0, 0, false},
	}
	for _, tc := range tests {
		req, _ := http.NewRequest("GET", "/", nil)

		q := req.URL.Query()
		q.Add("skip", tc.SrcSkip)
		q.Add("seed", tc.SrcSkip)
		req.URL.RawQuery = q.Encode()

		skip, seed, err := getPaginationParams(req)
		if tc.Valid {
			require.NoError(t, err)
		}
		assert.EqualValues(t, tc.ExpectedSkip, skip, "wrong skip")
		assert.EqualValues(t, tc.ExpectedSeed, seed, "wrong seed")
	}
}
