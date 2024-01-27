package server

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
	"github.com/stretchr/testify/require"
)

// func testRequest(t *testing.T, ts *httptest.Server, method, path string, requestBody io.Reader) (*http.Response, string) {
// 	req, err := http.NewRequest(method, ts.URL+path, requestBody)
// 	require.NoError(t, err)

// 	client := &http.Client{
// 		CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 			return http.ErrUseLastResponse
// 		},
// 	}

// 	result, err := client.Do(req)
// 	require.NoError(t, err)
// 	defer result.Body.Close()

// 	resultBody, err := io.ReadAll(result.Body)
// 	require.NoError(t, err)

// 	return result, string(resultBody)
// }

// func TestRouter(t *testing.T) {
// 	flagConfig := config.ParseFlags()
// 	var l *logger.Logger
// 	var err error
// 	if l, err = logger.CreateLogger(flagConfig.FlagLogLevel); err != nil {
// 		log.Fatal("Failed to create logger:", err)
// 	}

// 	storage := storage.NewStorage()
// 	app := app.NewApp(storage)
// 	serv := NewServer(app, flagConfig, l)
// 	testServer := httptest.NewServer(serv.newRouter())
// 	defer testServer.Close()

// 	type expectedData struct {
// 		expectedContentType string
// 		expectedStatusCode  int
// 		expectedBody        string
// 		expectedLocation    string
// 	}

// 	testCases := []struct {
// 		name         string
// 		method       string
// 		requestBody  io.Reader
// 		requestPath  string
// 		expectedData expectedData
// 	}{
// 		{
// 			name:        "handler: ShortenerHandler, test: StatusCreated",
// 			method:      http.MethodPost,
// 			requestBody: bytes.NewBuffer([]byte("https://practicum.yandex.ru/")),
// 			requestPath: "",
// 			expectedData: expectedData{
// 				expectedContentType: "text/plain",
// 				expectedStatusCode:  http.StatusCreated,
// 				expectedBody:        "http://localhost:8080/d41d8cd98f",
// 				expectedLocation:    "",
// 			},
// 		},
// 		{
// 			name:        "handler: OriginalHandler, test: StatusTemporaryRedirect",
// 			method:      http.MethodGet,
// 			requestBody: nil,
// 			requestPath: "/d41d8cd98f",
// 			expectedData: expectedData{
// 				expectedContentType: "",
// 				expectedStatusCode:  http.StatusTemporaryRedirect,
// 				expectedBody:        "",
// 				expectedLocation:    "https://practicum.yandex.ru/",
// 			},
// 		},
// 		{
// 			name:        "handler: shortenerHandlerJSON, test: StatusCreated",
// 			method:      http.MethodPost,
// 			requestBody: bytes.NewBuffer([]byte("{\"url\":\"https://practicum.yandex.ru/\"} ")),
// 			requestPath: "/api/shorten",
// 			expectedData: expectedData{
// 				expectedContentType: "application/json",
// 				expectedStatusCode:  http.StatusCreated,
// 				expectedBody:        "{\"result\":\"http://localhost:8080/d41d8cd98f\"}",
// 				expectedLocation:    "",
// 			},
// 		},
// 	}
// 	for _, test := range testCases {
// 		t.Run(test.name, func(t *testing.T) {
// 			result, resultBody := testRequest(t, testServer, test.method, test.requestPath, test.requestBody)
// 			defer result.Body.Close()
// 			assert.Equal(t, test.expectedData.expectedStatusCode, result.StatusCode)
// 			assert.Equal(t, test.expectedData.expectedLocation, result.Header.Get("Location"))
// 			assert.Equal(t, test.expectedData.expectedContentType, result.Header.Get("Content-Type"))
// 			assert.Equal(t, test.expectedData.expectedBody, string(resultBody))
// 		})
// 	}
// }

func TestGzipCompression(t *testing.T) {
	flagConfig := config.ParseFlags()
	var l *logger.Logger
	var err error
	if l, err = logger.CreateLogger(flagConfig.FlagLogLevel); err != nil {
		log.Fatal("Failed to create logger:", err)
	}

	storage := storage.NewStorage()
	app := app.NewApp(storage)
	serv := NewServer(app, flagConfig, l)
	testServer := httptest.NewServer(serv.newRouter())
	defer testServer.Close()

	requestBody := `{"url": "https://practicum.yandex.ru"}`
	successBody := `{"result":"http://localhost:8080/6bdb5b0e26"}`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", testServer.URL, buf)
		r.RequestURI = "/api/shorten"
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", testServer.URL, buf)
		r.RequestURI = "/api/shorten"
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.JSONEq(t, successBody, string(b))
	})
}