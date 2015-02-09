package todo

import (
	"fmt"
	"time"
)

type Todo struct {
	ID      string    `json:"id" gorethink:"id,omitempty"`
	Title   string    `json:"title"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}

type Todos []Todo

func NewTodo(title string) *Todo {
	return &Todo{
		Title:   title,
		Status:  "active",
		Created: time.Now().UTC(),
	}
}

func (t *Todo) Complete() {
	t.Status = "completed"
}

func (t Todo) Completed() bool {
	return t.Status == "completed"
}

func (t Todo) Equal(other Todo) bool {
	return t.ID == other.ID &&
		t.Title == other.Title &&
		t.Status == other.Status &&
		t.Created.Unix() == other.Created.Unix()
}

func (t Todo) String() string {
	return fmt.Sprintf("id:%s, title:%s, status:%s, created:%s",
		t.ID, t.Title, t.Status, t.Created)
}

type ByCreated Todos

func (bc ByCreated) Len() int {
	return len(bc)
}
func (bc ByCreated) Swap(i, j int) {
	bc[i], bc[j] = bc[j], bc[i]
}
func (bc ByCreated) Less(i, j int) bool {
	return bc[i].Created.After(bc[j].Created)
}
