package app

import (
	"crypto/md5"
	"fmt"
)

type URL struct {
	Data map[string]string
}

func encodeString(data string) string {
	encodedMD5 := md5.Sum([]byte(data))
	encodedMD5Trimed := encodedMD5[:5]
	return fmt.Sprintf("%x", encodedMD5Trimed)

}

func (url *URL) ToSortenURL(longURL string) (shortURL string) {
	if url.Data == nil {
		url.Data = make(map[string]string)
	}

	if value, ok := url.Data[longURL]; ok {
		shortURL = value
		return
	}

	shortURL = encodeString(longURL)
	url.Data[longURL] = shortURL
	return
}

func (url *URL) ToOriginalURL(shortURL string) (longURL string) {
	for lURL, sURL := range url.Data {
		if sURL == shortURL {
			longURL = lURL
			return
		}
	}
	return ""
}

// func main() {
// 	var data = map[string]string{"https://practicum.yandex.ru/": "d41d8cd98f"}
// 	var myURL = URL{Data: data}
// 	fmt.Println(myURL)
// 	// myURL := URL{}
// 	res := myURL.ToSortenURL("https://practicum.yandex.ru/")
// 	fmt.Println(res)
// 	fmt.Println(myURL)

// 	res1 := myURL.ToOriginalURL("d41d8cd98f")
// 	fmt.Println(res1)
// 	fmt.Println(myURL)

// }
