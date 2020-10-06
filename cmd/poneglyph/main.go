package main

import (
	"fmt"
	"os"

	"laugh-tale/pkg/poneglyph"
)

var appName = "poneglyph"

var (
	Version   = "0.0.0"
	Build     = "0000000"
	BuildTime = "Sun Mar 1 15:04:05 PDT 2020"
)

func main() {
	poneglyph.ServerCACert = "key-retriever-ca.crt"

	poneglyph.ConsoleWriter = os.Stdout
	poneglyph.ErrorWriter = os.Stderr

	poneglyph.Version = Version
	poneglyph.Build = Build
	poneglyph.BuildTime = BuildTime

	if err := poneglyph.Run(os.Args); err != nil {
		fmt.Printf("Failed to run application \"%s\": %s", appName, err.Error())
		os.Exit(1)
	}
}
