package api

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	file_storage "github.com/DanilLagunov/jokes-api/pkg/storage/file-storage"
	"github.com/DanilLagunov/jokes-api/pkg/views"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJokes(t *testing.T) {
	storage := file_storage.NewFileStorage("./test-data/test_jokes.json")
	template := views.NewTemptale("../../templates/")
	h := NewHandler(storage, template)

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
	h := NewHandler(storage, template)

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
	h := NewHandler(storage, template)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/random", nil)

	data, err := os.ReadFile("./test-data/index_test.html")
	if err != nil {
		log.Fatal(err)
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
	h := NewHandler(storage, template)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/search", nil)
	q := req.URL.Query()
	q.Add("text", "say")
	req.URL.RawQuery = q.Encode()

	h.getJoke(recorder, req)

	data, err := os.ReadFile("./test-data/get-jokes-by-text_test.html")
	if err != nil {
		log.Fatal(err)
	}

	assert.EqualValues(t, data, recorder.Body.Bytes())

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestGetJokeByID(t *testing.T) {
	storage := file_storage.NewFileStorage("./test-data/test_jokes.json")
	template := views.NewTemptale("../../templates/")
	h := NewHandler(storage, template)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jokes/search", nil)
	q := req.URL.Query()
	q.Add("id", "1a7xnd")
	req.URL.RawQuery = q.Encode()

	h.getJoke(recorder, req)

	data, err := os.ReadFile("./test-data/get-joke-by-id_test.html")
	if err != nil {
		log.Fatal(err)
	}

	assert.EqualValues(t, data, recorder.Body.Bytes())

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func TestAddJoke(t *testing.T) {
	storage := file_storage.NewFileStorage("./test-data/test_jokes.json")
	template := views.NewTemptale("../../templates/")
	h := NewHandler(storage, template)

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
