package storage

import (
	"encoding/json"
	"log"
	"os"
)

type fileStorage struct {
	producer *producer
	consumer *consumer
	fileName string
}

func newFileStorage(fileName string) *fileStorage {
	producer, err := newProducer(fileName)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := newConsumer(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return &fileStorage{
		producer: producer,
		consumer: consumer,
		fileName: fileName,
	}
}

type fileLine struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func newProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) writeURL(url *fileLine) error {
	return p.encoder.Encode(&url)
}

func (p *producer) close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func newConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) readURLs() ([]*fileLine, error) {

	var urls []*fileLine
	for c.decoder.More() {
		url := &fileLine{}
		if err := c.decoder.Decode(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (c *consumer) close() error {
	return c.file.Close()
}

func addURLsToMap(urls []*fileLine, originalToShort, shortToOriginal map[string]string) (map[string]string, map[string]string) {
	for _, url := range urls {
		originalToShort[url.OriginalURL] = url.ShortURL
		shortToOriginal[url.ShortURL] = url.OriginalURL
	}
	return originalToShort, shortToOriginal
}
