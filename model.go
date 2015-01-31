package todo

import (
	"fmt"
	"time"
)

type Todo struct {
	ID      string    `json:"id" gorethink:"id,omitempty"`
	Text    string    `json:"text"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}

type Todos []Todo

func NewTodo(text string) *Todo {
	return &Todo{
		Text:    text,
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
		t.Text == other.Text &&
		t.Status == other.Status &&
		t.Created.Unix() == other.Created.Unix()
}

func (t Todo) String() string {
	return fmt.Sprintf("id:%s, text:%s, status:%s, created:%s",
		t.ID, t.Text, t.Status, t.Created)
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
