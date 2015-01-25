package todo

import (
	"net/http"

	"github.com/justinas/alice"
)

// NewAppHandler creates a new http.ServeMux,
// and registers the application handlers.
func NewAppHandler(store Store) http.Handler {
	var chain = alice.New(LoggingHandler, RecoverHandler)

	var router = http.NewServeMux()
	router.Handle("/", chain.ThenFunc(Index))
	router.Handle("/about", chain.ThenFunc(About))

	var todoRouter = NewRouter()
	NewContext(store).Register(todoRouter)
	router.Handle("/api/", chain.Then(todoRouter))

	return router
}
