package storage

import (
	"encoding/json"
	"log"
	"os"
)

type FileStorage struct {
	producer *Producer
	consumer *Consumer
	fileName string
}

func NewFileStorage(fileName string) *FileStorage {
	producer, err := NewProducer(fileName)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := NewConsumer(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return &FileStorage{
		producer: producer,
		consumer: consumer,
	}
}

type fileLine struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteURL(url *fileLine) error {
	return p.encoder.Encode(&url)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadURLs() ([]*fileLine, error) {

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

func (c *Consumer) Close() error {
	return c.file.Close()
}

func AddURLsToMap(urls []*fileLine, originalToShort, shortToOriginal map[string]string) (map[string]string, map[string]string) {
	for _, url := range urls {
		originalToShort[url.OriginalURL] = url.ShortURL
		shortToOriginal[url.ShortURL] = url.OriginalURL
	}
	return originalToShort, shortToOriginal
}
