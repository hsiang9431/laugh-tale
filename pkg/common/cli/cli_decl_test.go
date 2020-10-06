package cli_test

import "laugh-tale/pkg/common/cli"

// output of appCmd.PrintHelp() should be:
// Usage
//     ./app <cmd> args

// Avaliable Commands:
//     help|-h|--help      Show this help message
//     setup <task>        Run setup task
//     start               Start the app
//     stop                Stop the app
//     status              Show status of the app
//     uninstall           Uninstall app
//     version             Show version of the app

// Avaliable Tasks for setup:
//     setup task1 [--arg=val]     Description 1
//         --arg11=<val>   arg1 description <Required>
//         --arg12=<val>   arg2 description
//         --arg13=<val>   arg3 description
//     setup task2 [--arg=val]     Description 2
//         --arg21=<val>   arg1 description <Required>
//         --arg22=<val>   arg2 description
//     setup task3                 Description 3
//     setup all                   Run all setups, arguments should be defined in env

var printSettings = cli.PrintSettings{
	HelpStr:          "Usage",
	UnknownCmdErrStr: "Invalid usage of command: ",

	DefaultIndent:    4,
	DefaultTabWidth:  8,
	PrintRequiredStr: "<Required>",
}

var appCmd = cli.Command{Name: "app", DispStr: "./app <cmd> args", SubCmdDesc: "Avaliable Commands:",
	SubCmd: []cli.Command{helpShow, help1, help2, run, setup, start, stop, status, uninstall, version}}

var helpShow = cli.Command{Name: "help", DispStr: "help|-h|--help", Description: "Show this help message"}
var help1 = cli.Command{Name: "-h"}
var help2 = cli.Command{Name: "--help"}

var run = cli.Command{Name: "run"}

var (
	start     = cli.Command{Name: "start", DispStr: "start", Description: "Start the app"}
	stop      = cli.Command{Name: "stop", DispStr: "stop", Description: "Stop the app"}
	status    = cli.Command{Name: "status", DispStr: "status", Description: "Show status of the app"}
	uninstall = cli.Command{Name: "uninstall", DispStr: "uninstall", Description: "Uninstall app"}
	version   = cli.Command{Name: "version", DispStr: "version", Description: "Show version of the app"}

	setup = cli.Command{Name: "setup",
		DispStr:     "setup <task>",
		Description: "Run setup task",
		SubCmdDesc:  "Avaliable Tasks for setup:",
		SubCmd:      []cli.Command{task1, task2, task3, all}}
)

var (
	task1Fg1 = cli.Flag{Name: "arg11",
		DispStr:     "--arg11=<val>",
		Description: "arg1 description",
		Required:    true}
	task1Fg2 = cli.Flag{Name: "arg12",
		DispStr:     "--arg12=<val>",
		Description: "arg2 description",
		Required:    false}
	task1Fg3 = cli.Flag{Name: "arg13",
		DispStr:     "--arg13=<val>",
		Description: "arg3 description",
		Required:    false}

	task1 = cli.Command{Name: "task1",
		DispStr:     "setup task1 [--arg=val]",
		Description: "Description 1",
		Flags:       []cli.Flag{task1Fg1, task1Fg2, task1Fg3}}
)

var (
	task2Fg1 = cli.Flag{Name: "arg21",
		DispStr:     "--arg21=<val>",
		Description: "arg1 description",
		Required:    true}

	task2Fg2 = cli.Flag{Name: "arg22",
		DispStr:     "--arg22=<val>",
		Description: "arg2 description",
		Required:    false}

	task2 = cli.Command{Name: "task2",
		DispStr:     "setup task2 [--arg=val]",
		Description: "Description 2",
		Flags:       []cli.Flag{task2Fg1, task2Fg2}}
)

// adding '\t' at the end of shorter command display strings
// can helps aligning description strings.
// example without '\t':
// Avaliable Tasks for setup:
//     setup task1 [--arg=val]     Description 1
//         --arg11=<val>   arg1 description <Required>
//         --arg12=<val>   arg2 description
//         --arg13=<val>   arg3 description
//     setup task2 [--arg=val]     Description 2
//         --arg21=<val>   arg1 description <Required>
//         --arg22=<val>   arg2 description
//     setup task3         Description 3
//     setup all           Run all setups, arguments should be defined in env
//
// example with '\t':
// Avaliable Tasks for setup:
//     setup task1 [--arg=val]     Description 1
//         --arg11=<val>   arg1 description <Required>
//         --arg12=<val>   arg2 description
//         --arg13=<val>   arg3 description
//     setup task2 [--arg=val]     Description 2
//         --arg21=<val>   arg1 description <Required>
//         --arg22=<val>   arg2 description
//     setup task3                 Description 3
//     setup all                   Run all setups, arguments should be defined in env
var (
	task3 = cli.Command{Name: "task3", DispStr: "setup task3\t", Description: "Description 3"}

	all = cli.Command{Name: "all", DispStr: "setup all\t",
		Description: "Run all setups, arguments should be defined in env"}
)
