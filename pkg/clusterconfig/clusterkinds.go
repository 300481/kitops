package clusterconfig

// namespaced struct holds the dynamic information if Kind is Namespaced
import (
	"log"
	"os/exec"
	"strings"
)

// declare package variable
var namespaced = make(map[string]bool)

// initialize namespaced
func init() {
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

		namespaced[kind] = s[len(s)-2] == "true"
	}
}

// invalidKind checks if the given Kind is an invalid one for a collection,
// returns a bool
func invalidKind(kind string) bool {
	kinds := []string{"ComponentStatus"}

	for _, item := range kinds {
		if kind == item {
			return true
		}
	}

	return false
}
