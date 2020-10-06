package roger

import "laugh-tale/pkg/common/cli"

var appCmd = cli.Command{Name: "roger", DispStr: "./roger <cmd> args", SubCmdDesc: "Avaliable Commands:",
	SubCmd: []cli.Command{helpCmd, helpCmd1, helpCmd2, versionCmd, encryptCmd}}

var helpCmd = cli.Command{Name: "help", DispStr: "help|-h|--help", Description: "Show this help message"}
var helpCmd1 = cli.Command{Name: "-h"}
var helpCmd2 = cli.Command{Name: "--help"}

var versionCmd = cli.Command{Name: "version", DispStr: "version", Description: "Show current version", RunFunc: version}

var (
	keyIDFlag = cli.Flag{Name: "key-id",
		DispStr:     "[--key-id=<key id>]",
		Description: "The uuid of key to bind",
		Required:    true}
	imageIDFlag = cli.Flag{Name: "image-id",
		DispStr:     "[--image-id=<image id>]",
		Description: "The docker image id of image to bind, in the format of sha256:...",
		Required:    true}
	imageFlag = cli.Flag{Name: "image",
		DispStr:     "[--image=<base image>]",
		Description: "The base image to reference in Dockerfile, uses ubuntu:latest image if not provided",
		Required:    false}
	volumeFlag = cli.Flag{Name: "volume",
		DispStr:     "[--volume=<allowed volume>]",
		Description: "Allowed volumes in encrypted container, use this flag multiple times for adding more volumes",
		Required:    false}
	keyServerFlag = cli.Flag{Name: "key-server",
		DispStr:     "[--key-server=<ker server address>]",
		Description: "The address of key server",
		Required:    true}
	verboseFlag = cli.Flag{Name: "verbose",
		DispStr:     "[--verbose]",
		Description: "Print detailed log to console",
		Required:    false}

	encryptCmd = cli.Command{Name: "encrypt",
		DispStr:     "roger encrypt <--arg=val>",
		Description: "Generate all required files for building encrypted image",
		Flags:       []cli.Flag{imageFlag, volumeFlag, keyServerFlag, verboseFlag},
		RunFunc:     encrypt}

	bindIDCmd = cli.Command{Name: "bind-image-id",
		DispStr:     "roger bind-image-id <--arg=val>",
		Description: "Bind the image ID of newly generated to tag",
		Flags:       []cli.Flag{keyIDFlag, imageIDFlag, keyServerFlag, verboseFlag},
		RunFunc:     bind}
)
