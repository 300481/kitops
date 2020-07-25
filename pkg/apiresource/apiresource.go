package apiresource

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// APIResource holds the object information of the API Object
type APIResource struct {
	Kind     string
	Metadata struct {
		Name      string
		Namespace string
	}
}

// NewResource parses a YAML of a Kubernetes Resource description
// returns an initialized API Resource
// returns an error, if the Reader contains no valid yaml
func NewResource(r io.Reader) (resource *APIResource, err error) {
	dec := yaml.NewDecoder(r)

	var ar APIResource
	err = dec.Decode(&ar)
	if err != nil {
		return nil, err
	}

	if len(ar.Metadata.Namespace) == 0 {
		ar.Metadata.Namespace = "default"
	}

	return &ar, nil
}

// Exists returns a bool. true if the Object exists in the cluster, false if not
// It returns also false, when there is no information if the resource is namespaced
func (r *APIResource) Exists() bool {
	var commandArguments []string
	if ns.namespaced(r.Kind) {
		commandArguments = []string{
			"-n",
			r.Metadata.Namespace,
			"get",
			r.Kind,
			r.Metadata.Name,
		}
	} else {
		commandArguments = []string{
			"get",
			r.Kind,
			r.Metadata.Name,
		}
	}

	err := exec.Command("kubectl", commandArguments...).Run()
	if err != nil {
		log.Println("Error running command: kubectl ", commandArguments)
		return false
	}

	return true
}

// Label labels the resource in the cluster
func (r *APIResource) Label() {
	if !r.Exists() {
		log.Println("Warning: resource to label don't exists. Kind: " + r.Kind + " Name: " + r.Metadata.Name + " Namespace: " + r.Metadata.Namespace)
		return
	}

	var commandArguments []string
	if ns.namespaced(r.Kind) {
		commandArguments = []string{
			"-n",
			r.Metadata.Namespace,
			"label",
			"--overwrite",
			r.Kind,
			r.Metadata.Name,
			"managedBy=kitops",
		}
	} else {
		commandArguments = []string{
			"label",
			"--overwrite",
			r.Kind,
			r.Metadata.Name,
			"managedBy=kitops",
		}
	}

	err := exec.Command("kubectl", commandArguments...).Run()
	if err != nil {
		log.Println("Error running command: kubectl ", commandArguments)
		return
	}

	return
}

// Checksum returns a SHA256 checksum of the APIResource as a string
func (r *APIResource) Checksum() string {
	s := fmt.Sprintf("%v", *r)
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

// namespaced struct holds the dynamic information if Kind is Namespaced
type namespaced struct {
	resource map[string]bool
}

// declare package variable
var ns *namespaced

func (n *namespaced) namespaced(kind string) bool {
	n.update()
	return n.resource[kind]
}

func (n *namespaced) set(kind string, namespaced bool) {
	n.resource[kind] = namespaced
}

func (n *namespaced) update() {
	output, err := exec.Command("kubectl", "api-resources").Output()
	if err != nil {
		log.Println(err)
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		if len(line) == 0 {
			break
		}

		s := strings.Fields(line)
		kind := s[len(s)-1]
		namespaced := s[len(s)-2] == "true"

		ns.set(kind, namespaced)
	}
}

// initialize namespaced
func init() {
	ns = &namespaced{
		resource: make(map[string]bool),
	}
	ns.update()
}
