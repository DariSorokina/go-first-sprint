package storage

import "sync"

type Storage struct {
	originalToShort map[string]string
	shortToOriginal map[string]string
	mutex           sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		originalToShort: map[string]string{"https://practicum.yandex.ru/": "d41d8cd98f"},
		shortToOriginal: map[string]string{"d41d8cd98f": "https://practicum.yandex.ru/"},
	}
}

func (storage *Storage) SetValue(shortURL, longURL string) {
	storage.originalToShort[longURL] = shortURL
	storage.shortToOriginal[shortURL] = longURL
}

func (storage *Storage) GetShort(longURL string) (shortURL string) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if value, ok := storage.originalToShort[longURL]; ok {
		shortURL = value
		return shortURL
	}
	return ""
}

func (storage *Storage) GetOriginal(shortURL string) (longURL string) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if value, ok := storage.shortToOriginal[shortURL]; ok {
		longURL = value
		return
	}
	return ""
}
