package server

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		clientIDcookie := cookie.CreateCookieClientID(clientID)
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
	flagConfig := &config.FlagConfig{
		FlagRunAddr:         ":8080",
		FlagBaseURL:         "http://localhost:8080/",
		FlagLogLevel:        "info",
		FlagFileStoragePath: "/Users/dariasorokina/Desktop/yp_golang/go-first-sprint/internal/storage/short-url-db.json",
		FlagPostgresqlDSN:   "host=localhost user=app password=123qwe dbname=urls_database sslmode=disable"}

	// flagConfig := config.ParseFlags()
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

	// data for shortenerBatchHandler test
	batchData := []struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}{
		{
			CorrelationID: "qwerty",
			OriginalURL:   "https://practicum.yandex.ru/",
		},
	}
	batchJSONData, err := json.Marshal(batchData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

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
				expectedStatusCode:  http.StatusCreated,
				expectedBody:        "http://localhost:8080/0dd198178d",
				expectedLocation:    "",
			},
		},
		{
			name:        "handler: OriginalHandler, test: StatusTemporaryRedirect",
			method:      http.MethodGet,
			clientID:    1,
			requestBody: nil,
			requestPath: "/0dd198178d",
			expectedData: expectedData{
				expectedContentType: "",
				expectedStatusCode:  http.StatusTemporaryRedirect,
				expectedBody:        "",
				expectedLocation:    "https://practicum.yandex.ru/",
			},
		},
		{
			name:        "handler: shortenerHandlerJSON, test: StatusConflict",
			method:      http.MethodPost,
			clientID:    1,
			requestBody: bytes.NewBuffer([]byte("{\"url\":\"https://practicum.yandex.ru/\"} ")),
			requestPath: "/api/shorten",
			expectedData: expectedData{
				expectedContentType: "application/json",
				expectedStatusCode:  http.StatusConflict,
				expectedBody:        "{\"result\":\"http://localhost:8080/0dd198178d\"}",
				expectedLocation:    "",
			},
		},
		{
			name:        "handler: shortenerBatchHandler, test: StatusCreated",
			method:      http.MethodPost,
			clientID:    1,
			requestBody: bytes.NewBuffer([]byte(batchJSONData)),
			requestPath: "/api/shorten/batch",
			expectedData: expectedData{
				expectedContentType: "application/json",
				expectedStatusCode:  http.StatusCreated,
				expectedBody:        "[{\"correlation_id\":\"qwerty\",\"short_url\":\"http://localhost:8080/0dd198178d\"}]",
				expectedLocation:    "",
			},
		},
		{
			name:        "handler: urlsByIDHandler, test: StatusCreated",
			method:      http.MethodGet,
			clientID:    1,
			requestBody: nil,
			requestPath: "/api/user/urls",
			expectedData: expectedData{
				expectedContentType: "application/json",
				expectedStatusCode:  http.StatusOK,
				expectedBody:        "[{\"short_url\":\"http://localhost:8080/0dd198178d\",\"original_url\":\"https://practicum.yandex.ru/\"}]",
				expectedLocation:    "",
			},
		},
		{
			name:        "handler: deleteURLsHandler, test: StatusCreated",
			method:      http.MethodDelete,
			clientID:    1,
			requestBody: bytes.NewBuffer([]byte("[\"0dd198178d\"]")),
			requestPath: "/api/user/urls",
			expectedData: expectedData{
				expectedContentType: "",
				expectedStatusCode:  http.StatusAccepted,
				expectedBody:        "",
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

func getTestServer() (flagConfig *config.FlagConfig, storageFile storage.Database, serv *Server) {
	flagConfig = &config.FlagConfig{
		FlagRunAddr:         ":8080",
		FlagBaseURL:         "http://localhost:8080/",
		FlagLogLevel:        "info",
		FlagFileStoragePath: "/Users/dariasorokina/Desktop/yp_golang/go-first-sprint/internal/storage/short-url-db.json"}

	var l *logger.Logger
	var err error
	if l, err = logger.CreateLogger(flagConfig.FlagLogLevel); err != nil {
		log.Fatal("Failed to create logger:", err)
	}

	storageFile, err = storage.SetStorage(flagConfig, l)
	if err != nil {
		panic(err)
	}

	app := app.NewApp(storageFile, l)
	serv = NewServer(app, flagConfig, l)

	return flagConfig, storageFile, serv
}

func BenchmarkHandlers(b *testing.B) {
	b.StopTimer()
	flagConfig, storage, serv := getTestServer()
	testServer := httptest.NewServer(serv.newRouter())
	defer testServer.Close()

	if flagConfig.FlagFileStoragePath != "" || flagConfig.FlagPostgresqlDSN != "" {
		defer storage.Close()
	}

	// data for shortenerBatchHandler test
	batchData := []struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}{
		{
			CorrelationID: "qwerty",
			OriginalURL:   "https://practicum.yandex.ru/",
		},
	}
	batchJSONData, err := json.Marshal(batchData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	b.Run("BenchmarkOriginalHandlerRequest", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			httpMethod := http.MethodGet
			requestPath := "/d41d8cd98f"
			requestBody := bytes.NewBuffer([]byte(""))
			httpRequest := httptest.NewRequest(httpMethod, testServer.URL+requestPath, requestBody)
			responseRecorder := httptest.NewRecorder()
			httpRequest.Header.Set("ClientID", "1")
			b.StartTimer()
			serv.handlers.originalHandler(responseRecorder, httpRequest)
		}
	})

	b.Run("BenchmarkShortenerHandlerRequest", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			httpMethod := http.MethodPost
			requestPath := ""
			requestBody := bytes.NewBuffer([]byte("{\"url\":\"https://practicum.yandex.ru/\"} "))
			httpRequest := httptest.NewRequest(httpMethod, testServer.URL+requestPath, requestBody)
			responseRecorder := httptest.NewRecorder()
			httpRequest.Header.Set("ClientID", "1")
			b.StartTimer()
			serv.handlers.shortenerHandler(responseRecorder, httpRequest)
		}
	})

	b.Run("BenchmarkShortenerHandlerJSONRequest", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			httpMethod := http.MethodPost
			requestPath := "/api/shorten"
			requestBody := bytes.NewBuffer([]byte("https://practicum.yandex.ru/"))
			httpRequest := httptest.NewRequest(httpMethod, testServer.URL+requestPath, requestBody)
			responseRecorder := httptest.NewRecorder()
			httpRequest.Header.Set("ClientID", "1")
			b.StartTimer()
			serv.handlers.shortenerHandlerJSON(responseRecorder, httpRequest)
		}
	})

	b.Run("BenchmarkShortenerBatchHandler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			httpMethod := http.MethodPost
			requestPath := "/api/shorten/batch"
			requestBody := bytes.NewBuffer([]byte(batchJSONData))
			httpRequest := httptest.NewRequest(httpMethod, testServer.URL+requestPath, requestBody)
			responseRecorder := httptest.NewRecorder()
			httpRequest.Header.Set("ClientID", "1")
			b.StartTimer()
			serv.handlers.shortenerBatchHandler(responseRecorder, httpRequest)
		}
	})

	b.Run("BenchmarkUrlsByIDHandler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			httpMethod := http.MethodGet
			requestPath := "/api/user/urls"
			requestBody := bytes.NewBuffer([]byte(""))
			httpRequest := httptest.NewRequest(httpMethod, testServer.URL+requestPath, requestBody)
			responseRecorder := httptest.NewRecorder()
			httpRequest.Header.Set("ClientID", "1")
			b.StartTimer()
			serv.handlers.urlsByIDHandler(responseRecorder, httpRequest)
		}
	})

	b.Run("BenchmarkDeleteURLsHandler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			httpMethod := http.MethodDelete
			requestPath := "/api/user/urls"
			requestBody := bytes.NewBuffer([]byte("[\"0dd198178d\"]"))
			httpRequest := httptest.NewRequest(httpMethod, testServer.URL+requestPath, requestBody)
			responseRecorder := httptest.NewRecorder()
			httpRequest.Header.Set("ClientID", "1")
			b.StartTimer()
			serv.handlers.deleteURLsHandler(responseRecorder, httpRequest)
		}
	})
}
