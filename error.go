package todo

import "log"
import "net/http"

// BadRequest defines a client error
type BadRequest struct{ error }

// NotFound defines a not found error
type NotFound struct{ error }

// ErrorFunc augments http.HandlerFunc with error return value
type ErrorFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r)
func (f ErrorFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err = f(w, r)
	if err == nil {
		return
	}

	var status int

	switch err.(type) {
	case BadRequest:
		status = http.StatusBadRequest
	case NotFound:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
		log.Printf("error: %s", err)
	}

	http.Error(w, err.Error(), status)
}
