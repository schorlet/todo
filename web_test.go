package todo

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func withClientContext(fn func(client *Client, store Store)) {
	var store = NewStore()
	defer store.Close()
	store.CreateTable()

	var handler = NewAppHandler(store)
	var server = httptest.NewServer(handler)
	defer server.Close()

	println(server.URL)
	var client = NewClient(server.URL)

	fn(client, store)
}

func assertStatus(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response status %d but was %d",
			expected, actual)
	}
}

func TestClientSimple(t *testing.T) {
	withClientContext(func(client *Client, store Store) {
		// create
		var todo = NewTodo("todo 1")
		var err = client.Create(todo)
		if err != nil {
			t.Fatal(err)
		}
		assertStatus(t, http.StatusCreated, client.Status)

		// list
		todos, err := client.List()
		if err != nil {
			t.Fatal(err)
		}
		assertStatus(t, http.StatusOK, client.Status)

		if len(todos) != 1 {
			t.Fatal("todos list error", todos)
		}

		// update
		todo.Complete()
		err = client.Update(todo)
		if err != nil {
			t.Fatal(err)
		}
		assertStatus(t, http.StatusOK, client.Status)

		if !todo.Completed() {
			t.Fatal("todo status error", todo.Status)
		}

		// find
		todo2, err := client.Find(todo.ID)
		if err != nil {
			t.Fatal(err)
		}
		assertStatus(t, http.StatusOK, client.Status)

		if !todo2.Equal(*todo) {
			t.Fatalf("equals error:\n%s\n%s\n", todo, todo2)
		}

		// filter
		todos, err = client.Filter(todo.Status)
		if err != nil {
			t.Fatal(err)
		}
		assertStatus(t, http.StatusOK, client.Status)

		if len(todos) != 1 {
			t.Fatal("todos filter error", todos)
		}

		// delete
		err = client.Delete(todo.ID)
		if err != nil {
			t.Fatal(err)
		}
		assertStatus(t, http.StatusNoContent, client.Status)

		// find
		_, err = client.Find(todo.ID)
		if err != nil {
			t.Error(err)
		}
		assertStatus(t, http.StatusNotFound, client.Status)

		// update
		err = client.Update(todo)
		if err != nil {
			t.Error(err)
		}
		assertStatus(t, http.StatusNotFound, client.Status)

		// delete
		err = client.Delete(todo.ID)
		if err != nil {
			t.Error(err)
		}
		assertStatus(t, http.StatusNotFound, client.Status)
	})
}

func TestClientFilter(t *testing.T) {
	withClientContext(func(client *Client, store Store) {
		// create
		client.Create(NewTodo("todo 1"))

		var todo2 = NewTodo("todo 2")
		client.Create(todo2)

		var todo3 = NewTodo("todo 3")
		todo3.Complete()
		client.Create(todo3)

		// list
		todos, _ := client.List()
		if len(todos) != 3 {
			t.Fatal("todos list error", todos)
		}

		if todos[0].Created.Before(todos[1].Created) {
			t.Fatal("todos sort error", todos)
		}

		// filter
		todos, _ = client.Filter(todo2.Status)
		if len(todos) != 2 {
			t.Fatal("todos filter error", todos)
		}

		// filter
		todos, _ = client.Filter(todo3.Status)
		if len(todos) != 1 {
			t.Fatal("todos filter error", todos)
		}

		// clear
		cleared, err := client.Clear(todo3.Status)
		if cleared != 1 {
			t.Fatal("todos clear error", cleared, err)
		}

		// filter
		todos, _ = client.Filter(todo2.Status)
		if len(todos) != 2 {
			t.Fatal("todos filter error", todos)
		}

		// filter
		todos, _ = client.Filter(todo3.Status)
		if len(todos) != 0 {
			t.Fatal("todos filter error", todos)
		}

		// clear
		cleared, _ = client.Clear(todo2.Status)
		if cleared != 2 {
			t.Fatal("todos clear error", cleared)
		}

		// filter
		todos, _ = client.Filter(todo2.Status)
		if len(todos) != 0 {
			t.Fatal("todos filter error", todos)
		}
	})
}

func BenchmarkClientR(b *testing.B) {
	withClientContext(func(client *Client, store Store) {
		var ids = saveTodos(b, store)

		rand.Seed(time.Now().UnixNano())
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var j = rand.Intn(b.N)
			var id = ids[j]

			var _, err = client.Find(id)
			if err != nil {
				b.Error(err)
			}

			if http.StatusOK != client.Status {
				b.Errorf("expected response status %d but was %d",
					http.StatusOK, client.Status)
			}
		}
	})
}
