package app

import (
	"crypto/md5"
	"fmt"
)

type URL struct {
	data map[string]string
}

func encodeString(data string) string {
	encodedMD5 := md5.Sum([]byte(data))
	encodedMD5Trimed := encodedMD5[:5]
	return fmt.Sprintf("%x", encodedMD5Trimed)

}

func (url *URL) ToSortenURL(longURL string) (shortURL string) {
	if url.data == nil {
		url.data = make(map[string]string)
	}

	if value, ok := url.data[longURL]; ok {
		shortURL = value
		return
	}

	shortURL = encodeString(longURL)
	url.data[longURL] = shortURL
	return
}

func (url *URL) ToOriginalURL(shortURL string) (longURL string) {
	for lURL, sURL := range url.data {
		if sURL == shortURL {
			longURL = lURL
			return
		}
	}
	return ""
}

// func main() {
// 	myURL := URL{}
// 	res := myURL.ToSortenURL("https://practicum.yandex.ru")
// 	fmt.Println(res)

// 	res1 := myURL.ToOriginalURL("6bdb5b0e26")
// 	fmt.Println(res1)

// }
