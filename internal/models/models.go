package models

type Request struct {
	OriginalURL string `json:"url"`
}

type Response struct {
	ShortenURL string `json:"result"`
}

type URLPair struct {
	ShortenURL  string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLsClientID struct {
	URL      string
	ClientID int
}
