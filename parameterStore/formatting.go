package parameterstore

import (
	"fmt"
	"strings"
)

var (
	formats             = []string{"line", "json", "pretty-json", "yaml", "env"}
	globalformatOptions = FormatOptions{}
)

type FormatOptions struct {
	Format    string
	Prefix    string
	UpperCase bool
}

func ValidFormats() []string {
	return formats
}

func FormatValidation(formatString string) bool {
	valid := false

	for _, formatValue := range formats {
		if formatString == formatValue {
			valid = true
		}
	}

	return valid
}

func upperIfNeeded(s string) string {
	if globalformatOptions.UpperCase {
		return strings.ToUpper(s)
	}
	return s
}

func prefixIfNeeded(s string) string {
	if globalformatOptions.Prefix != "" {
		return fmt.Sprintf("%s%s", globalformatOptions.Prefix, s)
	}
	return s
}

func lineFormat(data map[string]string) []byte {
	output := []string{}
	for key, value := range data {
		output = append(output, fmt.Sprintf("%s:%s", key, value))
	}

	return []byte(strings.Join(output, "\n"))
}

func envFormat(data map[string]string) []byte {
	out := []string{}

	for key, value := range data {
		brokenKeys := strings.Split(key, "/")
		out = append(out, fmt.Sprintf("%s=%s", upperIfNeeded(prefixIfNeeded(brokenKeys[len(brokenKeys)-1])), value))
	}

	return []byte(strings.Join(out, "\n"))
}

func convertTree(data map[string]string) map[string]interface{} {
	treeView := make(map[string]interface{})

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
