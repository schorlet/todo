package main

import (
    "log"
    "net/http"
    "os"
    "strings"
)

func main() {
    var store = NewStore()
    defer store.Close()

    var handler = NewAppHandler(store)

    var port = os.Getenv("PORT")
    port = strings.TrimSpace(port)
    if len(port) == 0 {
        port = "8000"
    }

    var err = http.ListenAndServe(":"+port, handler)
    if err != nil {
        log.Fatal(err)
    }
}
