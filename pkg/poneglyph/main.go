package poneglyph

import (
	"fmt"
	"io/ioutil"
	"laugh-tale/pkg/common/cli"

	"github.com/pkg/errors"
)

var (
	ServerURL    = ""
	ServerCACert = "ca.crt"

	WorkDir            = "/secret-container"
	EntrypointFileName = "entrypoint.sh"
	PayloadDirName     = "payload"

	ConsoleWriter = ioutil.Discard
	ErrorWriter   = ioutil.Discard

	Version   = "unknown"
	Build     = "unknown"
	BuildTime = "unknown"
)

func Run(args []string) error {
	cmd, ctx, err := cli.ParseCmd(&appCmd, args)
	if err != nil {
		cli.PrintMisuse(cmd, ErrorWriter, err, nil)
		return errors.Wrap(err, "Invalid command")
	}
	ctx, err = cli.ParseCliFlags(cmd, ctx)
	if err != nil {
		cli.PrintMisuse(cmd, ErrorWriter, err, nil)
		return errors.Wrap(err, "Invalid command line flags")
	}
	return cli.Call(cmd, ctx)
}

func version(ctx *cli.Context) error {
	fmt.Println("Poneglyph, encrypted container implant")
	fmt.Println("Version: ", Version)
	fmt.Println("Build: ", Build)
	fmt.Println("BuildTime: ", BuildTime)
	return nil
}
