// Package models defines the data structures used for handling URL shortening requests and responses.
package models

// Request represents a structure for incoming requests containing the original URL to be shortened.
type Request struct {
	OriginalURL string `json:"url"`
}

// Response represents a structure for outgoing responses containing the shortened URL as a result.
type Response struct {
	ShortenURL string `json:"result"`
}

// URLPair represents a structure for storing the association between a shortened URL and its corresponding original URL.
type URLPair struct {
	ShortenURL  string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// URLsClientID represents a structure for storing multiple URLs associated with a specific client identified by a ClientID.
type URLsClientID struct {
	URLs     []string
	ClientID int
}
