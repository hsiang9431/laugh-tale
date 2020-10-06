package console

import (
	"fmt"
	"io/ioutil"
	"laugh-tale/pkg/common/cli"

	"github.com/pkg/errors"
)

var (
	ServerBaseURL = "https://127.0.0.1:5566/kozuki/v1"

	CACertName  = "ca.crt"
	TLSCertName = "tls.crt"
	TLSKeyName  = "tls.pem"

	ConsoleWriter = ioutil.Discard
	ErrorWriter   = ioutil.Discard
)

var (
	Version   = "0.0.0"
	Build     = "0000000"
	BuildTime = "Sun Mar 1 15:04:05 PDT 2020"
)

func Run(args []string) error {
	cmd, ctx, err := cli.ParseCmd(&appCmd, args)
	if err != nil {
		cli.PrintMisuse(cmd, ErrorWriter, err, nil)
		return errors.Wrap(err, "Invalid command")
	}
	if cmd.Name == "help" ||
		cmd.Name == "-h" ||
		cmd.Name == "--help" {
		cli.PrintHelp(&appCmd, ConsoleWriter, nil)
		return nil
	}
	ctx, err = cli.ParseCliFlags(cmd, ctx)
	if err != nil {
		cli.PrintMisuse(cmd, ErrorWriter, err, nil)
		return errors.Wrap(err, "Invalid command line flags")
	}
	return cli.Call(cmd, ctx)
}

func version(ctx *cli.Context) error {
	fmt.Println("Kozuki, key management service:")
	fmt.Println("    console crud tool for management")
	fmt.Println("Version: ", Version)
	fmt.Println("Build: ", Build)
	fmt.Println("BuildTime: ", BuildTime)
	return nil
}
