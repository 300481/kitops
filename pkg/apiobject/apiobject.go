package apiobject

import (
	"errors"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type ErrorMessage string

const (
	ObjectTimeout                = time.Minute
	ErrKindNotFound ErrorMessage = "Kind not found"
)

type ApiObject struct {
	Kind     string
	Metadata struct {
		Name      string
		Namespace string
	}
	exists      bool
	lastUpdated time.Time
	timeout     time.Duration
}

// New parses a YAML of a Kubernetes Object description
// returns an initialized API Object
// returns an error, if the Reader contains no valid yaml object
func New(r io.Reader) (ao *ApiObject, err error) {
	dec := yaml.NewDecoder(r)

	var ob ApiObject
	err = dec.Decode(&ob)

	if err != nil {
		ob.exists = false // TODO needs implementation
		ob.lastUpdated = time.Now()
		ob.timeout = ObjectTimeout
	}

	return &ob, err
}

// Namespaced returns a bool if the API Object is namespaced or not
// returns an error if Kind is not found
// TODO must be made dynamic by fetching supported API-Resources from current cluster
func (ao *ApiObject) Namespaced() (b bool, err error) {
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

	v, ok := namespaced[ao.Kind]
	if ok {
		return v, nil
	}
	return false, errors.New(string(ErrKindNotFound) + ": " + ao.Kind)
}
