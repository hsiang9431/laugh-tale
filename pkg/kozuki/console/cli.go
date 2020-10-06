package console

import "laugh-tale/pkg/common/cli"

var appCmd = cli.Command{Name: "kozuki", DispStr: "./kozuki <cmd> args", SubCmdDesc: "Avaliable Commands:",
	SubCmd: []cli.Command{helpCmd, helpCmd1, helpCmd2, versionCmd}}

var helpCmd = cli.Command{Name: "help", DispStr: "help|-h|--help", Description: "Show this help message"}
var helpCmd1 = cli.Command{Name: "-h"}
var helpCmd2 = cli.Command{Name: "--help"}

var versionCmd = cli.Command{Name: "version", DispStr: "version", Description: "Show current version", RunFunc: version}

var (
	filenameFlag = cli.Flag{Name: "filename",
		DispStr:     "[--filename=<key id>]",
		Description: "The yaml file for running command",
		Required:    true}
	jsonFlag = cli.Flag{Name: "json",
		DispStr:     "[--json]",
		Description: "Output json instead of yaml",
		Required:    false}
	verboseFlag = cli.Flag{Name: "verbose",
		DispStr:     "[--verbose]",
		Description: "Print details to console",
		Required:    false}
)

var (
	createCmd = cli.Command{Name: "create",
		DispStr:     "create  <--arg=val>",
		Description: "Create object in database",
		Flags:       []cli.Flag{filenameFlag, jsonFlag, verboseFlag},
		RunFunc:     create}

	retrieveCmd = cli.Command{Name: "retrieve",
		DispStr:     "retrieve  <--arg=val>",
		Description: "Retrieve object from database",
		Flags:       []cli.Flag{filenameFlag, jsonFlag, verboseFlag},
		RunFunc:     retrieve}

	updateCmd = cli.Command{Name: "update",
		DispStr:     "update  <--arg=val>",
		Description: "Update object in database",
		Flags:       []cli.Flag{filenameFlag, jsonFlag, verboseFlag},
		RunFunc:     update}

	deleteCmd = cli.Command{Name: "delete",
		DispStr:     "delete  <--arg=val>",
		Description: "Delete object in database",
		Flags:       []cli.Flag{filenameFlag, jsonFlag, verboseFlag},
		RunFunc:     delete}
)

var keyCommand = cli.Command{Name: "key",
	DispStr:     "kozuki key <crud> <--arg=val>",
	Description: "Generate all required files for building encrypted image",
	SubCmdDesc:  "Supported CRUD actions:",
	SubCmd:      []cli.Command{createCmd, retrieveCmd, updateCmd, deleteCmd}}
