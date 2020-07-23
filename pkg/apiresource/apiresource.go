package apiresource

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	errKindNotFound = "Kind not found"
	errInvalidYaml  = "Invalid YAML"
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

			if err := c.Add(manifest, path); err != nil {
				log.Printf("error adding manifest: %v\n", err)
				return nil
			}
		}
		return nil
	})
}

// Add adds API Resources from the given manifest
// manifest is a byte array containing the manifest
// path is the path of the manifest file
// returns an error if the manifest is invalid
func (c *Collection) Add(manifest []byte, path string) error {
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

// Exists returns a bool. true if the resources exists in the cluster, false if not
// It returns also false, when there is no information if the resource is namespaced
func (c *Collection) Exists() bool {
	var b = true
	for _, resource := range c.Items {
		b = b && resource.Exists()
	}
	return b
}

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

// Namespaced returns a bool if the API Resource is namespaced or not
// returns an error if Kind is not found
// TODO must be made dynamic by fetching supported API-Resources from current cluster
func (r *APIResource) Namespaced() (b bool, err error) {
	namespaced := map[string]bool{
		"Alertmanager":                   true,
		"Binding":                        true,
		"Certificate":                    true,
		"CertificateRequest":             true,
		"Challenge":                      true,
		"ConfigMap":                      true,
		"ControllerRevision":             true,
		"CronJob":                        true,
		"DaemonSet":                      true,
		"Deployment":                     true,
		"EndpointSlice":                  true,
		"Endpoints":                      true,
		"Event":                          true,
		"HelmRelease":                    true,
		"HorizontalPodAutoscaler":        true,
		"Ingress":                        true,
		"Issuer":                         true,
		"Job":                            true,
		"Lease":                          true,
		"LimitRange":                     true,
		"LocalSubjectAccessReview":       true,
		"NetworkPolicy":                  true,
		"NetworkSet":                     true,
		"Order":                          true,
		"PersistentVolumeClaim":          true,
		"Pod":                            true,
		"PodDisruptionBudget":            true,
		"PodMonitor":                     true,
		"PodTemplate":                    true,
		"Prometheus":                     true,
		"PrometheusRule":                 true,
		"ReplicaSet":                     true,
		"ReplicationController":          true,
		"ResourceQuota":                  true,
		"Role":                           true,
		"RoleBinding":                    true,
		"Secret":                         true,
		"Service":                        true,
		"ServiceAccount":                 true,
		"ServiceMonitor":                 true,
		"StatefulSet":                    true,
		"ThanosRuler":                    true,
		"APIService":                     false,
		"BGPConfiguration":               false,
		"BGPPeer":                        false,
		"BlockAffinity":                  false,
		"CSIDriver":                      false,
		"CSINode":                        false,
		"CertificateSigningRequest":      false,
		"ClusterInformation":             false,
		"ClusterIssuer":                  false,
		"ClusterRole":                    false,
		"ClusterRoleBinding":             false,
		"ComponentStatus":                false,
		"CustomResourceDefinition":       false,
		"FelixConfiguration":             false,
		"GlobalNetworkPolicy":            false,
		"GlobalNetworkSet":               false,
		"HostEndpoint":                   false,
		"IPAMBlock":                      false,
		"IPAMConfig":                     false,
		"IPAMHandle":                     false,
		"IPPool":                         false,
		"MutatingWebhookConfiguration":   false,
		"Namespace":                      false,
		"Node":                           false,
		"PersistentVolume":               false,
		"PodSecurityPolicy":              false,
		"PriorityClass":                  false,
		"RuntimeClass":                   false,
		"SelfSubjectAccessReview":        false,
		"SelfSubjectRulesReview":         false,
		"StorageClass":                   false,
		"SubjectAccessReview":            false,
		"TokenReview":                    false,
		"ValidatingWebhookConfiguration": false,
		"VolumeAttachment":               false,
	}

	v, ok := namespaced[r.Kind]
	if ok {
		return v, nil
	}
	return false, errors.New(errKindNotFound + ": " + r.Kind)
}

// Exists returns a bool. true if the Object exists in the cluster, false if not
// It returns also false, when there is no information if the resource is namespaced
func (r *APIResource) Exists() bool {
	namespaced, err := r.Namespaced()
	if err != nil {
		return false
	}

	var commandArguments []string
	if namespaced {
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

	err = exec.Command("kubectl", commandArguments...).Run()
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

	namespaced, err := r.Namespaced()
	if err != nil {
		log.Println(err)
		return
	}

	var commandArguments []string
	if namespaced {
		commandArguments = []string{
			"-n",
			r.Metadata.Namespace,
			"label",
			r.Kind,
			r.Metadata.Name,
			"managedBy=kitops",
		}
	} else {
		commandArguments = []string{
			"label",
			r.Kind,
			r.Metadata.Name,
			"managedBy=kitops",
		}
	}

	err = exec.Command("kubectl", commandArguments...).Run()
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
