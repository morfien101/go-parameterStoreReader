package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	formats = []string{"line", "json", "pretty-json", "yaml"}
)

func formatOutput(data map[string]string, format string, createTree bool) ([]byte, error) {
	switch format {
	case "json":
		return json.Marshal(convertTree(data, createTree))
	case "pretty-json":
		return json.MarshalIndent(convertTree(data, createTree), "", "  ")
	case "yaml":
		return yaml.Marshal(convertTree(data, createTree))
	case "line":
		return lineFormat(data), nil
	default:
		return []byte{}, fmt.Errorf("output format '%s' is not valid", format)
	}
}

func formatValidation(formatString string) bool {
	valid := false

	for _, formatValue := range formats {
		if formatString == formatValue {
			valid = true
		}
	}

	return valid
}

func lineFormat(data map[string]string) []byte {
	output := []string{}
	for key, value := range data {
		output = append(output, fmt.Sprintf("%s:%s", key, value))
	}

	return []byte(strings.Join(output, "\n"))
}

func convertTree(data map[string]string, split bool) map[string]interface{} {
	treeView := make(map[string]interface{})

	if split {
		for keySolid, value := range data {
			keys := strings.Split(keySolid, "/")
			// clean keys
			cleanKeys := []string{}
			for _, key := range keys {
				if key != "" {
					cleanKeys = append(cleanKeys, key)
				}
			}
			keys = []string{}

			mapper(treeView, cleanKeys, value)
		}
	} else {
		for key, value := range data {
			treeView[key] = value
		}
	}

	return treeView

}

func mapper(current map[string]interface{}, keys []string, value string) {
	if len(keys) == 1 {
		current[keys[0]] = value
		return
	}

	if _, ok := current[keys[0]]; !ok {
		current[keys[0]] = make(map[string]interface{})
	}

	mapper(current[keys[0]].(map[string]interface{}), keys[1:], value)
}
