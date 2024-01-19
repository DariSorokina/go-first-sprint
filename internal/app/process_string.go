package app

import (
	"crypto/md5"
	"fmt"

	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
)

func encodeString(data string) string {
	encodedMD5 := md5.Sum([]byte(data))
	encodedMD5Trimed := encodedMD5[:5]
	return fmt.Sprintf("%x", encodedMD5Trimed)

}

func ToSortenURL(url *storage.URL, longURL string) (shortURL string) {
	url.Mutex.Lock()
	defer url.Mutex.Unlock()

	if value, ok := url.OriginalToShort[longURL]; ok {
		shortURL = value
		return
	}

	shortURL = encodeString(longURL)
	url.OriginalToShort[longURL] = shortURL
	url.ShortToOriginal[shortURL] = longURL
	return
}

func ToOriginalURL(url *storage.URL, shortURL string) (longURL string) {
	url.Mutex.Lock()
	defer url.Mutex.Unlock()

	if value, ok := url.ShortToOriginal[shortURL]; ok {
		longURL = value
		return
	}
	return ""
}
