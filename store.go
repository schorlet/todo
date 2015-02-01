package todo

import (
	"os"
	"strings"
)

// Store manages todos storage.
type Store interface {
	List() Todos
	Find(id string) (Todo, error)
	Save(t *Todo) error
	Delete(id string) error
	// status
	Filter(status string) Todos
	Clear(status string) (int64, error)
	// store
	Close()
	CreateTable()
}

// NewStore returns a new Store.
// NewStore reads DATABASE_URL environment variable.
func NewStore() Store {
	// var driver = "sqlite3"
	var driver = "rethink"

	var url = os.Getenv("DATABASE_URL")
	url = strings.TrimSpace(url)

	if len(url) == 0 {
		// url = ":memory:"
		url = "localhost:28015/test"
	}

	return OpenStore(driver, url)
}

// OpenStore returns a new Store.
func OpenStore(driver, url string) Store {
	if len(driver) == 0 {
		panic("store: driver is empty")
	}
	if len(url) == 0 {
		panic("store: url is empty")
	}

	var store Store
	if driver == "rethink" {
		store = NewRethinkStore(url)
	} else {
		store = NewSqlStore(driver, url)
	}

	return store
}
