package kitops

import (
	"io"
	"log"
	"net/http"
)

// routes sets the routes
func (k *Kitops) routes() {
	k.router.HandleFunc("/healthz", k.healthHandler)
}

// healthHandler handles the /healthz endpoint
func (k *Kitops) healthHandler(w http.ResponseWriter, r *http.Request) {
	// respond OK
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// error handling function
func handleError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	log.Printf("error: %s", err.Error())
}
