package middlewares

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
)

func Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no token"))
			return
		}
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			log.Fatalln(err)
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			user := claims["sub"]
			err := r.ParseForm()
			if err != nil {
				log.Fatalln("Form parsing error")
				return
			}
			r.Form.Add("user", user.(string))
		} else {
			log.Println("invalid token")
			return
		}
		next.ServeHTTP(w, r)
	})
}
