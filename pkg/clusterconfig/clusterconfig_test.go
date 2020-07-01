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
	commitID := "f0ae1a86eed1923b09dfe3e55b9d657c7dec18ff"
	resourceDirectory := "namespaces"
	sourceDirectory := "testdir"

	log.Println("Clone Repository.")
	repository, err := sourcerepo.New(repoURL, sourceDirectory)
	if err != nil {
		t.Error(err)
	}

	log.Println("Test Checkout.")
	if err := repository.Checkout(commitID); err != nil {
		t.Error(err)
	}

	testClusterConfig := New(repository, commitID, resourceDirectory)
	if err := testClusterConfig.Apply(); err != nil {
		t.Error(err)
	}
}
