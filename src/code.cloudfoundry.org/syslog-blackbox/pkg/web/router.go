package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// CountGetter is the function used by the routers handlers to get counts for a
// given ID.
type CountGetter func(id string) (primers, messages int)

// NewRouter configures and returns a new router as an http.Handler
func NewRouter(getter CountGetter) http.Handler {
	r := mux.NewRouter()

	r.Handle("/get/{id}", MessageCountHandler(getter)).Methods(http.MethodGet)
	r.Handle("/get-prime/{id}", PrimeCountHandler(getter)).Methods(http.MethodGet)

	return handlers.LoggingHandler(os.Stdout, r)
}

// MessageCountHandler will return the number of messages for a given ID.
func MessageCountHandler(getter CountGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, msgs := getter(id)

		w.Write([]byte(fmt.Sprint(msgs)))
	})
}

// PrimeCountHandler will return the number of primers for a given ID.
func PrimeCountHandler(getter CountGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		primers, _ := getter(id)

		w.Write([]byte(fmt.Sprint(primers)))
	})
}
