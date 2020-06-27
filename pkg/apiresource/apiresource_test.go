package apiresource

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	var TestYaml = `
---
kind: Test
metadata:
  name: Test
  namespace: Test
`

	var TestResource = APIResource{
		Kind: "Test",
		Metadata: struct {
			Name      string
			Namespace string
		}{
			Name:      "Test",
			Namespace: "Test",
		},
	}

	r := bytes.NewReader([]byte(TestYaml))
	resource, err := New(r)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Testing 'TestNew()'\nTestResource:\n%+v\nNew resource:\n%+v\n", TestResource, *resource)
	if *resource != TestResource {
		t.Errorf("Test resources don't match.\nTestResource:\n%+v\nNew resource:\n%+v\n", TestResource, *resource)
	}
}

func TestNamespaced(t *testing.T) {
	fmt.Printf("Testing 'Namespaced()'\n")

	var TestYaml1 = `
---
kind: Test
metadata:
  name: Test
  namespace: Test
`
	r1 := bytes.NewReader([]byte(TestYaml1))
	resource1, err := New(r1)
	if err != nil {
		t.Error(err)
	}

	_, err = resource1.Namespaced()
	if err == nil {
		t.Errorf("Kind 'Test' should not be found.\n%+v\n", resource1)
	}

	var TestYaml2 = `
---
kind: Deployment
metadata:
  name: Test
  namespace: Test
`
	r2 := bytes.NewReader([]byte(TestYaml2))
	resource2, err := New(r2)
	if err != nil {
		t.Error(err)
	}

	namespaced2, _ := resource2.Namespaced()
	if !namespaced2 {
		t.Errorf("Kind 'Deployment' should be namespaced.\n%+v\n", resource2)
	}

	var TestYaml3 = `
---
kind: StorageClass
metadata:
  name: Test
`
	r3 := bytes.NewReader([]byte(TestYaml3))
	resource3, err := New(r3)
	if err != nil {
		t.Error(err)
	}

	namespaced3, _ := resource3.Namespaced()
	if namespaced3 {
		t.Errorf("Kind 'StorageClass' should not be namespaced.\n%+v\n", resource3)
	}

	var TestYaml4 = `
kind= StorageClass
metadata=
  name= Test
`
	r4 := bytes.NewReader([]byte(TestYaml4))
	_, err = New(r4)
	if err == nil {
		t.Errorf("Wrong formatted YAML should return an error\n%s\n", TestYaml4)
	}
}
