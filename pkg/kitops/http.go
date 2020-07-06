package kitops

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// routes sets the routes
func (k *Kitops) routes() {
	k.router.HandleFunc("/healthz", k.healthHandler)
	k.router.HandleFunc("/apply", k.applyHandler)
}

// healthHandler handles the /healthz endpoint
func (k *Kitops) healthHandler(w http.ResponseWriter, r *http.Request) {
	// respond OK
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// applyHandler handles the /apply endpoint
func (k *Kitops) applyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("apply.handler:", r.Method, "request from ", r.RemoteAddr)

	commitID := r.URL.Query().Get("commitid")

	if len(commitID) < 1 {
		log.Println("apply.handler got no commitID")
		handleError(fmt.Errorf("apply.handler got no commitID"), w)
		return
	}

	log.Printf("apply.handler got commitID: %s\n", commitID)

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
