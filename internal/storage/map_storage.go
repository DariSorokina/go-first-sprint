// Package storage provides primitives for connecting to data storages.
package storage

import (
	"context"
	"log"
	"sync"

	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/models"
)

// Storage represents a storage structure for managing file storage, mappings between original and short URLs,
// synchronization with a mutex, and logging functionality.
type Storage struct {
	fileStorage     *fileStorage      // File storage instance.
	originalToShort map[string]string // Mapping of original URLs to short URLs.
	shortToOriginal map[string]string // Mapping of short URLs to original URLs.
	mutex           sync.RWMutex      // Mutex for synchronization.
	log             *logger.Logger    // Logger for recording events and errors.
}

// NewStorage creates a new Storage instance with the provided file name and logger.
func NewStorage(fileName string, l *logger.Logger) *Storage {
	fileStorage := newFileStorage(fileName, l)
	if fileName != "" {
		var url = []*fileLine{
			{
				ShortURL:    "d41d8cd98f",
				OriginalURL: "https://practicum.yandex.ru/",
			},
		}

		readURLs, err := fileStorage.consumer.readURLs()
		if err != nil {
			log.Println(err)
		}
		obtainedUrls := append(url, readURLs...)

		originalToShort := make(map[string]string)
		shortToOriginal := make(map[string]string)
		originalToShort, shortToOriginal = addURLsToMap(obtainedUrls, originalToShort, shortToOriginal)

		return &Storage{
			fileStorage:     fileStorage,
			originalToShort: originalToShort,
			shortToOriginal: shortToOriginal,
			log:             l,
		}
	}
	// Default mappings if fileName is empty.
	return &Storage{
		originalToShort: map[string]string{"https://practicum.yandex.ru/": "d41d8cd98f"},
		shortToOriginal: map[string]string{"d41d8cd98f": "https://practicum.yandex.ru/"},
		log:             l,
	}
}

// SetValue sets a value in the map storage for a given shortURL, longURL, and userID.
func (storage *Storage) SetValue(ctx context.Context, shortURL, longURL string, userID int) {
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
			log.Println(err)
		}
	}

	storage.originalToShort, storage.shortToOriginal = addURLsToMap(url, storage.originalToShort, storage.shortToOriginal)
}

// GetShort retrieves the short URL corresponding to a given long URL from the map storage.
func (storage *Storage) GetShort(ctx context.Context, longURL string) (shortURL string, err error) {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	if value, ok := storage.originalToShort[longURL]; ok {
		shortURL = value
		return shortURL, ErrShortURLAlreadyExist
	}
	return "", nil
}

// GetOriginal retrieves the original long URL corresponding to a given short URL from the map storage.
func (storage *Storage) GetOriginal(ctx context.Context, shortURL string) (longURL string, getOriginalErr error) {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	if value, ok := storage.shortToOriginal[shortURL]; ok {
		longURL = value
		return
	}
	return "", nil
}

// GetURLsByUserID retrieves URLs associated with a given user ID from the mp storage.
func (storage *Storage) GetURLsByUserID(ctx context.Context, userID int) (urls []models.URLPair) {
	return
}

func (storage *Storage) DeleteURLsWorker(shortURLs []string, userID int) {
}

func (storage *Storage) Ping(ctx context.Context) error {
	return nil
}

// Close closes the file storage producer and consumer.
func (storage *Storage) Close() {
	if storage.fileStorage.producer != nil {
		storage.fileStorage.producer.close()
	}
	if storage.fileStorage.consumer != nil {
		storage.fileStorage.consumer.close()
	}
}
