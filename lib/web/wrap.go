package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Wrap the "httprouter" call for use with a context.
func Wrap(handler func(*Context)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Make sure the request body is close upon completion of the handler.
		defer r.Body.Close()

		// Extract username (for logging) and password (reduce redundant calls) and
		// create context object.
		user, password, _ := r.BasicAuth()
		ctx := New(w, r, ps, user, password)

		// Call handler.
		handler(ctx)
	}
}
