package main

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestMapper(t *testing.T) {
	output := make(map[string]interface{})
	matcher := "matcher"
	mapper(output, []string{"one", "two", "three", "four"}, "value")
	mapper(output, []string{"one", "a"}, matcher)

	s, err := yaml.Marshal(output)
	if err != nil {
		t.Logf("Failed to marshal mapper output. Error: %s", err)
		t.Fail()
	}
	if output["one"].(map[string]interface{})["a"].(string) != matcher {
		t.Logf("Didn't get the correct value from the new map")
		t.Fail()
	}
	t.Logf("\n%s", s)
}
