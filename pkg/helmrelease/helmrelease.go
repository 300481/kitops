package helmrelease

import (
	"io"

	yaml "gopkg.in/yaml.v3"
)

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

// New parses a YAML of a Kubernetes Resource description
// returns an initialized HelmRelease or nil if no Helm Release found in the reader
// returns an error, if the Reader contains no valid yaml
func New(r io.Reader) (helmRelease *HelmRelease, err error) {
	dec := yaml.NewDecoder(r)

	var hr HelmRelease
	err = dec.Decode(&hr)
	if err != nil {
		return nil, err
	}

	if hr.Kind != "HelmRelease" {
		return nil, nil
	}

	if len(hr.Metadata.Namespace) == 0 {
		hr.Metadata.Namespace = "default"
	}

	return &hr, nil
}
