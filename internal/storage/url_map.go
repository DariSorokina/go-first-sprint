package storage

import "sync"

type URL struct {
	OriginalToShort map[string]string
	ShortToOriginal map[string]string
	Mutex           sync.RWMutex
}

func NewURL() *URL {
	return &URL{
		OriginalToShort: map[string]string{"https://practicum.yandex.ru/": "d41d8cd98f"},
		ShortToOriginal: map[string]string{"d41d8cd98f": "https://practicum.yandex.ru/"},
	}
}
