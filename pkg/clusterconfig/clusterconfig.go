package clusterconfig

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// ClusterConfig holds the desired and actual API Resources Information
type ClusterConfig struct {
	desired       *Collection
	actual        *Collection
	resourceLabel string
}

// NewClusterConfig returns an initialized ClusterConfig
func NewClusterConfig(label string) *ClusterConfig {
	return &ClusterConfig{
		desired:       NewCollection(),
		actual:        NewCollection(),
		resourceLabel: label,
	}
}

// LoadDesired loads the desired configuration from the given directory
func (c *ClusterConfig) LoadDesired(directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		// if it is a YAML file, load the resources
		if !info.IsDir() {
			matched, _ := filepath.Match("*.yaml", info.Name())
			if !matched {
				return nil
			}

			manifest, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("error reading from file %q: %v\n", path, err)
				return nil
			}

			if err := c.desired.Load(manifest); err != nil {
				log.Printf("error adding manifest: %v\n", err)
				return nil
			}
		}
		return nil
	})
}

// LoadActual loads the actual configuration from the cluster
func (c *ClusterConfig) LoadActual() {
	for kind, isNamespaced := range namespaced {
		var commandArguments []string
		if isNamespaced {
			commandArguments = []string{
				"get",
				kind,
				"-l",
				c.resourceLabel,
				"-o",
				"yaml",
				"-A",
			}
		} else {
			commandArguments = []string{
				"get",
				kind,
				"-l",
				c.resourceLabel,
				"-o",
				"yaml",
			}
		}

		list, err := exec.Command("kubectl", commandArguments...).Output()
		if err != nil {
			log.Println("Error running command: kubectl ", commandArguments)
			continue
		}

		err = c.actual.LoadFromList(list)
		if err != nil {
			log.Println("Error loading API resources from list.")
			continue
		}
	}
}
