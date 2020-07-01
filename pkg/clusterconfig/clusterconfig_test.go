package clusterconfig

import (
	"log"
	"testing"

	"github.com/300481/kitops/pkg/sourcerepo"
)

func TestNew(t *testing.T) {
	log.Println("Test New()")
	commitID := "f0ae1a86eed1923b09dfe3e55b9d657c7dec18ff"
	resourceDirectory := "releases"

	testClusterConfig := New(nil, commitID, resourceDirectory)

	if testClusterConfig.CommitID != commitID || testClusterConfig.ResourceDirectory != resourceDirectory {
		t.Error("ClusterConfig not initialized correctly.")
		log.Printf("%+v\n", testClusterConfig)
	}
}

func TestApply(t *testing.T) {
	log.Println("Test Apply()")

	repoURL := "https://github.com/300481/kitops-test.git"
	commitID := "045d4485d54af656b11b05b2e26697cac7df8b76"
	resourceDirectory := "namespaces"
	sourceDirectory := "testdir"

	log.Println("Clone Repository.")
	repository, err := sourcerepo.New(repoURL, sourceDirectory)
	if err != nil {
		t.Error(err)
	}

	testClusterConfig := New(repository, commitID, resourceDirectory)
	if err := testClusterConfig.Apply(); err != nil {
		t.Error(err)
	}

	if testClusterConfig.CommitID != commitID || testClusterConfig.ResourceDirectory != resourceDirectory {
		t.Errorf("CommitID or ResourceDirectory wrong.\nClusterConfig.CommitID: %s\nClusterConfig.ResourceDirectory: %s", testClusterConfig.CommitID, testClusterConfig.ResourceDirectory)
	}

	m := make(map[string]string)
	m["default"] = "Namespace"
	m["default-test"] = "Namespace"
	m["kube-node-lease"] = "Namespace"
	m["kube-public"] = "Namespace"
	m["kube-system"] = "Namespace"

	for _, resource := range testClusterConfig.APIResources {
		if m[resource.Metadata.Name] != resource.Kind {
			t.Errorf("Resource Name: %s is not the right Kind: %s", resource.Metadata.Name, resource.Kind)
		}
		if !resource.Exists() {
			t.Errorf("Resource Name: %s Kind: %s was not created in the cluster.", resource.Metadata.Name, resource.Kind)
		}
	}
}
