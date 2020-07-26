package kitops

// namespaced struct holds the dynamic information if Kind is Namespaced
import (
	"log"
	"os/exec"
	"strings"
)

// clusterKinds holds the namespaced information for Kitops
type clusterKinds struct {
	isNamespaced map[string]bool
}

// declare package variable
var kinds *clusterKinds

// getAll returns the isNamespaced map
func (ck *clusterKinds) getAll() map[string]bool {
	return ck.isNamespaced
}

// namespaced returns a bool if the Kind is namespaced or not
func (ck *clusterKinds) namespaced(kind string) bool {
	ck.update()
	return ck.isNamespaced[kind]
}

// update updates the namespaced information for Kitops
func (ck *clusterKinds) update() {
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

		// skip kind TokenReview
		if kind == "TokenReview" {
			continue
		}
		ck.isNamespaced[kind] = namespaced
	}
}

// initialize namespaced
func init() {
	kinds = &clusterKinds{
		isNamespaced: make(map[string]bool),
	}
	kinds.update()
}
