package apiobject_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/300481/kitops/pkg/apiobject"
)

var TestYaml = `
---
kind: Test
metadata:
  name: Test
  namespace: Test
`

var TestObject = apiobject.ApiObject{
	Kind: "Test",
	Metadata: struct {
		Name      string
		Namespace string
	}{
		Name:      "Test",
		Namespace: "Test",
	},
}

func TestNew(t *testing.T) {
	r := bytes.NewReader([]byte(TestYaml))
	object, err := apiobject.New(r)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Testing 'TestNew()'\nTestObject:\n%+v\nNew Object:\n%+v\n", TestObject, *object)
	if *object != TestObject {
		t.Errorf("Test Objects don't match.\nTestObject:\n%+v\nNew Object:\n%+v\n", TestObject, *object)
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
	object1, err := apiobject.New(r1)
	if err != nil {
		t.Error(err)
	}

	_, err = object1.Namespaced()
	if err == nil {
		t.Errorf("Kind 'Test' should not be found.\n%+v\n", object1)
	}

	var TestYaml2 = `
---
kind: Deployment
metadata:
  name: Test
  namespace: Test
`
	r2 := bytes.NewReader([]byte(TestYaml2))
	object2, err := apiobject.New(r2)
	if err != nil {
		t.Error(err)
	}

	namespaced2, err := object2.Namespaced()
	if !namespaced2 {
		t.Errorf("Kind 'Deployment' should be namespaced.\n%+v\n", object2)
	}

	var TestYaml3 = `
---
kind: StorageClass
metadata:
  name: Test
  namespace: Test
`
	r3 := bytes.NewReader([]byte(TestYaml3))
	object3, err := apiobject.New(r3)
	if err != nil {
		t.Error(err)
	}

	namespaced3, err := object3.Namespaced()
	if namespaced3 {
		t.Errorf("Kind 'StorageClass' should not be namespaced.\n%+v\n", object3)
	}
}
