package parameterstore

import (
	"testing"
	"regexp"

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

func TestLineFormat(t *testing.T) {
	in := map[string]string{
		"/key1/key2/final1": "value1",
	}

	out := string(lineFormat(in))
	expected := "/key1/key2/final1:value1"
	if out != expected {
		t.Logf("lineformat did not give the expected result. Want: %s, Got: %s", expected, out)
		t.Fail()
	}
}

func TestEnvFormat(t *testing.T) {
	in := map[string]string{
		"/key1/key2/FINAL1": "value1",
		"/key1/key2/FINAL2": "value2",
	}

	out := string(envFormat(in))
	expected := "FINAL[1-2]=value[1-2]\nFINAL[1-2]=value[1-2]"
	re := regexp.MustCompile(expected)
	if !re.MatchString(out) {
		t.Logf("format 'env' did not give the expected result. Want something like: %s, Got: %s", expected, out)
		t.Fail()
	}
}