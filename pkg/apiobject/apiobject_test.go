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
