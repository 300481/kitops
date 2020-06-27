package clusterconfig

import (
	"github.com/300481/kitops/pkg/apiresource"
	"github.com/300481/kitops/pkg/sourcerepo"
)

// ClusterConfig holds all API Resources for a commit id
type ClusterConfig struct {
	CommitID          string
	APIResources      []*apiresource.APIResource
	SourceRepository  *sourcerepo.SourceRepo
	ResourceDirectory string
}

// New returns an initialized *ClusterConfig
// sourceRepo is the Repository with the configuration
// commitID is the commit id of the source repository.
func New(sourceRepo *sourcerepo.SourceRepo, commitID string, resourceDirectory string) (cc *ClusterConfig, err error) {
	sourceRepo.Checkout(commitID)

	config := &ClusterConfig{
		CommitID:          commitID,
		APIResources:      []*apiresource.APIResource{},
		SourceRepository:  sourceRepo,
		ResourceDirectory: resourceDirectory,
	}
	return config, nil
}
