package todo

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

// Client communicates with the application API.
type Client struct {
    BaseURL string
    router  *mux.Router
    // context
    Status int
    body   []byte
}

// NewClient creates a new todo client with specified baseURL.
func NewClient(baseURL string) *Client {
    return &Client{
        BaseURL: baseURL,
        router:  NewRouter(),
    }
}

func (c *Client) do(method, url string, payload interface{}) error {
    c.Status = http.StatusBadRequest
    c.body = []byte{}

    var body io.Reader
    if payload != nil {
        var b, err = json.Marshal(payload)
        if err != nil {
            return err
        }
        body = bytes.NewReader(b)
    }

    var req, err = http.NewRequest(method, url, body)
    if err != nil {
        return err
    }

    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    res, err := http.DefaultClient.Do(req)
    c.Status = res.StatusCode
    if err != nil {
        return err
    }
    defer res.Body.Close()

    c.body, err = ioutil.ReadAll(res.Body)
    return err
}

// GET /api/todos
func (c *Client) List() (Todos, error) {
    var path, _ = c.router.Get(RouteList).URLPath()
    var url = c.BaseURL + path.String()

    var todos = make(Todos, 0)

    var err = c.do("GET", url, nil)
    if err != nil {
        return todos, err
    }

    if c.Status == http.StatusOK {
        err = json.Unmarshal(c.body, &todos)
    }

    return todos, err
}

// POST /api/todos
func (c *Client) Create(todo *Todo) error {
    var path, _ = c.router.Get(RouteCreate).URLPath()
    var url = c.BaseURL + path.String()

    var err = c.do("POST", url, todo)
    if err != nil {
        return err
    }

    if c.Status == http.StatusCreated {
        err = json.Unmarshal(c.body, todo)
    }

    return err
}

// GET /api/todos{id}
func (c *Client) Find(id int64) (Todo, error) {
    var pairs = []string{"id", strconv.FormatInt(id, 10)}
    var path, _ = c.router.Get(RouteFind).URLPath(pairs...)
    var url = c.BaseURL + path.String()

    var todo = Todo{}

    var err = c.do("GET", url, nil)
    if err != nil {
        return todo, err
    }
    // println(string(c.body))

    if c.Status == http.StatusOK {
        err = json.Unmarshal(c.body, &todo)
    }

    return todo, err
}

// PUT /api/todos/{id}
func (c *Client) Update(todo *Todo) error {
    var pairs = []string{"id", strconv.FormatInt(todo.ID, 10)}
    var path, _ = c.router.Get(RouteUpdate).URLPath(pairs...)
    var url = c.BaseURL + path.String()

    var err = c.do("PUT", url, todo)
    if err != nil {
        return err
    }

    if c.Status == http.StatusOK {
        err = json.Unmarshal(c.body, todo)
    }

    return err
}

// DELETE /api/todos/{id}
func (c *Client) Delete(id int64) error {
    var pairs = []string{"id", strconv.FormatInt(id, 10)}
    var path, _ = c.router.Get(RouteDelete).URLPath(pairs...)
    var url = c.BaseURL + path.String()

    var err = c.do("DELETE", url, nil)
    return err
}

// GET /api/todos/{status}
func (c *Client) Filter(status string) (Todos, error) {
    var pairs = []string{"status", status}
    var path, _ = c.router.Get(RouteFilter).URLPath(pairs...)
    var url = c.BaseURL + path.String()

    var todos = make(Todos, 0)

    var err = c.do("GET", url, nil)
    if err != nil {
        return todos, err
    }

    if c.Status == http.StatusOK {
        err = json.Unmarshal(c.body, &todos)
    }

    return todos, err
}

// DELETE /api/todos/{status}
func (c *Client) Clear(status string) (int64, error) {
    var pairs = []string{"status", status}
    var path, _ = c.router.Get(RouteClear).URLPath(pairs...)
    var url = c.BaseURL + path.String()

    var err = c.do("DELETE", url, nil)
    if err != nil {
        return 0, err
    }

    if c.Status != http.StatusOK {
        return 0, fmt.Errorf("client: expected response status %d but was %d",
            http.StatusOK, c.Status)
    }

    var result map[string]interface{}
    err = json.Unmarshal(c.body, &result)
    if err != nil {
        return 0, err
    }

    value, ok := result["count"]
    if !ok {
        return 0, fmt.Errorf("client: expected count value")
    }

    number, ok := value.(float64)
    if !ok {
        return 0, fmt.Errorf("client: expected count as float64")
    }

    count := int64(number)
    return count, nil
}
