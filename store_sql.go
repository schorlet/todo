package todo

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type sqlStore struct {
	db  *sqlx.DB
	url string
}

// NewSqlStore connects to the database specified by driver and url
// and returns a new Store.
func NewSqlStore(driver, url string) Store {
	// log.Printf("database: Opening connection to %s\n", url)
	var db = sqlx.MustOpen(driver, url)

	var tx = db.MustBegin()
	tx.MustExec(SqlSchema)

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

// Find returns the todo with the given id.
func (s sqlStore) Find(id int64) (Todo, error) {
	var t Todo

	var query = `SELECT id, text, status, created
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

	var query = `SELECT id, text, status, created
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

	var query = `SELECT id, text, status, created
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
	var query string

	if t.ID == 0 {
		t.Created = time.Now().UTC()

		query = `INSERT INTO todo (text, status, created)
                VALUES ($1, $2, $3)`

	} else {
		query = `UPDATE todo SET text = $1, status = $2, created = $3
                WHERE id = $4`
	}

	// println(query)

	var tx = s.db.MustBegin()
	defer tx.Rollback()

	var r, err = tx.Exec(query, t.Text, t.Status, t.Created, t.ID)
	if err != nil {
		log.Printf("store: save - %s\n%s\n%s\n", err, query, t)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	if t.ID == 0 {
		t.ID, err = r.LastInsertId()

	} else {
		var count, errc = r.RowsAffected()
		if errc == nil && count == 0 {
			err = NotFound{sql.ErrNoRows}
		} else if err != nil {
			log.Printf("store: save - %s\n", err)
		}
	}

	return err
}

// Delete deletes the todo with the given id.
func (s sqlStore) Delete(id int64) error {
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
	} else if err != nil {
		log.Printf("store: delete - %s\n", err)
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
