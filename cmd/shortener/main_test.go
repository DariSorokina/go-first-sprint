package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, requestBody io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, requestBody)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	result, err := client.Do(req)
	require.NoError(t, err)
	defer result.Body.Close()

	resultBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	return result, string(resultBody)
}

func TestRouter(t *testing.T) {
	testServer := httptest.NewServer(LinkRouter())
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
		requestBody  io.Reader
		requestPath  string
		expectedData expectedData
	}{
		{
			name:        "handler: ShortenerHandler, test: StatusCreated",
			method:      http.MethodPost,
			requestBody: bytes.NewBuffer([]byte("https://practicum.yandex.ru/")),
			requestPath: "",
			expectedData: expectedData{
				expectedContentType: "text/plain",
				expectedStatusCode:  http.StatusCreated,
				expectedBody:        "http://localhost:8080/d41d8cd98f",
				expectedLocation:    "",
			},
		},
		{
			name:        "handler: ShortenerHandler, test: StatusBadRequest",
			method:      http.MethodPut,
			requestBody: bytes.NewBuffer([]byte("https://practicum.yandex.ru/")),
			requestPath: "",
			expectedData: expectedData{
				expectedContentType: "text/plain; charset=utf-8",
				expectedStatusCode:  http.StatusBadRequest,
				expectedBody:        "Only POST requests are allowed!\n",
				expectedLocation:    "",
			},
		},
		{
			name:        "handler: OriginalHandler, test: StatusTemporaryRedirect",
			method:      http.MethodGet,
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
			name:        "handler: OriginalHandler, test: StatusBadRequest",
			method:      http.MethodPut,
			requestBody: nil,
			requestPath: "/d41d8cd98f",
			expectedData: expectedData{
				expectedContentType: "text/plain; charset=utf-8",
				expectedStatusCode:  http.StatusBadRequest,
				expectedBody:        "Only GET requests are allowed!\n",
				expectedLocation:    "",
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result, resultBody := testRequest(t, testServer, test.method, test.requestPath, test.requestBody)
			defer result.Body.Close()
			assert.Equal(t, test.expectedData.expectedStatusCode, result.StatusCode)
			assert.Equal(t, test.expectedData.expectedLocation, result.Header.Get("Location"))
			assert.Equal(t, test.expectedData.expectedContentType, result.Header.Get("Content-Type"))
			assert.Equal(t, test.expectedData.expectedBody, string(resultBody))
		})
	}
}
