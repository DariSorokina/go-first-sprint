package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/DariSorokina/go-first-sprint/internal/cookie"
)

func exampleRequest(ts *httptest.Server, method, path string, requestBody io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, requestBody)
	if err != nil {
		log.Println(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	_, clientIDcookie := cookie.СreateCookieClientID("test")
	req.AddCookie(clientIDcookie)

	result, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer result.Body.Close()

	resultBody, err := io.ReadAll(result.Body)
	if err != nil {
		log.Println(err)
	}

	return result, string(resultBody)
}

func Example() {
	flagConfig, storage, serv := getTestServer()
	testServer := httptest.NewServer(serv.newRouter())
	defer testServer.Close()

	if flagConfig.FlagFileStoragePath != "" || flagConfig.FlagPostgresqlDSN != "" {
		defer storage.Close()
	}

	// example of using originalHandler
	httpMethod := http.MethodGet
	requestPath := "/d41d8cd98f"
	requestBody := bytes.NewBuffer([]byte(""))
	result, resultBody := exampleRequest(testServer, httpMethod, requestPath, requestBody)
	defer result.Body.Close()

	fmt.Println(result.StatusCode)
	fmt.Println(string(resultBody))

	// example of using shortenerHandler
	httpMethod = http.MethodPost
	requestPath = ""
	requestBody = bytes.NewBuffer([]byte("{\"url\":\"https://practicum.yandex.ru/\"} "))
	result, resultBody = exampleRequest(testServer, httpMethod, requestPath, requestBody)
	defer result.Body.Close()

	fmt.Println(result.StatusCode)
	fmt.Println(string(resultBody))

	// example of using shortenerHandlerJSON
	httpMethod = http.MethodPost
	requestPath = "/api/shorten"
	requestBody = bytes.NewBuffer([]byte("https://practicum.yandex.ru/"))
	result, resultBody = exampleRequest(testServer, httpMethod, requestPath, requestBody)
	defer result.Body.Close()

	fmt.Println(result.StatusCode)
	fmt.Println(string(resultBody))

	// data for shortenerBatchHandler
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

	// example of using shortenerBatchHandler
	httpMethod = http.MethodPost
	requestPath = "/api/shorten/batch"
	requestBody = bytes.NewBuffer([]byte(batchJSONData))
	result, resultBody = exampleRequest(testServer, httpMethod, requestPath, requestBody)
	defer result.Body.Close()

	fmt.Println(result.StatusCode)
	fmt.Println(string(resultBody))

}
