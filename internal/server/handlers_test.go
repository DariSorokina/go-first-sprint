package server

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/cookie"
	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, clientID int, requestBody io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, requestBody)
	require.NoError(t, err)

	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if clientID != 0 {
		_, clientIDcookie := cookie.СreateCookieClientID("test")
		req.AddCookie(clientIDcookie)
	}

	result, err := client.Do(req)
	require.NoError(t, err)
	defer result.Body.Close()

	resultBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	return result, string(resultBody)
}

func TestRouter(t *testing.T) {
	flagConfig := config.ParseFlags()
	var l *logger.Logger
	var err error
	if l, err = logger.CreateLogger(flagConfig.FlagLogLevel); err != nil {
		log.Fatal("Failed to create logger:", err)
	}

	storage, err := storage.SetStorage(flagConfig, l)
	if err != nil {
		panic(err)
	}

	if flagConfig.FlagFileStoragePath != "" || flagConfig.FlagPostgresqlDSN != "" {
		defer storage.Close()
	}

	app := app.NewApp(storage, l)
	serv := NewServer(app, flagConfig, l)
	testServer := httptest.NewServer(serv.newRouter())
	defer testServer.Close()

	type expectedData struct {
		expectedContentType string
		expectedStatusCode  int
		expectedBody        string
		expectedLocation    string
	}

	testCases := []struct {
		name         string
		method       string
		clientID     int
		requestBody  io.Reader
		requestPath  string
		expectedData expectedData
	}{
		{
			name:        "handler: ShortenerHandler, test: StatusCreated",
			method:      http.MethodPost,
			clientID:    1,
			requestBody: bytes.NewBuffer([]byte("https://practicum.yandex.ru/")),
			requestPath: "",
			expectedData: expectedData{
				expectedContentType: "text/plain",
				expectedStatusCode:  http.StatusConflict,
				expectedBody:        "http://localhost:8080/d41d8cd98f",
				expectedLocation:    "",
			},
		},
		{
			name:        "handler: OriginalHandler, test: StatusTemporaryRedirect",
			method:      http.MethodGet,
			clientID:    1,
			requestBody: nil,
			requestPath: "/d41d8cd98f",
			expectedData: expectedData{
				expectedContentType: "",
				expectedStatusCode:  http.StatusTemporaryRedirect,
				expectedBody:        "",
				expectedLocation:    "https://practicum.yandex.ru/",
			},
		},
		{
			name:        "handler: shortenerHandlerJSON, test: StatusCreated",
			method:      http.MethodPost,
			clientID:    1,
			requestBody: bytes.NewBuffer([]byte("{\"url\":\"https://practicum.yandex.ru/\"} ")),
			requestPath: "/api/shorten",
			expectedData: expectedData{
				expectedContentType: "application/json",
				expectedStatusCode:  http.StatusConflict,
				expectedBody:        "{\"result\":\"http://localhost:8080/d41d8cd98f\"}",
				expectedLocation:    "",
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result, resultBody := testRequest(t, testServer, test.method, test.requestPath, test.clientID, test.requestBody)
			defer result.Body.Close()
			assert.Equal(t, test.expectedData.expectedStatusCode, result.StatusCode)
			assert.Equal(t, test.expectedData.expectedLocation, result.Header.Get("Location"))
			assert.Equal(t, test.expectedData.expectedContentType, result.Header.Get("Content-Type"))
			assert.Equal(t, test.expectedData.expectedBody, string(resultBody))
		})
	}
}
