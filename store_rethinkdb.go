package todo

import (
	"fmt"
	"log"
	"strings"
	"time"

	r "github.com/dancannon/gorethink"
)

type rethinkStore struct {
	session *r.Session
	url     string
}

// NewRethinkStore connects to the database specified by url
// and returns a new Store.
func NewRethinkStore(url string) Store {
	// log.Printf("rethink: Opening connection to %s\n", url)

	// url = "localhost:28015/test"
	var split = strings.Split(url, "/")
	var address = split[0]
	var database = split[1]

	var session, err = r.Connect(r.ConnectOpts{
		Address:  address,
		Database: database,
	})

	if err != nil {
		log.Fatalf("%+v", err)
	}

	return rethinkStore{
		session: session,
		url:     url,
	}
}

// Close closes connection to the rethinkdb server.
func (s rethinkStore) Close() {
	// log.Printf("rethink: Closing connection to %s\n", s.url)
	// s.session.NoReplyWait()
	s.session.Close()
}

// CreateTable drop and create the todo table.
func (s rethinkStore) CreateTable() {
	var err = r.Db("test").TableDrop("Todo").Exec(s.session)
	if err != nil {
		log.Println(err)
	}

	err = r.Db("test").TableCreate("Todo").Exec(s.session)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Find returns the todo with the given id.
func (s rethinkStore) Find(id string) (Todo, error) {
	var t Todo

	var cur, err = r.Table("Todo").Get(id).Run(s.session)
	if err != nil {
		return t, err
	}

	err = cur.One(&t)
	if err == r.ErrEmptyResult {
		err = NotFound{r.ErrEmptyResult}
	}

	return t, err
}

// List returns a list of all todos.
func (s rethinkStore) List() Todos {
	var todos = make(Todos, 0)

	var cur, err = r.Table("Todo").OrderBy(r.Desc("Created")).Run(s.session)
	if err != nil {
		log.Printf("rethink: list - %s\n", err)
		return todos
	}

	err = cur.All(&todos)
	if err != nil {
		log.Printf("rethink: list - %s\n", err)
	}
	return todos
}

// Filter returns a list of todos with the specified status.
func (s rethinkStore) Filter(status string) Todos {
	var todos = make(Todos, 0)

	var cur, err = r.Table("Todo").Filter(
		r.Row.Field("Status").Eq(status),
	).OrderBy(r.Desc("Created")).Run(s.session)

	if err != nil {
		log.Printf("rethink: filter - %s\n", err)
		return todos
	}

	err = cur.All(&todos)
	if err != nil {
		log.Printf("rethink: filter - %s\n", err)
	}
	return todos
}

// Save saves the given todo.
func (s rethinkStore) Save(t *Todo) error {
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
func (s rethinkStore) Insert(t *Todo) error {
	t.Created = time.Now().UTC()

	var res, err = r.Table("Todo").Insert(t).RunWrite(s.session)
	if err != nil {
		log.Printf("rethink: insert - %s\n%s\n", err, t)
		return err
	}

	if res.Errors != 0 {
		return fmt.Errorf(res.FirstError)
	}

	if len(res.GeneratedKeys) == 0 {
		err = fmt.Errorf("GeneratedKeys == 0; %+s", res)
	}

	t.ID = res.GeneratedKeys[0]
	return nil
}

// Update saves the given todo.
func (s rethinkStore) Update(t *Todo) error {
	var cols = map[string]interface{}{
		"Text":   t.Text,
		"Status": t.Status,
	}

	var res, err = r.Table("Todo").Get(t.ID).Update(cols).RunWrite(s.session)
	// log.Printf("%+v", res)

	if err != nil {
		log.Printf("rethink: update - %s\n%s\n", err)
		return err
	}

	if res.Errors != 0 {
		err = fmt.Errorf(res.FirstError)
	} else if res.Replaced == 0 {
		err = NotFound{r.ErrEmptyResult}
	}

	return err
}

// Delete deletes the todo with the given id.
func (s rethinkStore) Delete(id string) error {
	var res, err = r.Table("Todo").Get(id).Delete().RunWrite(s.session)

	if err != nil {
		log.Printf("rethink: delete - %s\n", err)
		return err
	}

	if res.Errors != 0 {
		return fmt.Errorf(res.FirstError)
	}

	if res.Deleted == 0 {
		err = NotFound{r.ErrEmptyResult}
	}

	return err
}

// Clear deletes the todos with the specified status.
func (s rethinkStore) Clear(status string) (int64, error) {
	var res, err = r.Table("Todo").Filter(
		r.Row.Field("Status").Eq(status),
	).Delete().RunWrite(s.session)

	if err != nil {
		log.Printf("rethink: clear - %s\n", err)
		return 0, err
	}

	if res.Errors != 0 {
		return 0, fmt.Errorf(res.FirstError)
	}

	return int64(res.Deleted), nil
}
