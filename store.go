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

type Options struct {
	Driver      string
	Url         string
	CreateTable bool
}

// NewStore returns a new Store.
// NewStore reads DATABASE_URL environment variable.
func NewStore(options ...Options) Store {
	var option = Options{}
	if len(options) >= 1 {
		option = options[0]
	}

	var driver = "sqlite3"
	if len(option.Driver) > 0 {
		driver = option.Driver
	}

	var url = os.Getenv("DATABASE_URL")
	url = strings.TrimSpace(url)

	if len(option.Url) > 0 {
		url = option.Url
	} else if len(url) == 0 {
		url = ":memory:"
	}

	var store = NewSqlStore(driver, url)

	if option.CreateTable {
		store.CreateTable()
	}

	return store
}
