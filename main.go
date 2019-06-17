package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/morfien101/go-parameterStoreReader/awsSession"
	parameterstore "github.com/morfien101/go-parameterStoreReader/parameterStore"
)

var (
	// VERSION is the application version
	VERSION = "0.0.1"

	// Flags
	flagPath       = flag.String("path", "", "Parameter Store path")
	flagBase64     = flag.Bool("base64", false, "Base64 encode collected values.")
	flagJSONMulti  = flag.Bool("json-output", false, "Output multiple key, values as JSON")
	flagJSONPretty = flag.Bool("json-pretty", false, "JSON output is displayed in pretty print")
	flagRecursive  = flag.Bool("recursive", false, "Look up all keys in branch")
	flagDecrypt    = flag.Bool("decrypt", false, "Request decrypted keys")
	flagAccessKey  = flag.String("access-key", "", "Access key for AWS API")
	flagSecretKey  = flag.String("secret-key", "", "Secret key for AWS API")
	flagRegion     = flag.String("region", "", "Region for AWS API")
	flagHelp       = flag.Bool("h", false, "Help menu")
	flagVersion    = flag.Bool("v", false, "Show Version")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println(VERSION)
		return
	}

	if *flagHelp {
		flag.PrintDefaults()
		return
	}

	if *flagPath == "" {
		log.Fatal("--path can not be empty")
	}

	if err := awsSession.SetAccessKey(*flagAccessKey); err != nil {
		log.Fatal("Failed to set environment variable AWS_ACCESS_KEY_ID for access to AWS")
	}
	if err := awsSession.SetSecretKey(*flagSecretKey); err != nil {
		log.Fatal("Failed to set environment variable AWS_SECRET_ACCESS_KEY for access to AWS")
	}
	if err := awsSession.SetRegion(*flagRegion); err != nil {
		log.Fatal("Failed to set environment variable AWS_REGION for access to AWS")
	}

	session, err := awsSession.New()
	if err != nil {
		log.Fatalf("Failed to create AWS Session. Error: %s", err)
	}

	ps := parameterstore.New(session, *flagPath, *flagRecursive, *flagDecrypt)
	if *flagRecursive {
		values, err := ps.CollectPath()
		if err != nil {
			log.Fatalf("Failed to read from parameter store. Error: %s", err)
		}

		// If base64 is requested then convert the values and store them in place of the current values.
		if *flagBase64 {
			for key, value := range values {
				values[key] = b64(value)
			}
		}
		switch {
		case *flagJSONMulti:
			var out []byte
			var err error
			if *flagJSONPretty {
				out, err = json.MarshalIndent(values, "", "  ")
			} else {
				out, err = json.Marshal(values)
			}

			if err != nil {
				log.Fatalf("Failed to convert to JSON. Error: %s", err)
			}
			fmt.Println(string(out))
		default:
			for key, value := range values {
				fmt.Printf("%s:%s\n", key, value)
			}
		}

	} else {
		value, err := ps.CollectSingle()
		if err != nil {
			log.Fatalf("Failed to read from parameter store. Error: %s", err)
		}
		if *flagBase64 {
			value = b64(value)
		}
		fmt.Println(value)
	}
}

func b64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}
