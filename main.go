package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/morfien101/go-parameterStoreReader/awsSession"
	parameterstore "github.com/morfien101/go-parameterStoreReader/parameterStore"
)

var (
	// VERSION is the application version
	VERSION = "0.0.1"
)

func main() {
	// Flags
	flagPath := flag.String("path", "", "Parameter Store path.")
	flagBase64 := flag.Bool("base64", false, "Base64 encode collected values.")
	flagFormat := flag.String("format", "line", fmt.Sprintf("Format for output. Supported values: %s.", strings.Join(parameterstore.ValidFormats(), ",")))
	flagRecursive := flag.Bool("recursive", false, "Look up all keys in branch.")
	flagDecrypt := flag.Bool("decrypt", false, "Request decrypted keys.")
	flagAccessKey := flag.String("access-key", "", "Access key for AWS API.")
	flagSecretKey := flag.String("secret-key", "", "Secret key for AWS API.")
	flagProfile := flag.String("profile", "", "AWS Profile to use.")
	flagCredsFile := flag.String("config-file", "", "AWS Config file override, only valid with -profile.")
	flagRegion := flag.String("region", "", "Region for AWS API.")
	flagIncludePath := flag.Bool("include-path", false, "Include the passed in path in the output. Only used with recursive lookups.")
	flagFileOutput := flag.String("f", "", "Output to specified file.")
	flagHelp := flag.Bool("h", false, "Help menu.")
	flagVersion := flag.Bool("v", false, "Show application Version.")
	flag.Parse()

	if *flagVersion {
		fmt.Println(VERSION)
		return
	}

	if *flagHelp {
		flag.PrintDefaults()
		return
	}

	if !parameterstore.FormatValidation(*flagFormat) {
		fmt.Printf("Format %s is not valid.\n", *flagFormat)
		os.Exit(1)
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

	var session *session.Session
	var err error
	if *flagProfile == "" {
		session, err = awsSession.New()
	} else {
		session, err = awsSession.NewWithOptions(*flagProfile, *flagCredsFile)
	}
	if err != nil {
		log.Fatalf("Failed to create AWS Session. Error: %s", err)
	}

	ps := parameterstore.New(session, *flagPath, *flagRecursive, *flagDecrypt, *flagIncludePath, *flagBase64)

	var output string

	if *flagRecursive {
		psMap, err := ps.CollectPath()
		if err != nil {
			log.Fatalf("Failed to read from parameter store. Error: %s", err)
		}

		formattedOutput, err := ps.FormatOutput(psMap, *flagFormat)
		if err != nil {
			log.Fatal(err)
		}
		output = string(formattedOutput)
	} else {
		output, err = ps.CollectSingle()
		if err != nil {
			log.Fatalf("Failed to read from parameter store. Error: %s\n", err)
		}
	}

	if *flagFileOutput != "" {
		err := ioutil.WriteFile(*flagFileOutput, []byte(output), 0644)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	fmt.Println(output)
}
