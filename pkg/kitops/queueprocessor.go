package kitops

import (
	"log"

	"github.com/300481/kitops/pkg/clusterconfig"
	"github.com/300481/kitops/pkg/queue"
	"github.com/300481/kitops/pkg/sourcerepo"
)

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
