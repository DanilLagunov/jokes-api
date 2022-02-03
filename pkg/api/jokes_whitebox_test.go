package api

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/cache/memcache"
	file_storage "github.com/DanilLagunov/jokes-api/pkg/storage/file-storage"
	"github.com/DanilLagunov/jokes-api/pkg/views"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJokes(t *testing.T) {
	storage := file_storage.NewFileStorage("./test-data/test_jokes.json")
	template := views.NewTemptale("../../templates/")
	cache := memcache.NewMemCache(20*time.Second, 1*time.Minute)
	h := NewHandler(storage, template, cache)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes", nil)

	h.getJokes(recorder, req)
	data, err := os.ReadFile("./test-data/index_test.html")
	if err != nil {
		log.Fatal(err)
	}
	assert.EqualValues(t, data, recorder.Body.Bytes())

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestGetFunniestJokes(t *testing.T) {
	storage := file_storage.NewFileStorage("./test-data/test_jokes.json")
	template := views.NewTemptale("../../templates/")
	cache := memcache.NewMemCache(20*time.Second, 1*time.Minute)
	h := NewHandler(storage, template, cache)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/funniest", nil)

	h.getFunniestJokes(recorder, req)

	data, err := os.ReadFile("./test-data/funniest_test.html")
	if err != nil {
		log.Fatal(err)
	}
	assert.EqualValues(t, data, recorder.Body.Bytes())

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestGetRandomJokes(t *testing.T) {
	storage := file_storage.NewFileStorage("./test-data/test_jokes.json")
	template := views.NewTemptale("../../templates/")
	cache := memcache.NewMemCache(20*time.Second, 1*time.Minute)
	h := NewHandler(storage, template, cache)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/random", nil)

	data, err := os.ReadFile("./test-data/index_test.html")
	if err != nil {
		t.Fatal(err)
	}

	h.getRandomJokes(recorder, req)
	if res := bytes.Compare(data, recorder.Body.Bytes()); res == 0 {
		t.Fatal("jokes not in random order")
	}

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestGetJokeByText(t *testing.T) {
	storage := file_storage.NewFileStorage("./test-data/test_jokes.json")
	template := views.NewTemptale("../../templates/")
	cache := memcache.NewMemCache(20*time.Second, 1*time.Minute)
	h := NewHandler(storage, template, cache)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/search", nil)
	q := req.URL.Query()
	q.Add("text", "say")
	req.URL.RawQuery = q.Encode()

	h.getJokesByText(recorder, req)

	data, err := os.ReadFile("./test-data/get-jokes-by-text_test.html")
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, data, recorder.Body.Bytes())

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestGetJokeByID(t *testing.T) {
	storage := file_storage.NewFileStorage("./test-data/test_jokes.json")
	template := views.NewTemptale("../../templates/")
	cache := memcache.NewMemCache(20*time.Second, 1*time.Minute)
	h := NewHandler(storage, template, cache)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/", nil)

	vars := map[string]string{
		"id": "1a7xnd",
	}

	// CHANGE THIS LINE!!!
	req = mux.SetURLVars(req, vars)

	h.getJokeByID(recorder, req)

	data, err := os.ReadFile("./test-data/get-joke-by-id_test.html")
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, data, recorder.Body.Bytes())

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestAddJoke(t *testing.T) {
	storage := file_storage.NewFileStorage("./test-data/test_jokes.json")
	template := views.NewTemptale("../../templates/")
	cache := memcache.NewMemCache(20*time.Second, 1*time.Minute)
	h := NewHandler(storage, template, cache)

	form := url.Values{}
	form.Set("title", "Test Joke")
	form.Set("body", "Test joke body")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/add", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	h.addJoke(recorder, req)
	assert.EqualValues(t, http.StatusFound, recorder.Code)
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
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)

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
