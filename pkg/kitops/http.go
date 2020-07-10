package kitops

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// routes sets the routes
func (k *Kitops) routes() {
	k.router.HandleFunc("/healthz", k.healthHandler).Methods("GET")
	k.router.HandleFunc("/apply", k.applyHandler).Methods("GET")
	k.router.HandleFunc("/clusterconfig", k.clusterConfigHandler).Methods("GET")
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

	if len(commitID) != 40 {
		handleError(fmt.Errorf("apply.handler got no or wrong commitID"), w)
		return
	}

	log.Printf("apply.handler got commitID: %s\n", commitID)

	k.queue.Add(commitID)

	// respond OK
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// clusterConfigHandler writes the ClusterConfig as response
func (k *Kitops) clusterConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(k.queueProcessor.ClusterConfigs)
}

// error handling function
func handleError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	log.Printf("error: %s", err.Error())
}
