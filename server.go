package todo

import (
	"net/http"

	"github.com/justinas/alice"
)

// NewAppHandler creates a new http.ServeMux,
// and registers the application handlers.
func NewAppHandler(store Store) http.Handler {
	var router = http.NewServeMux()
	var chain = alice.New(LoggingHandler, RecoverHandler)

	// todos api
	var todoRouter = NewRouter()
	var todoContext = NewContext(store)
	todoContext.Register(todoRouter)

	router.Handle("/api/", chain.Then(todoRouter))

	// static pages
	router.Handle("/index.html", chain.Then(HomePage(store)))
	router.Handle("/about", chain.ThenFunc(AboutPage))
	router.Handle("/", chain.ThenFunc(StaticPages))

	return router
}
