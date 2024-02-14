package storage

import (
	"log"
	"sync"
)

type Storage struct {
	fileStorage     *fileStorage
	originalToShort map[string]string
	shortToOriginal map[string]string
	mutex           sync.RWMutex
}

func NewStorage(fileName string) *Storage {
	fileStorage := newFileStorage(fileName)
	if fileName != "" {
		var url = []*fileLine{
			{
				ShortURL:    "d41d8cd98f",
				OriginalURL: "https://practicum.yandex.ru/",
			},
		}

		readURLs, err := fileStorage.consumer.readURLs()
		if err != nil {
			log.Fatal(err)
		}
		obtainedUrls := append(url, readURLs...)

		originalToShort := make(map[string]string)
		shortToOriginal := make(map[string]string)
		originalToShort, shortToOriginal = addURLsToMap(obtainedUrls, originalToShort, shortToOriginal)

		return &Storage{
			fileStorage:     fileStorage,
			originalToShort: originalToShort,
			shortToOriginal: shortToOriginal,
		}
	}

	return &Storage{
		originalToShort: map[string]string{"https://practicum.yandex.ru/": "d41d8cd98f"},
		shortToOriginal: map[string]string{"d41d8cd98f": "https://practicum.yandex.ru/"},
	}
}

func (storage *Storage) SetValue(shortURL, longURL string) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	var url = []*fileLine{
		{
			ShortURL:    shortURL,
			OriginalURL: longURL,
		},
	}

	if storage.fileStorage.fileName != "" {
		err := storage.fileStorage.producer.writeURL(url[0])
		if err != nil {
			log.Fatal(err)
		}
	}

	storage.originalToShort, storage.shortToOriginal = addURLsToMap(url, storage.originalToShort, storage.shortToOriginal)
}

func (storage *Storage) GetShort(longURL string) (shortURL string) {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	if value, ok := storage.originalToShort[longURL]; ok {
		shortURL = value
		return shortURL
	}
	return ""
}

func (storage *Storage) GetOriginal(shortURL string) (longURL string) {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	if value, ok := storage.shortToOriginal[shortURL]; ok {
		longURL = value
		return
	}
	return ""
}

func (storage *Storage) CloseFile() {
	storage.fileStorage.producer.close()
	storage.fileStorage.consumer.close()
}
