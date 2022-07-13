package runner

import (
	"testing"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("InitializeGinkgoParams should should set up some default parameters for ginkgo", func(t *testing.T) {
		defaultParams := InitializeGinkgoParams()
		assert.Equal(t, "", defaultParams["GinkgoTestPackage"])
		assert.Equal(t, "-r", defaultParams["GinkgoRecursive"])
		assert.Equal(t, "-p", defaultParams["GinkgoParallel"])
		assert.Equal(t, "--randomize-all", defaultParams["GinkgoRandomize"])
		assert.Equal(t, "--randomize-suites", defaultParams["GinkgoRandomizeSuites"])
		assert.Equal(t, "--trace", defaultParams["GinkgoTrace"])
		assert.Equal(t, "--junit-report report.xml", defaultParams["GinkgoJunitReport"])

	})

	t.Run("FindGoinkgoParams should override default params when provided with new value", func(t *testing.T) {
		defaultParams := InitializeGinkgoParams()
		variables := make(map[string]testkube.Variable)
		variable_one := testkube.Variable{
			Name:  "GinkgoTestPackage",
			Value: "e2e",
		}
		variable_two := testkube.Variable{
			Name:  "GinkgoRecursive",
			Value: "",
		}
		variables["GinkgoTestPackage"] = variable_one
		variables["GinkgoRecursive"] = variable_two
		execution := testkube.Execution{
			Variables: variables,
		}
		mappedParams := FindGinkgoParams(&execution, defaultParams)
		assert.Equal(t, "e2e", mappedParams["GinkgoTestPackage"])
		assert.Equal(t, "", mappedParams["GinkgoRecursive"])
	})

	t.Run("BuildGinkgoArgs should build ginkgo args slice", func(t *testing.T) {
		defaultParams := InitializeGinkgoParams()
		argSlice := BuildGinkgoArgs(defaultParams)
		assert.Contains(t, argSlice, "-r")
		assert.Contains(t, argSlice, "-p")
		assert.Contains(t, argSlice, "--randomize-all")
		assert.Contains(t, argSlice, "--randomize-suites")
		assert.Contains(t, argSlice, "--trace")
		assert.Contains(t, argSlice, "--junit-report")
		assert.Contains(t, argSlice, "report.xml")
	})

	t.Run("BuildGinkgoPassThroughFlags should build pass through flags slice from leftover Variables and from Args", func(t *testing.T) {
		variables := make(map[string]testkube.Variable)
		variable_one := testkube.Variable{
			Name:  "one",
			Value: "one",
		}
		variable_two := testkube.Variable{
			Name:  "two",
			Value: "two",
		}
		variables["GinkgoPassThroughOne"] = variable_one
		variables["GinkgoPassThroughTwo"] = variable_two

		args := []string{
			"--three",
			"--four=four",
		}

		execution := testkube.Execution{
			Variables: variables,
			Args:      args,
		}
		passThroughs := BuildGinkgoPassThroughFlags(execution)
		assert.Contains(t, passThroughs, "--")
		assert.Contains(t, passThroughs, "--one=one")
		assert.Contains(t, passThroughs, "--two=two")
		assert.Contains(t, passThroughs, "--three")
		assert.Contains(t, passThroughs, "--four=four")
	})
}
