package app

import (
	"crypto/md5"
	"fmt"
	"sync"
)

var URLMap = NewURL()

type URL struct {
	originalToShort map[string]string
	shortToOriginal map[string]string
	mutex           sync.Mutex
}

func NewURL() *URL {
	var originalToShortData = map[string]string{"https://practicum.yandex.ru/": "d41d8cd98f"}
	var shortToOriginalData = map[string]string{"d41d8cd98f": "https://practicum.yandex.ru/"}
	return &URL{originalToShort: originalToShortData, shortToOriginal: shortToOriginalData}
}

func encodeString(data string) string {
	encodedMD5 := md5.Sum([]byte(data))
	encodedMD5Trimed := encodedMD5[:5]
	return fmt.Sprintf("%x", encodedMD5Trimed)

}

func (url *URL) ToSortenURL(longURL string) (shortURL string) {
	url.mutex.Lock()
	defer url.mutex.Unlock()

	if value, ok := url.originalToShort[longURL]; ok {
		shortURL = value
		return
	}

	shortURL = encodeString(longURL)
	url.originalToShort[longURL] = shortURL
	url.shortToOriginal[shortURL] = longURL
	return
}

func (url *URL) ToOriginalURL(shortURL string) (longURL string) {
	url.mutex.Lock()
	defer url.mutex.Unlock()

	if value, ok := url.shortToOriginal[shortURL]; ok {
		longURL = value
		return
	}
	return ""
}
