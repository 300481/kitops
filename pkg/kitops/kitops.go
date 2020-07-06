package kitops

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Kitops is the instance type
type Kitops struct {
	router *mux.Router
}

// New returns a new Kitops instance
func New() *Kitops {
	return &Kitops{
		router: mux.NewRouter(),
	}
}

// Serve runs the application in server mode
func (k *Kitops) Serve() {
	log.Fatal(http.ListenAndServe(":8080", k.router))
}
