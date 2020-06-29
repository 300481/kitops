package clusterconfig

import (
	"fmt"
	"os"
	"path/filepath"

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

	err = config.loadResources()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// loadResources walks the resource directory and loads the resources from yaml files
func (cc *ClusterConfig) loadResources() error {
	walkPath := cc.SourceRepository.Directory + "/" + cc.ResourceDirectory
	err := filepath.Walk(walkPath, cc.readFile)
	return err
}

// readFile loads the resources from yaml files
func (cc *ClusterConfig) readFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}

	if !info.IsDir() {
		matched, _ := filepath.Match("*.yaml", info.Name())
		if matched {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			resource, err := apiresource.New(file)
			if err != nil {
				return err
			}
			cc.APIResources = append(cc.APIResources, resource)
		}
	}

	return nil
}
