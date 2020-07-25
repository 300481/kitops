package kitops

// namespaced struct holds the dynamic information if Kind is Namespaced
import (
	"log"
	"os/exec"
	"strings"
)

type namespaced struct {
	resource map[string]bool
}

// declare package variable
var ns *namespaced

func (n *namespaced) namespaced(kind string) bool {
	n.update()
	return n.resource[kind]
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

		n.resource[kind] = namespaced
	}
}

// initialize namespaced
func init() {
	ns = &namespaced{
		resource: make(map[string]bool),
	}
	ns.update()
}
