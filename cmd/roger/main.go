package main

import (
	"fmt"
	"os"

	"laugh-tale/pkg/roger"
)

// Kubernetes secrets can be distributed
// either in mounted volumes or environment variables
const tlsKeyPassEnv = "SECRET_TLS_KEY_PASS_B64"
const signingKeyPassEnv = "SECRET_SIGNING_KEY_PASS_B64"

var appName = "roger"

var (
	Version   = "0.0.0"
	Build     = "0000000"
	BuildTime = "Sun Mar 1 15:04:05 PDT 2020"
)

func main() {
	// mounted volumes
	roger.InputDirPath = "/input"
	roger.OutputDirPath = "/output"
	// mounted volume at /home/roger
	roger.SecretDirPath = "/home/roger/secret"
	roger.CACertPath = "/home/roger/trusted-ca"
	// implant location and filenames
	roger.WorkDirPath = "/work"
	roger.KeyRetriverCACertName = "key-retriever-ca.crt"

	roger.ConsoleWriter = os.Stdout
	roger.ErrorWriter = os.Stderr
	roger.Version = Version

	roger.TLSKeyPass = os.Getenv(tlsKeyPassEnv)
	roger.SignKeyPass = os.Getenv(signingKeyPassEnv)

	roger.Version = Version
	roger.Build = Build
	roger.BuildTime = BuildTime

	if err := roger.Run(os.Args); err != nil {
		fmt.Printf("Failed to run application \"%s\": %s\n", appName, err.Error())
		os.Exit(1)
	}
}
