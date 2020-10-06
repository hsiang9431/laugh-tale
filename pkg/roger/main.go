package roger

import (
	"fmt"
	"io/ioutil"
	"laugh-tale/pkg/common/cli"
	"sync"

	"github.com/pkg/errors"
)

var (
	InputDirPath  = "/input"
	OutputDirPath = "/output"
	// implant is in this path
	WorkDirPath = "/work"

	SecretDirPath = "/home/roger/secret"
	CACertPath    = "/home/roger/trusted-ca"

	// tls key passphrase
	TLSKeyPass string

	// signing key passphrase
	SignKeyPass string

	// file under CACertPath
	KeyRetriverCACertName = "ca.crt"

	// files under SecretDirPath
	TLSCertName = "tls.crt"
	TLSKeyName  = "tls.pem"

	// files generated to output
	PayloadDirName     = "payload"
	EntrypointFileName = "entrypoint.sh"

	Version   = "unknown"
	Build     = "unknown"
	BuildTime = "unknown"

	ConsoleWriter = ioutil.Discard
	ErrorWriter   = ioutil.Discard
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
	fmt.Println("Roger, encrypted image prepare agent")
	fmt.Println("Version: ", Version)
	fmt.Println("Build: ", Build)
	fmt.Println("BuildTime: ", BuildTime)
	return nil
}

var logLock = sync.Mutex{}
var errLock = sync.Mutex{}
var verbose = false

func logInfo(msg string) {
	if verbose {
		logLock.Lock()
		fmt.Fprintf(ConsoleWriter, "[INFO] %s", msg)
		logLock.Unlock()
	}
}

func logWarning(msg string) {
	logLock.Lock()
	fmt.Println(ConsoleWriter, "[WARN] %s", msg)
	logLock.Unlock()
}

func logError(msg string) {
	errLock.Lock()
	fmt.Println(ErrorWriter, "[ERRO] %s", msg)
	errLock.Unlock()
}
