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

// Claims struct includes the registered claims from jwt package and a custom UserID field.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// Constants for token expiration time and the secret key used for signing JWT.
const TOKENEXP = time.Hour * 3
const SECRETKEY = "supersecretkey"

var generatedUsersIDs = []int{1}

func generateUserID() int {
	randomNumber := rand.Intn(1000001)
	return randomNumber
}

func createJWTString(task string) (generatedUserID int, tokenString string, err error) {
	if task != "test" {
		generatedUserID = generateUserID()
		generatedUsersIDs = append(generatedUsersIDs, generatedUserID)
	} else {
		generatedUserID = 1
	}

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

// СreateCookieClientID creates a JWT token and embeds it in an HTTP cookie.
// Returns the generated user ID and the created cookie.
func СreateCookieClientID(task string) (generatedUserID int, cookie *http.Cookie) {
	generatedUserID, JWTString, err := createJWTString(task)
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

// CookieMiddleware returns a middleware that ensures each request has a valid JWT in the "ClientID" cookie.
// If the cookie is missing or invalid, a new JWT is created and set as a cookie.
func CookieMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			reseivedCookie, err := r.Cookie("ClientID")
			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					userID, createdCookie := СreateCookieClientID("")
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
				_, createdCookie := СreateCookieClientID("")
				http.SetCookie(w, createdCookie)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID := getUserID(clientID)
			validUserID := validateUserID(userID)

			userIDString := strconv.Itoa(userID)

			if validUserID {
				r.Header.Set("ClientID", userIDString)
				h.ServeHTTP(w, r)
			} else {
				_, createdCookie := СreateCookieClientID("")
				http.SetCookie(w, createdCookie)
				h.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}
