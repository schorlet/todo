package todo

import "github.com/gorilla/mux"

const (
	RouteList   = "Todo.List"
	RouteCreate = "Todo.Create"

	// _/{id}
	RouteFind   = "Todo.Find"
	RouteUpdate = "Todo.Update"
	RouteDelete = "Todo.Delete"

	// _/{status}
	RouteFilter = "Todo.Filter"
	RouteClear  = "Todo.Clear"
)

// NewRouter creates a new mux.Router and defines HTTP methods
// with URL paths starting with "/api/todos".
func NewRouter() *mux.Router {
	return NewRouterPrefix("/api/todos")
}

// NewRouterPrefix creates a new mux.Router and defines HTTP methods
// with URL paths starting with the specified prefix.
// Basically, the prefix would be "/api/todos" without trailing slash.
func NewRouterPrefix(prefix string) *mux.Router {
	var router = mux.NewRouter()

	router.Methods("GET").Path(prefix).Name(RouteList)
	router.Methods("POST").Path(prefix).Name(RouteCreate)

	router.Methods("GET").Path(prefix + "/{id:[0-9]+}").Name(RouteFind)
	router.Methods("PUT").Path(prefix + "/{id:[0-9]+}").Name(RouteUpdate)
	router.Methods("DELETE").Path(prefix + "/{id:[0-9]+}").Name(RouteDelete)

	router.Methods("GET").Path(prefix + "/{status:[a-z]+}").Name(RouteFilter)
	router.Methods("DELETE").Path(prefix + "/{status:[a-z]+}").Name(RouteClear)

	return router
}
