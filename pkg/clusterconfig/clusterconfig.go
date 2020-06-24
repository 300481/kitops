package clusterconfig

import (
	"github.com/300481/kitops/pkg/apiobject"
)

// ClusterConfig holds all API Objects for a commit id
type ClusterConfig struct {
	CommitId   string
	ApiObjects []*apiobject.ApiObject
}

// New returns an initialized *ClusterConfig
// commitId is the commit id of the source repository.
func New(commitId string) (cc *ClusterConfig, err error) {
	config := &ClusterConfig{
		CommitId:   commitId,
		ApiObjects: []*apiobject.ApiObject{},
	}
	return config, nil
}
