package clusterconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/300481/kitops/pkg/apiresource"
	"github.com/300481/kitops/pkg/helmrelease"
	"github.com/300481/kitops/pkg/sourcerepo"
)

// ClusterConfig holds all API Resources for a commit id
type ClusterConfig struct {
	CommitID          string
	APIResources      []*apiresource.APIResource
	HelmReleases      []*helmrelease.HelmRelease
	SourceRepository  *sourcerepo.SourceRepo
	ResourceDirectory string
}

// New returns an initialized *ClusterConfig
// sourceRepo is the Repository with the configuration
// commitID is the commit id of the source repository.
func New(sourceRepo *sourcerepo.SourceRepo, commitID string, resourceDirectory string) *ClusterConfig {
	return &ClusterConfig{
		CommitID:          commitID,
		APIResources:      []*apiresource.APIResource{},
		HelmReleases:      []*helmrelease.HelmRelease{},
		SourceRepository:  sourceRepo,
		ResourceDirectory: resourceDirectory,
	}
}

// Apply applies the configuration stored in the repositories
// and checked out with the commitID.
// It also loads all Resources into the ClusterConfig.
func (cc *ClusterConfig) Apply() error {
	if err := cc.SourceRepository.Checkout(cc.CommitID); err != nil {
		return err
	}
	walkPath := cc.SourceRepository.Directory + "/" + cc.ResourceDirectory
	walkErr := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		// if it is a YAML file, load the resources
		if !info.IsDir() {
			matched, _ := filepath.Match("*.yaml", info.Name())
			if matched {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				// Create a Reader for the API Resources
				APIResourceContentReader := bytes.NewReader(content)
				for {
					resource, err := apiresource.New(APIResourceContentReader)
					if err != nil {
						break
					}
					cc.APIResources = append(cc.APIResources, resource)
				}

				// Create a Reader for the HelmReleases
				HelmReleaseContentReader := bytes.NewReader(content)
				for {
					helmRelease, err := helmrelease.New(HelmReleaseContentReader)
					if err != nil {
						break
					}
					if helmRelease != nil {
						cc.HelmReleases = append(cc.HelmReleases, helmRelease)
					}
				}
			}
			// if it is a directory, apply the YAML files with kubectl
		} else {
			// write out errors if they occur
			if err := cc.applyKubectl(path); err != nil {
				log.Println(err)
			}
		}
		return cc.loadHelmResources()
	})
	return walkErr
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

// loadHelmResources loads the resources which results from a HelmRelease
func (cc *ClusterConfig) loadHelmResources() error {
	var command string
	commandArguments := []string{
		"get",
		"manifest",
	}

	for _, helmrelease := range cc.HelmReleases {
		if helmrelease.Spec.HelmVersion == "v3" {
			command = "helm3"
			commandArguments = append(commandArguments, "--namespace", helmrelease.Metadata.Namespace)
		} else {
			command = "helm2"
		}

		b, err := exec.Command(command, commandArguments...).Output()
		if err != nil {
			return err
		}

		APIResourceContentReader := bytes.NewReader(b)
		for {
			resource, err := apiresource.New(APIResourceContentReader)
			if err != nil {
				break
			}
			cc.APIResources = append(cc.APIResources, resource)
		}
	}
	return nil
}
