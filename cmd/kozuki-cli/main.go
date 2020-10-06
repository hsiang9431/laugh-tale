package main

import (
	"fmt"
	"laugh-tale/pkg/kozuki/console"
	"os"
)

var (
	Version   = "0.0.0"
	Build     = "0000000"
	BuildTime = "Sun Mar 1 15:04:05 PDT 2020"
)

var appName = "kozuki-cli"

func loadEnv() {
	console.ServerBaseURL = os.Getenv("KOZUKI_BASE_URL")
	console.CACertName = os.Getenv("KOZUKI_SERVER_CA_CERT")
	console.TLSCertName = os.Getenv("KOZUKI_CLI_TLS_CERT")
	console.TLSKeyName = os.Getenv("KOZUKI_CLI_TLS_KEY")

	console.Version = Version
	console.Build = Build
	console.BuildTime = BuildTime

	if err := console.Run(os.Args); err != nil {
		fmt.Printf("Failed to run application \"%s\": %s\n", appName, err.Error())
		os.Exit(1)
	}
}
