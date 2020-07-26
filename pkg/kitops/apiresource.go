package kitops

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os/exec"

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
	if kinds.namespaced(r.Kind) {
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
	if kinds.namespaced(r.Kind) {
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

	log.Printf("Label Resource %s %s %s %s", r.Checksum(), r.Kind, r.Metadata.Name, r.Metadata.Namespace)

	return
}

// Checksum returns a SHA256 checksum of the APIResource as a string
func (r *APIResource) Checksum() string {
	s := fmt.Sprintf("%v", *r)
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

// Delete deletes the resource from the cluster
func (r *APIResource) Delete() {
	if !r.Exists() {
		return
	}

	var commandArguments []string
	if kinds.namespaced(r.Kind) {
		commandArguments = []string{
			"-n",
			r.Metadata.Namespace,
			"delete",
			r.Kind,
			r.Metadata.Name,
		}
	} else {
		commandArguments = []string{
			"delete",
			r.Kind,
			r.Metadata.Name,
		}
	}

	err := exec.Command("kubectl", commandArguments...).Run()
	if err != nil {
		log.Println("Error running command: kubectl ", commandArguments)
	}

	log.Printf("Cleanup Resource Kind: %s Name: %s Namespace: %s", r.Kind, r.Metadata.Name, r.Metadata.Namespace)

	return
}
