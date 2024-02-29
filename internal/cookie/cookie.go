package cookie

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TOKENEXP = time.Hour * 3
const SECRETKEY = "supersecretkey"

var generatedUsersIDs = []int{1}

func generateUserID() int {
	randomNumber := rand.Intn(1000001)
	return randomNumber
}

func createJWTString() (generatedUserID int, tokenString string, err error) {
	generatedUserID = generateUserID()
	generatedUsersIDs = append(generatedUsersIDs, generatedUserID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKENEXP)),
		},
		UserID: generatedUserID,
	})

	tokenString, err = token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return generatedUserID, "", err
	}

	return generatedUserID, tokenString, nil
}

func getUserID(tokenString string) int {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRETKEY), nil
		})

	if err != nil {
		return -1
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return -1
	}

	fmt.Println("Token is valid")
	return claims.UserID
}

func createCookieClientID() (generatedUserID int, cookie *http.Cookie) {
	generatedUserID, JWTString, err := createJWTString()
	if err != nil {
		log.Println(err)
	}
	cookie = &http.Cookie{
		Name:     "ClientID",
		Value:    JWTString,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	}

	return generatedUserID, cookie
}

func validateUserID(userID int) bool {
	for _, id := range generatedUsersIDs {
		if userID == id {
			return true
		}
	}
	return false
}

func CookieMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			reseivedCookie, err := r.Cookie("ClientID")
			fmt.Println(err)
			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					userID, createdCookie := createCookieClientID()
					http.SetCookie(w, createdCookie)
					userIDString := strconv.Itoa(userID)
					r.Header.Set("ClientID", userIDString)
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusUnauthorized)
					}
					h.ServeHTTP(w, r)
				default:
					log.Println(err)
					http.Error(w, "server error", http.StatusInternalServerError)
				}
				return
			}

			clientID := reseivedCookie.Value
			if clientID == "" {
				fmt.Println("clientID == ")
			}
			if clientID == "" {
				_, createdCookie := createCookieClientID()
				http.SetCookie(w, createdCookie)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID := getUserID(clientID)
			validUserID := validateUserID(userID)

			userIDString := strconv.Itoa(userID)
			fmt.Println(userIDString)

			if validUserID {
				r.Header.Set("ClientID", userIDString)
				h.ServeHTTP(w, r)
			} else {
				_, createdCookie := createCookieClientID()
				http.SetCookie(w, createdCookie)
				h.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}
