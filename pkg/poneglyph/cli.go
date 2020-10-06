package poneglyph

import "laugh-tale/pkg/common/cli"

var appCmd = cli.Command{Name: "poneglyph",
	SubCmd: []cli.Command{versionCmd, startCmd, runCmd}}

var versionCmd = cli.Command{Name: "version", RunFunc: version}
var startCmd = cli.Command{Name: "start", RunFunc: start}

var (
	decKeyFlag = cli.Flag{Name: "dec-key-b64",
		Required: true}

	runCmd = cli.Command{Name: "run", RunFunc: run,
		Flags: []cli.Flag{decKeyFlag}}
)
