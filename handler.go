package todo

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Here is the about page.")
}

func LoggerHandler(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stderr, next)
}

func RecoverHandler(next http.Handler) http.Handler {
	var fn = func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %s", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// func AuthHandler(next http.Handler) http.Handler {
// var fn = func(w http.ResponseWriter, r *http.Request) {
// var authToken = r.Header().Get("Authorization")
// var user, err := getUser(authToken)
//
// if err != nil {
// http.Error(w, http.StatusText(401), 401)
// return
// }
//
// context.Set(r, "user", user)
// next.ServeHTTP(w, r)
// }
//
// return http.HandlerFunc(fn)
// }
