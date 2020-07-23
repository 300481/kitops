package kitops

import (
	"log"
	"net/http"
	"os"

	"github.com/300481/kitops/pkg/clusterconfig"
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
		ClusterConfigs: make(map[string]*clusterconfig.ClusterConfig),
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

// QueueProcessor is the instance for processsing the queue items
type QueueProcessor struct {
	ClusterConfigs map[string]*clusterconfig.ClusterConfig
	repository     *sourcerepo.SourceRepo
}

// Process processes new queued commitIDs
func (qp *QueueProcessor) Process(q *queue.Queue) {
	commitID := q.StartNext().(string)

	// create a new ClusterConfig
	cc := clusterconfig.New(qp.repository, commitID)
	qp.ClusterConfigs[commitID] = cc

	// apply the manifests
	if err := cc.ApplyManifests(); err != nil {
		log.Printf("failed to apply manifests of commitID: %s", commitID)
		q.Finish(false)
	} else {
		q.Finish(true)
	}

	// load the manifests in the ClusterConfig
	if err := cc.LoadManifests(); err != nil {
		log.Printf("failed to load manifests of commitID: %s", commitID)
	}

	// label the api resources
	qp.ClusterConfigs[commitID].Label()
}
