package clusterconfig

import (
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
	APIResources     *apiresource.Collection
	SourceRepository *sourcerepo.SourceRepo
	CommitID         string
}

// New returns an initialized *ClusterConfig
// sourceRepo is the Repository with the configuration
// commitID is the commit id of the source repository.
func New(sourceRepo *sourcerepo.SourceRepo, commitID string) *ClusterConfig {
	return &ClusterConfig{
		APIResources:     apiresource.NewCollection(),
		SourceRepository: sourceRepo,
		CommitID:         commitID,
	}
}

// ApplyManifests applies the manifests stored in the repository
// and checked out with the commitID.
// It returns an error if something goes wrong on apply.
func (cc *ClusterConfig) ApplyManifests() error {
	log.Println("Apply the ClusterConfig.")

	if err := cc.SourceRepository.Checkout(cc.CommitID); err != nil {
		log.Printf("Checkout of repository failed. Commit: %s", cc.CommitID)
		return err
	}

	walkPath := cc.SourceRepository.Directory
	return filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		// if it is a directory, apply the manifests
		if info.IsDir() {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				log.Printf("error reading directory: %q: %v\n", path, err)
				return nil
			}

			// check if it contains YAML files
			var containsYAML = false
			for _, file := range files {
				matched, _ := filepath.Match("*.yaml", file.Name())
				containsYAML = containsYAML || matched
				if matched {
					break
				}
			}

			if containsYAML {
				commandArguments := []string{
					"apply",
					"-f",
					path + "/",
				}

				if err := exec.Command("kubectl", commandArguments...).Run(); err != nil {
					log.Println("Error running command: kubectl ", commandArguments)
					return nil
				}
			}
		}
		return nil
	})
}

// LoadManifests loads the manifests of the checked out repository
// into the ClusterConfig
func (cc *ClusterConfig) LoadManifests() error {
	return cc.APIResources.LoadFromDirectory(cc.SourceRepository.Directory)
}
