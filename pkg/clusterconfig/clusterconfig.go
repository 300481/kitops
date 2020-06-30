package clusterconfig

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

// helmRelease holds the neccessary information for helm to get
// the cluster resources of the Helm Release
type helmRelease struct {
	Kind     string
	Metadata struct {
		Namespace string
	}
	Spec struct {
		releaseName string
	}
}

// New returns an initialized *ClusterConfig
// sourceRepo is the Repository with the configuration
// commitID is the commit id of the source repository.
func New(sourceRepo *sourcerepo.SourceRepo, commitID string, resourceDirectory string) *ClusterConfig {
	return &ClusterConfig{
		CommitID:          commitID,
		APIResources:      []*apiresource.APIResource{},
		SourceRepository:  sourceRepo,
		ResourceDirectory: resourceDirectory,
	}
}

// Apply applies the configuration stored in the repositories
// and checked out with the commitID
func (cc *ClusterConfig) Apply() error {
	if err := cc.SourceRepository.Checkout(cc.CommitID); err != nil {
		return err
	}

	if err := cc.loadKubectlResourcesAndApply(); err != nil {
		return err
	}

	if err := cc.loadHelmResources(); err != nil {
		return err
	}

	return nil
}

// loadResources walks the resource directory and loads the resources from yaml files
func (cc *ClusterConfig) loadKubectlResourcesAndApply() error {
	walkPath := cc.SourceRepository.Directory + "/" + cc.ResourceDirectory
	return filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		// if it is a YAML file, load the resources
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
			// if it is a directory, apply the YAML files with kubectl
		} else {
			// write out errors if they occur
			if err := cc.applyKubectl(path); err != nil {
				log.Println(err)
			}
		}
		return nil
	})
}

// loadHelmResources loads the resources which results from a HelmRelease
func (cc *ClusterConfig) loadHelmResources() error {
	// TODO: needs implementation
	return nil
}

// applyKubectl runs kubectl apply with the given path
// returns error if something goes wrong with this path
func (cc *ClusterConfig) applyKubectl(path string) error {
	// first check if directory contains yaml files
	var containsYAML = false

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		matched, _ := filepath.Match("*.yaml", file.Name())
		if matched {
			containsYAML = true
		}
	}

	// if it contains YAML files, apply them
	if containsYAML {
		commandArguments := []string{
			"apply",
			"-f",
			path + "/",
		}

		if err := exec.Command("kubectl", commandArguments...).Run(); err != nil {
			log.Println("Error running command: kubectl ", commandArguments)
			return err
		}
	}

	return nil
}
