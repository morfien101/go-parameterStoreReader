package parameterstore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/yaml.v2"
)

type ParameterStore struct {
	ssmSession *ssm.SSM
	config     config
	format     *FormatOptions
}

type config struct {
	recursive    bool
	path         string
	decrypt      bool
	includePath  bool
	base64Values bool
}

func New(
	session *session.Session,
	path string,
	recursive bool,
	decrypt bool,
	includePath bool,
	b64Values bool,
	formatOptions *FormatOptions,
) *ParameterStore {
	ps := &ParameterStore{
		ssmSession: ssm.New(session),
		config: config{
			path:         path,
			recursive:    recursive,
			decrypt:      decrypt,
			includePath:  includePath,
			base64Values: b64Values,
		},
		format: formatOptions,
	}
	return ps
}

func (ps *ParameterStore) CollectPath(upperCase bool) (map[string]string, error) {
	inputObjects := ps.pathInput()
	values, err := ps.values(inputObjects)
	if err != nil {
		return nil, err
	}
	return values, err
}

func (ps *ParameterStore) CollectSingle() (string, error) {
	return ps.value(ps.singleInput())
}

// PathInput will create a input object for passing to Values
func (ps *ParameterStore) singleInput() *ssm.GetParameterInput {
	return &ssm.GetParameterInput{
		Name:           aws.String(ps.config.path),
		WithDecryption: aws.Bool(ps.config.decrypt),
	}
}

// PathInput will create a input object for passing to Values
func (ps *ParameterStore) pathInput() *ssm.GetParametersByPathInput {
	parameterInputPath := &ssm.GetParametersByPathInput{
		Path:           aws.String(ps.config.path),
		Recursive:      aws.Bool(ps.config.recursive),
		WithDecryption: aws.Bool(ps.config.decrypt),
	}

	return parameterInputPath
}

func (ps *ParameterStore) value(pip *ssm.GetParameterInput) (string, error) {
	output, err := ps.ssmSession.GetParameter(pip)
	if err != nil {
		return "", err
	}
	value := *output.Parameter.Value
	if ps.config.base64Values {
		value = ps.b64(value)
	}

	return value, nil
}

// Values will get keys and values from Parameter Store and return them as a map[string]string
func (ps *ParameterStore) values(pip *ssm.GetParametersByPathInput) (map[string]string, error) {
	out := make(map[string]string)

	for {
		// Try get the parameters
		output, err := ps.ssmSession.GetParametersByPath(pip)
		if err != nil {
			return out, fmt.Errorf("failed to get parameters. Error: %s", err)
		}

		if len(output.Parameters) == 0 {
			return nil, fmt.Errorf("no parameter found using %s", *pip.Path)
		}

		if ps.config.includePath {
			for _, parameter := range output.Parameters {
				out[*parameter.Name] = *parameter.Value
			}
		} else {
			// The splitter is used to show the values key minus what the user gave us.
			givenPath := ps.config.path
			if givenPath[len(givenPath)-1] != byte('/') {
				givenPath = givenPath + "/"
			}

			for _, parameter := range output.Parameters {
				// Break the parameter name in the given path and discovered path.
				// Use the discovered path in the out map
				splitPath := strings.Split(*parameter.Name, givenPath)
				discoveredPath := fmt.Sprintf("%s%s", "/", splitPath[len(splitPath)-1])
				out[discoveredPath] = *parameter.Value
			}
		}

		// Check to see if we need to go again
		if aws.StringValue(output.NextToken) != "" {
			// If we have a next token then we have more to collect.
			// Set next token and go again
			pip.NextToken = output.NextToken
		} else {
			// If nothing is there then break out the loop
			break
		}
	}

	// If base64 is requested then convert the values and store them in place of the current values.
	if ps.config.base64Values {
		for key, value := range out {
			out[key] = ps.b64(value)
		}
	}

	if ps.format.UpperCase || ps.format.Prefix != "" {
		formattedOutMap := map[string]string{}
		// Generate new keys
		for key, value := range out {
			newKey := key
			if ps.format.Prefix != "" {
				newKey = prefixPath(ps.format.Prefix, newKey)
			}
			if ps.format.UpperCase {
				newKey = strings.ToUpper(newKey)
			}
			if _, ok := formattedOutMap[newKey]; !ok {
				formattedOutMap[newKey] = value
			}
		}
		out = formattedOutMap
	}

	return out, nil
}

func prefixPath(prefix string, path string) string {
	sep := "/"
	splitPath := strings.Split(path, sep)
	splitPath[len(splitPath)-1] = fmt.Sprintf("%s%s", prefix, splitPath[len(splitPath)-1])
	return strings.Join(splitPath, sep)
}

func (ps *ParameterStore) b64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func (ps ParameterStore) FormatOutput(data map[string]string) ([]byte, error) {
	switch ps.format.Format {
	case "json":
		return json.Marshal(convertTree(data))
	case "pretty-json":
		return json.MarshalIndent(convertTree(data), "", "  ")
	case "yaml":
		return yaml.Marshal(convertTree(data))
	case "line":
		return lineFormat(data), nil
	case "env":
		return envFormat(data), nil
	default:
		return []byte{}, fmt.Errorf("output format '%s' is not valid", globalformatOptions.Format)
	}
}
