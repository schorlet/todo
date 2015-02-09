package todo

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/handlers"
)

var templates *template.Template

func init() {
	templates = template.Must(template.New("base").Funcs(
		template.FuncMap{"jsstr": templateJSStr}).ParseFiles("static/index.html"))
}

// AboutPage handles about page.
func AboutPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Here is the about page.")
}

// HomePage handles index.html page.
func HomePage(store Store) http.Handler {
	var fn = func(w http.ResponseWriter, r *http.Request) {
		var todos = store.List()
		var err = templates.ExecuteTemplate(w, "index.html", todos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	return http.HandlerFunc(fn)
}

// StaticPages handles static pages.
func StaticPages(w http.ResponseWriter, r *http.Request) {
	var upath = r.URL.Path

	if upath == "/" {
		http.Redirect(w, r, "index.html", http.StatusMovedPermanently)
		return
	}

	upath = "static/" + upath
	http.ServeFile(w, r, path.Clean(upath))
}

// LoggingHandler
func LoggingHandler(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stderr, next)
}

// RecoverHandler
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
