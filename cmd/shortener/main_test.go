package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ShortenerHandler(t *testing.T) {
	type expectedData struct {
		expectedContentType string
		expectedStatusCode  int
		expectedBody        string
	}

	testCases := []struct {
		name         string
		method       string
		request      string
		expectedData expectedData
	}{
		{
			name:    "StatusCreated",
			method:  http.MethodPost,
			request: "https://practicum.yandex.ru/",
			expectedData: expectedData{
				expectedContentType: "text/plain",
				expectedStatusCode:  http.StatusCreated,
				expectedBody:        "http://localhost:8080/d41d8cd98f",
			},
		},
		{
			name:    "StatusBadRequest",
			method:  http.MethodPut,
			request: "https://practicum.yandex.ru/",
			expectedData: expectedData{
				expectedContentType: "text/plain; charset=utf-8",
				expectedStatusCode:  http.StatusBadRequest,
				expectedBody:        "Only POST requests are allowed!\n",
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			w := httptest.NewRecorder()
			ShortenerHandler(w, request)

			result := w.Result()
			assert.Equal(t, test.expectedData.expectedStatusCode, result.StatusCode)
			assert.Equal(t, test.expectedData.expectedContentType, result.Header.Get("Content-Type"))

			defer result.Body.Close()
			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			assert.Equal(t, test.expectedData.expectedBody, string(resBody))
		})
	}
}

func Test_OriginalHandler(t *testing.T) {
	type expectedData struct {
		expectedContentType string
		expectedStatusCode  int
		expectedBody        string
		expectedLocation    string
	}

	testCases := []struct {
		name         string
		method       string
		request      string
		expectedData expectedData
	}{
		{
			name:    "StatusTemporaryRedirect",
			method:  http.MethodGet,
			request: "/d41d8cd98f",
			expectedData: expectedData{
				expectedContentType: "",
				expectedStatusCode:  http.StatusTemporaryRedirect,
				expectedBody:        "",
				expectedLocation:    "https://practicum.yandex.ru/",
			},
		},
		{
			name:    "StatusBadRequest",
			method:  http.MethodPut,
			request: "/d41d8cd98f",
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
			router := mux.NewRouter()
			router.HandleFunc("/{id}", OriginalHandler).Methods(test.method)
			request := httptest.NewRequest(test.method, test.request, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			result := w.Result()
			fmt.Println(result)
			assert.Equal(t, test.expectedData.expectedStatusCode, result.StatusCode)
			assert.Equal(t, test.expectedData.expectedLocation, result.Header.Get("Location"))
			assert.Equal(t, test.expectedData.expectedContentType, result.Header.Get("Content-Type"))

			defer result.Body.Close()
			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			assert.Equal(t, test.expectedData.expectedBody, string(resBody))
		})
	}
}
