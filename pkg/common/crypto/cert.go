package crypto

import (
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetCertPoolFromPath(path string, sysCert bool) (*x509.CertPool, error) {
	var err error
	var ret *x509.CertPool
	if _, err = os.Stat(path); err != nil {
		return nil, errors.Wrap(err, "Failed to traverse given path")
	}
	if sysCert {
		if ret, err = x509.SystemCertPool(); err != nil {
			return nil, errors.Wrap(err, "Failed to get system cert pool")
		}
	} else {
		ret = &x509.CertPool{}
	}
	if err := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Error traversing trusted CA path")
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		certPem, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.Wrapf(err, "Failed to read cert %s", file)
		}
		if !ret.AppendCertsFromPEM(certPem) {
			return errors.New("Failed to append server ca cert")
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "Failed to load all certs in path")
	}
	return ret, nil
}
