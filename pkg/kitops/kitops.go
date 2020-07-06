package kitops

import (
	"log"
	"net/http"
	"os"

	"github.com/300481/kitops/pkg/sourcerepo"
	"github.com/gorilla/mux"
)

// Kitops is the instance type
type Kitops struct {
	router     *mux.Router
	repository *sourcerepo.SourceRepo
}

// New returns a new Kitops instance
func New() *Kitops {
	url := os.Getenv("KITOPS_DEPLOYMENTS_URL")
	if len(url) == 0 {
		url = "https://github.com/300481/kitops-test.git"
	}

	repo, err := sourcerepo.New(url, "/tmp/repo")

	if err != nil {
		log.Printf("unable to get repository: %s\n%v", url, err)
		return nil
	}

	return &Kitops{
		router:     mux.NewRouter(),
		repository: repo,
	}
}

// Serve runs the application in server mode
func (k *Kitops) Serve() {
	k.routes()
	log.Fatal(http.ListenAndServe(":8080", k.router))
}
