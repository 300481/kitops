package kitops

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/300481/kitops/pkg/sourcerepo"
)

// ClusterConfig holds all API Resources for a commit id
type ClusterConfig struct {
	APIResources     *Collection
	SourceRepository *sourcerepo.SourceRepo
	CommitID         string
	ResourceLabel    string
}

// NewClusterConfig returns an initialized *ClusterConfig
// sourceRepo is the Repository with the configuration
// commitID is the commit id of the source repository.
func NewClusterConfig(sourceRepo *sourcerepo.SourceRepo, commitID string) *ClusterConfig {
	resourceLabel := "managedBy=" + strings.ReplaceAll(strings.ReplaceAll(sourceRepo.URL, ":", "-"), "/", "-")
	return &ClusterConfig{
		APIResources:     NewCollection(resourceLabel),
		SourceRepository: sourceRepo,
		CommitID:         commitID,
		ResourceLabel:    resourceLabel,
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

				log.Println("Running command: kubectl ", commandArguments)
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
	tempCollection := NewCollection(cc.ResourceLabel)
	clusterkinds := kinds.getAll()

	// get all labelled resources, put them in a temporary collection
	for kind, namespaced := range clusterkinds {
		var commandArguments []string
		if namespaced {
			commandArguments = []string{
				"get",
				kind,
				"-l",
				cc.ResourceLabel,
				"-o",
				"yaml",
				"-A",
			}
		} else {
			commandArguments = []string{
				"get",
				kind,
				"-l",
				cc.ResourceLabel,
				"-o",
				"yaml",
			}
		}

		list, err := exec.Command("kubectl", commandArguments...).Output()
		if err != nil {
			//log.Println("Error running command: kubectl ", commandArguments)
			continue
		}

		err = tempCollection.LoadFromList(list)
		if err != nil {
			log.Println("Error loading API resources from list.")
			continue
		}
	}

	// compare them with the resources of the current ClusterConfig
	// if not in the current ClusterConfig, delete it
	for hash, item := range tempCollection.Items {
		log.Printf("Cleanup check for Checksum: %s %s %s %s", hash, item.Kind, item.Metadata.Name, item.Metadata.Namespace)
		_, ok := cc.APIResources.Items[hash]
		if !ok {
			item.Delete()
		}
	}

	return
}
