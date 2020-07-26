package kitops

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/300481/kitops/pkg/sourcerepo"
)

// ClusterConfig holds all API Resources for a commit id
type ClusterConfig struct {
	APIResources     *Collection
	SourceRepository *sourcerepo.SourceRepo
	CommitID         string
}

// NewClusterConfig returns an initialized *ClusterConfig
// sourceRepo is the Repository with the configuration
// commitID is the commit id of the source repository.
func NewClusterConfig(sourceRepo *sourcerepo.SourceRepo, commitID string) *ClusterConfig {
	return &ClusterConfig{
		APIResources:     NewCollection(),
		SourceRepository: sourceRepo,
		CommitID:         commitID,
	}
}

// checkout checks the commitID out
func (cc *ClusterConfig) checkout() error {
	if err := cc.SourceRepository.Checkout(cc.CommitID); err != nil {
		log.Printf("checkout of repository failed. Commit: %s", cc.CommitID)
		return err
	}
	return nil
}

// ApplyManifests applies the manifests stored in the repository
// and checked out with the commitID.
// It returns an error if something goes wrong on apply.
func (cc *ClusterConfig) ApplyManifests() error {
	log.Println("Apply the ClusterConfig.")

	if err := cc.checkout(); err != nil {
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
	if err := cc.checkout(); err != nil {
		return err
	}
	return cc.APIResources.LoadFromDirectory(cc.SourceRepository.Directory)
}

// Label labels all resources of this ClusterConfig in the Cluster
func (cc *ClusterConfig) Label() {
	cc.APIResources.Label()
	return
}

// Clean cleans the cluster from resources which are not in the ClusterConfig,
// but managed by Kitops
func (cc *ClusterConfig) Clean() {
	tempCollection := NewCollection()
	clusterkinds := kinds.getAll()

	// get all labelled resources, put them in a temporary collection
	for kind, namespaced := range clusterkinds {
		var commandArguments []string
		if namespaced {
			commandArguments = []string{
				"get",
				kind,
				"-l",
				"managedBy=kitops",
				"-o",
				"yaml",
				"-A",
			}
		} else {
			commandArguments = []string{
				"get",
				kind,
				"-l",
				"managedBy=kitops",
				"-o",
				"yaml",
			}
		}

		list, err := exec.Command("kubectl", commandArguments...).Output()
		if err != nil {
			log.Println("Error running command: kubectl ", commandArguments)
			continue
		}

		err = tempCollection.LoadFromList(list) // TODO implement LoadFromList()
		if err != nil {
			log.Println("Error loading API resources from list.")
			return
		}
	}

	// compare them with the resources of the current ClusterConfig
	// if not in the current ClusterConfig, delete it
	for hash, item := range tempCollection.Items {
		_, ok := cc.APIResources.Items[hash]
		if !ok {
			item.Delete()
		}
	}

	return
}
