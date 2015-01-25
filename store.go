package todo

import (
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Store manages todos storage.
type Store interface {
	List() Todos
	Find(id int64) (Todo, error)
	Save(t *Todo) error
	Delete(id int64) error
	// status
	Filter(status string) Todos
	Clear(status string) (int64, error)
	// close
	Close()
}

// NewStore returns a new Store.
// NewStore reads DATABASE_URL environment variable.
func NewStore() Store {
	// ---------------------------------
	// | driver  | url                 |
	// ---------------------------------
	// | sqlite3 | :memory:            |
	// | sqlite3 | /tmp/todo.sqlite    |
	// ---------------------------------

	var driver = "sqlite3"
	var url = os.Getenv("DATABASE_URL")

	url = strings.TrimSpace(url)
	if len(url) == 0 {
		url = ":memory:"
	}

	return NewSqlStore(driver, url)
}

var DropSchema = `
DROP INDEX IF EXISTS todoStatus;
DROP TABLE IF EXISTS todo;
`

var SqlSchema = `
CREATE TABLE IF NOT EXISTS todo (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    text    TEXT NOT NULL,
    status  TEXT NOT NULL,
    created DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS todoStatus ON todo (status);
`
