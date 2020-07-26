package kitops

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

const (
	errInvalidYaml = "Invalid YAML"
)

// Collection holds a collection of resources and the
// the corresponding manifests
type Collection struct {
	Items     map[string]*APIResource
	manifests map[string][]byte
}

// NewCollection returns an empty collection of API resources
func NewCollection() *Collection {
	return &Collection{
		Items:     make(map[string]*APIResource),
		manifests: make(map[string][]byte),
	}
}

// LoadFromDirectory loads collections from YAML files
// within the directory recursively
// returns an error if something is wrong with the files
func (c *Collection) LoadFromDirectory(directory string) error {
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

			if err := c.AddFromFile(manifest, path); err != nil {
				log.Printf("error adding manifest: %v\n", err)
				return nil
			}
		}
		return nil
	})
}

// AddFromFile adds API Resources from the given manifest
// manifest is a byte array containing the manifest
// path is the path of the manifest file
// returns an error if the manifest is invalid
func (c *Collection) AddFromFile(manifest []byte, path string) error {
	if len(manifest) == 0 {
		return errors.New(errInvalidYaml)
	}

	c.manifests[path] = make([]byte, len(manifest))
	copy(c.manifests[path], manifest)

	APIResourceContentReader := bytes.NewReader(c.manifests[path])
	for {
		resource, err := NewResource(APIResourceContentReader)
		if err != nil {
			break
		}
		c.Items[resource.Checksum()] = resource
	}

	return nil
}

// List holds the returned resources from kubectl get
type List struct {
	Kind  string
	Items []APIResource
}

// LoadFromList loads API Resources from a byte array
// with a List of API Resources
func (c *Collection) LoadFromList(listContent []byte) error {
	listContentReader := bytes.NewReader(listContent)
	dec := yaml.NewDecoder(listContentReader)

	var list List
	err := dec.Decode(&list)
	if err != nil {
		log.Println("Error loading resources from list of kubectl")
		return err
	}

	if list.Kind != "List" {
		return errors.New("Error: got no list from kubectl")
	}

	for _, apiresource := range list.Items {
		c.Items[apiresource.Checksum()] = &apiresource
	}

	return nil
}

// Exists returns a bool. true if the resources exists in the cluster, false if not
// It returns also false, when there is no information if the resource is namespaced
func (c *Collection) Exists() bool {
	var b = true
	for _, resource := range c.Items {
		b = b && resource.Exists()
	}
	return b
}

// Label labels all resources of the collection in the cluster
func (c *Collection) Label() {
	for _, resource := range c.Items {
		resource.Label()
	}
}
