package todo

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Context manages todos.
type Context struct {
	Store
}

// NewContext returns a new context.
func NewContext(store Store) Context {
	return Context{store}
}

// Register sets a handlers to the routes.
func (ctx Context) Register(router *mux.Router) {
	router.Get(RouteList).Handler(ErrorFunc(ctx.List))
	router.Get(RouteCreate).Handler(ErrorFunc(ctx.Create))

	// _/{id}
	router.Get(RouteFind).Handler(ErrorFunc(ctx.Find))
	router.Get(RouteUpdate).Handler(ErrorFunc(ctx.Update))
	router.Get(RouteDelete).Handler(ErrorFunc(ctx.Delete))

	// _/{status}
	router.Get(RouteFilter).Handler(ErrorFunc(ctx.Filter))
	router.Get(RouteClear).Handler(ErrorFunc(ctx.Clear))
}

// List handles todos listing.
func (ctx Context) List(w http.ResponseWriter, r *http.Request) error {
	var todos = ctx.Store.List()
	return writeJSON(w, todos, http.StatusOK) // 200
}

// Filter handles todos filtering.
func (ctx Context) Filter(w http.ResponseWriter, r *http.Request) error {
	var status = readStatus(w, r)
	var todos = ctx.Store.Filter(status)
	return writeJSON(w, todos, http.StatusOK) // 200
}

// Clear handles todos deletion.
func (ctx Context) Clear(w http.ResponseWriter, r *http.Request) error {
	var status = readStatus(w, r)
	var count, err = ctx.Store.Clear(status)
	if err != nil {
		return err // 500
	}
	var cleared = map[string]int64{"count": count}
	return writeJSON(w, cleared, http.StatusOK) // 200
}

// Create handles todo creation.
func (ctx Context) Create(w http.ResponseWriter, r *http.Request) error {
	var todo, err = readTodo(w, r)
	if err != nil {
		return BadRequest{err} // 400
	}

	todo.ID = ""

	err = ctx.Store.Save(todo)
	if err != nil {
		return err // 500
	}

	return writeJSON(w, todo, http.StatusCreated) // 201
}

// Find handles todo selection by ID.
func (ctx Context) Find(w http.ResponseWriter, r *http.Request) error {
	var id = readID(w, r)

	var todo, err = ctx.Store.Find(id)
	if err != nil {
		return err // 500
	}

	return writeJSON(w, todo, http.StatusOK) // 200
}

// Update handles todo update.
func (ctx Context) Update(w http.ResponseWriter, r *http.Request) error {
	var todo, err = readTodo(w, r)
	if err != nil {
		return BadRequest{err} // 400
	}

	var id = readID(w, r)
	if id != todo.ID {
		return BadRequest{fmt.Errorf("web: mismatch ids")} // 400
	}

	err = ctx.Store.Save(todo)
	if err != nil {
		return err // 500
	}

	return writeJSON(w, todo, http.StatusOK) // 200
}

// Delete handles todo deletion.
func (ctx Context) Delete(w http.ResponseWriter, r *http.Request) error {
	var id = readID(w, r)

	var err = ctx.Store.Delete(id)
	if err != nil {
		return err // 500
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// readTodo returns the todo from the given request.
func readTodo(w http.ResponseWriter, r *http.Request) (*Todo, error) {
	var todo = new(Todo)
	var err = readJSON(r, todo)
	return todo, err
}

// readID returns the "id" variable from the given request.
func readID(w http.ResponseWriter, r *http.Request) string {
	var params = mux.Vars(r)
	var idParam = params["id"]
	return idParam
}

// readStatus returns the "status" variable from the given request.
func readStatus(w http.ResponseWriter, r *http.Request) string {
	var params = mux.Vars(r)
	var status = params["status"]
	return status
}
