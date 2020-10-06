package poneglyph

import (
	"bytes"
	"fmt"
	"laugh-tale/pkg/common/cli"
	"laugh-tale/pkg/common/crypto"
	"laugh-tale/pkg/common/gzip"
	"laugh-tale/pkg/common/tar"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

func run(ctx *cli.Context) error {
	decKey, err := decKeyFlag.GetString(ctx)
	if err != nil {
		return err
	}
	decSrc := filepath.Join(WorkDir, encPrefix+PayloadDirName)
	decDst := filepath.Join(WorkDir, PayloadDirName+".tar.gz")
	if err := crypto.AESDecFile(decSrc, decDst, decKey); err != nil {
		return errors.Wrap(err, "Failed to decrypt payload")
	}
	unzipSrc := decDst
	unzipDst := filepath.Join(WorkDir, PayloadDirName+".tar")
	if err := gzip.Decompress(unzipSrc, unzipDst); err != nil {
		return errors.Wrap(err, "Failed to unzip payload")
	}
	untarSrc := unzipDst
	untarDst := filepath.Join(WorkDir, PayloadDirName)
	if err := tar.UnTarFileToDir(untarSrc, untarDst); err != nil {
		return errors.Wrap(err, "Failed to untar payload")
	}
	os.Remove(unzipSrc)
	os.Remove(untarSrc)
	// execute entrypoint
	entry := filepath.Join(WorkDir, EntrypointFileName)
	cmd := exec.Command(shellName, entry)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Entry point script returned with error:\n"+out.String())
	}
	fmt.Println("Payload terminated successfully")
	fmt.Println(out.String())
	return os.RemoveAll(WorkDir)
}
