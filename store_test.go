package todo

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func withStoreContext(fn func(store Store)) {
	var store = NewStore()
	defer store.Close()
	store.CreateTable()

	fn(store)
}

func findTodo(t *testing.T, store Store, id string) Todo {
	var todo, err = store.Find(id)
	if err != nil {
		// debug.PrintStack()
		t.Fatal(err)
	}

	return todo
}

func saveTodo(t *testing.T, store Store, todo *Todo) {
	var id = todo.ID

	var err = store.Save(todo)
	if err != nil {
		// debug.PrintStack()
		t.Log(todo)
		t.Fatal(err)
	}

	if len(todo.ID) == 0 {
		t.Fatal("todo id error", todo.ID)
	}

	if len(id) != 0 && id != todo.ID {
		t.Fatal("todo id error", todo.ID)
	}
}

func saveTodos(b *testing.B, store Store) []string {
	// fmt.Printf("saveTodos: %d\n", b.N)

	var ids = make([]string, b.N)
	for i := 0; i < b.N; i++ {
		var todo = NewTodo("todo " + strconv.Itoa(i))

		var err = store.Save(todo)
		if err != nil {
			b.Fatal(err)
		}
		ids[i] = todo.ID
	}

	return ids
}

func TestStoreSimple(t *testing.T) {
	withStoreContext(func(store Store) {

		// create
		var todo = &Todo{Title: "todo 1"}
		saveTodo(t, store, todo)

		if todo.Status != "active" {
			t.Fatal("todo status error", todo.Status)
		}
		var created = todo.Created
		if todo.Created.IsZero() {
			t.Fatal("todo created error", todo.Created)
		}

		// list
		var todos = store.List()
		if len(todos) != 1 {
			t.Fatal("todos count error", todos)
		}

		// update
		todo.Complete()
		saveTodo(t, store, todo)

		if !todo.Completed() {
			t.Fatal("todo status error", todo.Status)
		}

		if todo.Created != created {
			t.Fatal("todo created error", todo.Created)
		}

		// read
		var todo2 = findTodo(t, store, todo.ID)
		if !todo2.Equal(*todo) {
			t.Fatalf("equals error:\n%s\n%s\n", todo, todo2)
		}

		// delete
		var err = store.Delete(todo.ID)
		if err != nil {
			t.Fatal(err)
		}

		// read
		_, err = store.Find(todo.ID)
		switch err.(type) {
		case NotFound:
		default:
			t.Fatal("error type error", err)
		}

		// update
		err = store.Save(todo)
		switch err.(type) {
		case NotFound:
		default:
			t.Fatal("error type error", err)
		}

		// delete
		err = store.Delete(todo.ID)
		switch err.(type) {
		case NotFound:
		default:
			t.Fatal("error type error", err)
		}
	})
}

func TestStoreFilter(t *testing.T) {
	withStoreContext(func(store Store) {

		// create
		var todo1 = NewTodo("todo 1")
		saveTodo(t, store, todo1)

		// create
		var todo2 = NewTodo("todo 2")
		saveTodo(t, store, todo2)

		// create
		var todo3 = NewTodo("todo 3")
		todo3.Complete()
		saveTodo(t, store, todo3)

		// list
		var todos = store.List()
		if len(todos) != 3 {
			t.Fatal("todos count error", todos)
		}

		if todos[0].Created.Before(todos[1].Created) {
			t.Fatal("todos sort error", todos)
		}

		// filter
		todos = store.Filter(todo1.Status)
		if len(todos) != 2 {
			t.Fatal("todos count error", todos)
		}

		// filter
		todos = store.Filter(todo3.Status)
		if len(todos) != 1 {
			t.Fatal("todos count error", todos)
		}

		if !todos[0].Equal(*todo3) {
			t.Fatal("todos filter error", todos)
		}

		// clear
		var count, _ = store.Clear(todo3.Status)
		if count != 1 {
			t.Fatal("todos clear error", count)
		}

		// filter
		todos = store.Filter(todo1.Status)
		if len(todos) != 2 {
			t.Fatal("todos count error", todos)
		}

		// filter
		todos = store.Filter(todo3.Status)
		if len(todos) != 0 {
			t.Fatal("todos count error", todos)
		}
	})
}

func BenchmarkStoreC(b *testing.B) {
	withStoreContext(func(store Store) {
		saveTodos(b, store)
	})
}

func BenchmarkStoreR(b *testing.B) {
	withStoreContext(func(store Store) {

		var ids = saveTodos(b, store)
		rand.Seed(time.Now().UnixNano())

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var j = rand.Int63n(int64(b.N))
			var id = ids[j]

			var _, err = store.Find(id)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkStoreU(b *testing.B) {
	withStoreContext(func(store Store) {

		var ids = saveTodos(b, store)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var j = rand.Int63n(int64(b.N))
			var id = ids[j]

			var todo, err = store.Find(id)
			if err != nil {
				b.Fatal(err)
			}

			todo.Title = fmt.Sprintf("[%s]", todo.Title)
			err = store.Save(&todo)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkStoreD(b *testing.B) {
	withStoreContext(func(store Store) {

		var ids = saveTodos(b, store)
		rand.Seed(time.Now().UnixNano())

		var idm = make(map[string]bool)
		for _, id := range ids {
			idm[id] = true
		}

		b.ResetTimer()

		for id := range idm {
			var err = store.Delete(id)
			if err != nil {
				b.Fatal(err)
			}
		}

		var todos = store.List()
		if len(todos) != 0 {
			b.Fatal("todos count error", todos)
		}
	})
}
