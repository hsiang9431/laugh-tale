package roger

import (
	"io/ioutil"
	"laugh-tale/pkg/common/crypto"
	"laugh-tale/pkg/common/gzip"
	"laugh-tale/pkg/common/tar"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type tarAndEncTask struct {
	PayloadPassword string
	done            bool
}

type genSignedKeyCertPair struct {
	done bool
}

type cpFilesTask struct {
	done bool
}

func (t *tarAndEncTask) Init() error    { t.done = false; return nil }
func (t *tarAndEncTask) CleanUp() error { return nil }
func (t *tarAndEncTask) Done() bool     { return t.done }
func (t *tarAndEncTask) Run() error {
	tarSrc := filepath.Join(InputDirPath, PayloadDirName)
	tarDst := filepath.Join(WorkDirPath, PayloadDirName+".tar")
	if err := tar.TarDirToFile(tarSrc, tarDst); err != nil {
		logError("Failed to create tar file for payload")
		return errors.Wrap(err, "Failed to create tar file for payload")
	}
	logInfo("Payload tar created")
	zipSrc := tarDst
	zipDst := filepath.Join(WorkDirPath, PayloadDirName+".tar.gz")
	if err := gzip.Compress(zipSrc, zipDst); err != nil {
		logError("Failed to zip payload tar file")
		return errors.Wrap(err, "Failed to zip payload tar file")
	}
	encSrc := zipDst
	encDst := filepath.Join(OutputDirPath, encPrefix+PayloadDirName)
	err := crypto.AESEncFile(encSrc, encDst, t.PayloadPassword)
	if err != nil {
		logError("Failed to encrypt payload tar file")
		return errors.Wrap(err, "Failed to encrypt payload tar file")
	}
	logInfo("Payload tar encrypted")
	t.done = true
	return nil
}

func (t *cpFilesTask) Init() error    { t.done = false; return nil }
func (t *cpFilesTask) CleanUp() error { return nil }
func (t *cpFilesTask) Done() bool     { return t.done }
func (t *cpFilesTask) Run() error {
	epSrc := filepath.Join(InputDirPath, EntrypointFileName)
	epDst := filepath.Join(OutputDirPath, EntrypointFileName)
	if err := cp(epSrc, epDst, 0755); err != nil {
		return err
	}
	logInfo("Entry point script copied")
	bobSrc := filepath.Join(WorkDirPath, implantBinary)
	bobDst := filepath.Join(OutputDirPath, implantBinary)
	if err := cp(bobSrc, bobDst, 0755); err != nil {
		return err
	}
	bobShSrc := filepath.Join(WorkDirPath, implantScript)
	bobShDst := filepath.Join(OutputDirPath, implantScript)
	if err := cp(bobShSrc, bobShDst, 0755); err != nil {
		return err
	}
	gwenCertSrc := filepath.Join(CACertPath, KeyRetriverCACertName)
	gwenCertDst := filepath.Join(OutputDirPath, KeyRetriverCACertName)
	if err := cp(gwenCertSrc, gwenCertDst, 0644); err != nil {
		return err
	}
	t.done = true
	logInfo("Implant files copied")
	return nil
}

func cp(src string, dst string, mode os.FileMode) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.Wrap(err, "Error copying file, can not read file")
	}
	err = ioutil.WriteFile(dst, input, mode)
	if err != nil {
		return errors.Wrap(err, "Error copying file, can not create file")
	}
	return nil
}
