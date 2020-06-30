package clusterconfig

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/300481/kitops/pkg/apiresource"
	"github.com/300481/kitops/pkg/sourcerepo"
	yaml "gopkg.in/yaml.v2"
)

// ClusterConfig holds all API Resources for a commit id
type ClusterConfig struct {
	CommitID          string
	APIResources      []*apiresource.APIResource
	SourceRepository  *sourcerepo.SourceRepo
	ResourceDirectory string
	helmReleases      []*HelmRelease
}

// HelmRelease holds the neccessary information for helm to get
// the cluster resources of the Helm Release
type HelmRelease struct {
	Kind     string
	Metadata struct {
		Namespace string
	}
	Spec struct {
		ReleaseName     string
		HelmVersion     string
		TargetNamespace string
	}
}

// newHelmRelease parses a YAML of a Kubernetes Resource description
// returns an initialized HelmRelease
// returns an error, if the Reader contains no valid yaml
func newHelmRelease(r io.Reader) (helmRelease *HelmRelease, err error) {
	dec := yaml.NewDecoder(r)

	var hr HelmRelease
	err = dec.Decode(&hr)
	if err != nil {
		return nil, err
	}

	if len(hr.Metadata.Namespace) == 0 {
		hr.Metadata.Namespace = "default"
	}

	return &hr, nil
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
		helmReleases:      []*HelmRelease{},
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
				for {
					resource, err := apiresource.New(file)
					if err != nil {
						break
					}
					cc.APIResources = append(cc.APIResources, resource)
				}
				// TODO logic must be rearranged
				file.Close()
				// TODO needs to be included in the for loop or completely moved to loadHelmResources
				if resource.Kind == "HelmRelease" {
					file, err = os.Open(path)
					if err != nil {
						return err
					}
					helmrelease, err := newHelmRelease(file)
					if err != nil {
						return err
					}
					file.Close()
					cc.helmReleases = append(cc.helmReleases, helmrelease)
				}
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
	for _, helmrelease := range cc.helmReleases {
		var command string
		commandArguments := []string{
			"get",
			"manifest",
		}
		if helmrelease.Spec.HelmVersion == "v3" {
			command = "helm3"
			commandArguments = append(commandArguments, "--namespace", helmrelease.Metadata.Namespace)
		} else {
			command = "helm2"
		}
		cmd := exec.Command(command, commandArguments...)
		b, err := cmd.Output()
		if err != nil {
			return err
		}
		r := bytes.NewReader(b)

	}
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
