package kitops

import (
	"log"
	"net/http"
	"os"

	"github.com/300481/kitops/pkg/queue"
	"github.com/300481/kitops/pkg/sourcerepo"
	"github.com/gorilla/mux"
)

// Kitops is the instance type
type Kitops struct {
	router         *mux.Router
	queue          *queue.Queue
	queueProcessor *QueueProcessor
}

// New returns a new Kitops instance
func New() *Kitops {
	url := os.Getenv("KITOPS_DEPLOYMENTS_URL")
	if len(url) == 0 {
		// set default URL
		url = "https://github.com/300481/kitops-test.git"
	}

	repo, err := sourcerepo.New(url, "/tmp/repo")

	if err != nil {
		log.Printf("unable to get repository: %s\n%v", url, err)
		return nil
	}

	qp := &QueueProcessor{
		ClusterConfigs: make(map[string]*ClusterConfig),
		repository:     repo,
	}

	q := queue.New(qp)

	return &Kitops{
		router:         mux.NewRouter(),
		queue:          q,
		queueProcessor: qp,
	}
}

// Serve runs the application in server mode
func (k *Kitops) Serve() {
	k.routes()
	log.Fatal(http.ListenAndServe(":8080", k.router))
}
