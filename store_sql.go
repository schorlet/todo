package todo

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type sqlStore struct {
	db  *sqlx.DB
	url string
}

// NewSqlStore connects to the database specified by driver and url
// and returns a new Store.
func NewSqlStore(driver, url string) Store {
	// ---------------------------------
	// | driver  | url                 |
	// ---------------------------------
	// | sqlite3 | :memory:            |
	// | sqlite3 | /tmp/todo.sqlite    |
	// ---------------------------------

	// log.Printf("database: Opening connection to %s\n", url)
	var db = sqlx.MustOpen(driver, url)

	var tx = db.MustBegin()
	tx.MustExec(CreateTable)

	var err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return sqlStore{
		db:  db,
		url: url,
	}
}

// Close closes connection to the sql store.
func (s sqlStore) Close() {
	// log.Printf("database: Closing connection to %s\n", s.url)
	s.db.Close()
}

// CreateTable drop and create the todo table.
func (s sqlStore) CreateTable() {
	var tx = s.db.MustBegin()

	tx.MustExec(DropTable)
	tx.MustExec(CreateTable)

	var err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// Find returns the todo with the given id.
func (s sqlStore) Find(id string) (Todo, error) {
	var t Todo

	var query = `SELECT id, title, status, created
        FROM todo
        WHERE id = $1`
	// println(query)

	var err = s.db.Get(&t, query, id)

	if err == sql.ErrNoRows {
		err = NotFound{sql.ErrNoRows}
	} else if err != nil {
		log.Printf("store: find - %s\n", err)
	}

	return t, err
}

// List returns a list of all todos.
func (s sqlStore) List() Todos {
	var todos = make(Todos, 0)

	var query = `SELECT id, title, status, created
        FROM todo
        ORDER BY created DESC`
	// println(query)

	var err = s.db.Select(&todos, query)

	if err != nil {
		log.Printf("store: list - %s\n", err)
	}

	return todos
}

// Filter returns a list of todos with the specified status.
func (s sqlStore) Filter(status string) Todos {
	var todos = make(Todos, 0)

	var query = `SELECT id, title, status, created
        FROM todo
        WHERE status = $1
        ORDER BY created DESC`
	// println(query, status)

	var err = s.db.Select(&todos, query, status)

	if err != nil {
		log.Printf("store: filter - %s\n", err)
	}

	return todos
}

// Save saves the given todo.
func (s sqlStore) Save(t *Todo) error {
	if len(t.Status) == 0 {
		t.Status = "active"
	}

	if len(t.ID) == 0 {
		t.Created = time.Now().UTC()
		return s.Insert(t)
	}

	return s.Update(t)
}

// Insert saves the given todo.
func (s sqlStore) Insert(t *Todo) error {
	var query = `INSERT INTO todo (title, status, created)
                VALUES ($1, $2, $3)`

	var tx = s.db.MustBegin()
	defer tx.Rollback()

	var r, err = tx.Exec(query, t.Title, t.Status, t.Created)
	if err != nil {
		log.Printf("store: insert - %s\n%s\n%s\n", err, query, t)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	autoIncr, err := r.LastInsertId()
	t.ID = strconv.FormatInt(autoIncr, 10)
	return err
}

// Update saves the given todo.
func (s sqlStore) Update(t *Todo) error {
	var query = `UPDATE todo SET title = $1, status = $2
                WHERE id = $3`

	var tx = s.db.MustBegin()
	defer tx.Rollback()

	var r, err = tx.Exec(query, t.Title, t.Status, t.ID)
	if err != nil {
		log.Printf("store: update - %s\n%s\n%s\n", err, query, t)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err == nil && count == 0 {
		err = NotFound{sql.ErrNoRows}
	}
	return err
}

// Delete deletes the todo with the given id.
func (s sqlStore) Delete(id string) error {
	var query = `DELETE FROM todo WHERE id = $1`
	// println(query)

	var tx = s.db.MustBegin()
	defer tx.Rollback()

	var r, err = tx.Exec(query, id)
	if err != nil {
		log.Printf("store: delete - %s\n%s\n", err, query)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err == nil && count == 0 {
		err = NotFound{sql.ErrNoRows}
	}

	return err
}

// Clear deletes the todos with the specified status.
func (s sqlStore) Clear(status string) (int64, error) {
	var query = `DELETE FROM todo WHERE status = $1`
	// println(query)

	var tx = s.db.MustBegin()
	defer tx.Rollback()

	var r, err = tx.Exec(query, status)
	if err != nil {
		log.Printf("store: clear - %s\n%s\n", err, query)
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	count, err := r.RowsAffected()
	return count, err
}

const DropTable = `
DROP INDEX IF EXISTS todoStatus;
DROP TABLE IF EXISTS todo;
`

const CreateTable = `
CREATE TABLE IF NOT EXISTS todo (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    title   TEXT NOT NULL,
    status  TEXT NOT NULL,
    created DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS todoStatus ON todo (status);
`
