package clusterconfig

// APIResource is a basic API Resource description
type APIResource struct {
	Kind     string
	Metadata struct {
		Name      string
		Namespace string
	}
}

// NewAPIResource returns an initialized APIResource
func NewAPIResource() *APIResource {
	return &APIResource{}
}
