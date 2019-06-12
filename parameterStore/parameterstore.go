package parameterstore

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ParameterStore struct {
	ssmSession *ssm.SSM
	config     config
}

type config struct {
	recursive bool
	path      string
	decrypt   bool
}

func New(session *session.Session, path string, recursive, decrypt bool) *ParameterStore {
	ps := &ParameterStore{
		ssmSession: ssm.New(session),
		config: config{
			path:      path,
			recursive: recursive,
			decrypt:   decrypt,
		},
	}
	return ps
}

func (ps *ParameterStore) CollectPath() (map[string]string, error) {
	inputObjects := ps.pathInput("")
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
func (ps *ParameterStore) pathInput(nextToken string) *ssm.GetParametersByPathInput {
	parameterInputPath := &ssm.GetParametersByPathInput{
		Path:           aws.String(ps.config.path),
		Recursive:      aws.Bool(ps.config.recursive),
		WithDecryption: aws.Bool(ps.config.decrypt),
	}

	if len(nextToken) < 0 {
		parameterInputPath.NextToken = aws.String(nextToken)
	}

	return parameterInputPath
}

func (ps *ParameterStore) value(pip *ssm.GetParameterInput) (string, error) {
	output, err := ps.ssmSession.GetParameter(pip)
	if err != nil {
		return "", err
	}
	return *output.Parameter.Value, nil
}

// Values will get keys and values from Parameter Store and return them as a map[string]string
func (ps *ParameterStore) values(pip *ssm.GetParametersByPathInput) (map[string]string, error) {
	out := make(map[string]string)

	for {
		// Try get the parameters
		output, err := ps.ssmSession.GetParametersByPath(pip)
		if err != nil {
			return out, fmt.Errorf("Failed to get parameters. Error: %s", err)
		}

		if len(output.Parameters) == 0 {
			return nil, fmt.Errorf("No Secrets found using %s", *pip.Path)
		}

		// Extract the keys and values
		for _, parameter := range output.Parameters {
			SplitPath := strings.Split(*parameter.Name, "/")
			out[SplitPath[len(SplitPath)-1]] = *parameter.Value
		}

		// Check to see if we need to go again
		if output.NextToken != nil {
			// If we have a next token then we have more to collect.
			// Set next token and go again
			pip.NextToken = aws.String(*output.NextToken)
		} else {
			// If nothing is there then break out the loop
			break
		}
	}

	return out, nil
}
