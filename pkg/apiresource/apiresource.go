package apiresource

import (
	"errors"
	"io"
	"log"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type errorMessage string

const (
	errKindNotFound errorMessage = "Kind not found"
)

// APIResource holds the object information of the API Object
type APIResource struct {
	Kind     string
	Metadata struct {
		Name      string
		Namespace string
	}
}

// Collection holds a collection of resources and the content of
// New parses a YAML of a Kubernetes Resource description
// returns an initialized API Resource
// returns an error, if the Reader contains no valid yaml
func New(r io.Reader) (resource *APIResource, err error) {
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
	return false, errors.New(string(errKindNotFound) + ": " + r.Kind)
}

// Exists returns a bool. true if the Object exists in the cluster, false if not
// It returns also false, when there is no information if the object is namespaced
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

	cmd := exec.Command("kubectl", commandArguments...)

	err = cmd.Run()
	if err != nil {
		log.Println("Error running command: kubectl ", commandArguments)
		return false
	}

	return true
}
