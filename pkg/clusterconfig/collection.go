package clusterconfig

// Collection holds APIResources
type Collection struct {
	Items map[string]*APIResource
}

// NewCollection returns an initialized Collection
func NewCollection() *Collection {
	return &Collection{
		Items: make(map[string]*APIResource),
	}
}

// Load read out the API Resources from the given byte array which contains manifests
func (c *Collection) Load(manifest []byte) error {
	// TODO: implement
	return nil
}

// List holds the returned resources from kubectl get
type List struct {
	Kind  string
	Items []APIResource
}

// LoadFromList loads API Resources from a returned list of kubectl
func (c *Collection) LoadFromList(list []byte) error {
	// TODO: implement
	return nil
}
