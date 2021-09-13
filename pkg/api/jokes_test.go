package api

import (
	"bufio"
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
	data, err := os.Open("../../templates/index_test.html")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	reader := bufio.NewReader(data)
	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	assert.EqualValues(t, buffer, recorder.Body)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
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

	data, err := os.Open("../../templates/funniest_test.html")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	reader := bufio.NewReader(data)
	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	assert.EqualValues(t, buffer, recorder.Body)

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

	data, err := os.Open("../../templates/index_test.html")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	reader := bufio.NewReader(data)
	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		buffer.Write(part[:count])
	}

	h.getRandomJokes(recorder, req)
	if buffer == recorder.Body {
		t.Fatal("jokes not in random order")
	}

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
	q.Add("text", "say")
	req.URL.RawQuery = q.Encode()

	h.getJoke(recorder, req)

	data, err := os.Open("../../templates/get-jokes-by-text_test.html")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	reader := bufio.NewReader(data)
	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	assert.EqualValues(t, buffer, recorder.Body)

	assert.EqualValues(t, http.StatusOK, recorder.Code)

	req = httptest.NewRequest(http.MethodGet, "/jokes/search", nil)
	q = req.URL.Query()
	q.Add("id", "1a7xnd")
	req.URL.RawQuery = q.Encode()

	recorder.Body.Reset()
	h.getJoke(recorder, req)

	data, err = os.Open("../../templates/get-joke-by-id_test.html")
	if err != nil {
		log.Fatal(err)
	}

	reader = bufio.NewReader(data)
	buffer = bytes.NewBuffer(make([]byte, 0))
	part = make([]byte, 1024)
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	assert.EqualValues(t, buffer, recorder.Body)

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
