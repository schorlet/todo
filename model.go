package todo

import (
	"fmt"
	"time"
)

type Todo struct {
	ID      int64     `json:"id"`
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
	t.Status = "complete"
}

func (t Todo) Completed() bool {
	return t.Status == "complete"
}

func (t Todo) String() string {
	return fmt.Sprintf("id:%d, text:%s, status:%s, created:%s",
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
