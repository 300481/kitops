package clusterconfig

import (
	"github.com/300481/kitops/pkg/apiresource"
	"github.com/300481/kitops/pkg/sourcerepo"
)

// ClusterConfig holds all API Resources for a commit id
type ClusterConfig struct {
	CommitID         string
	APIResources     []*apiresource.APIResource
	SourceRepository *sourcerepo.SourceRepo
}

// New returns an initialized *ClusterConfig
// sourceRepo is the Repository with the configuration
// commitID is the commit id of the source repository.
func New(sourceRepo *sourcerepo.SourceRepo, commitID string) (cc *ClusterConfig, err error) {
	config := &ClusterConfig{
		CommitID:         commitID,
		APIResources:     []*apiresource.APIResource{},
		SourceRepository: sourceRepo,
	}
	return config, nil
}
