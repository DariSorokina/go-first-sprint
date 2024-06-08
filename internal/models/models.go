package models

type Request struct {
	OriginalURL string `json:"url"`
}

type Response struct {
	ShortenURL string `json:"result"`
}
