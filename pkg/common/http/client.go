package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"laugh-tale/pkg/common/crypto"
	"net/http"

	"github.com/pkg/errors"
)

func HTTPClientNoTLS() *http.Client {
	return &http.Client{}
}

// ref: https://gist.github.com/michaljemala/d6f4e01c4834bf47a9c4
// ref: https://stackoverflow.com/questions/56129533/tls-with-certificate-private-key-and-pass-phrase

func HTTPClientTLSFromPem(caCert string, tlsCert, tlsKey []byte, sysCert bool) (*http.Client, error) {
	x509Cert, err := tls.X509KeyPair(tlsCert, tlsKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load tls cert or key")
	}
	caCertPool, err := crypto.GetCertPoolFromPath(caCert, sysCert)
	if err != nil {
		return nil, err
	}
	return prepTLSConfig(x509Cert, caCertPool), nil
}

func HTTPClientTLSFromFile(caCert, tlsCert, tlsKey string, sysCert bool) (*http.Client, error) {
	x509Cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load tls cert or key")
	}
	caCertPool, err := crypto.GetCertPoolFromPath(caCert, sysCert)
	if err != nil {
		return nil, err
	}
	return prepTLSConfig(x509Cert, caCertPool), nil
}

func HTTPClientTLSFromFilePassphrase(caCert, tlsCert, tlsKey string, keyPass []byte, sysCert bool) (*http.Client, error) {
	encKeyPem, err := ioutil.ReadFile(tlsKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load tls key file")
	}
	encKeyPemBlock, _ := pem.Decode(encKeyPem)
	if !x509.IsEncryptedPEMBlock(encKeyPemBlock) {
		return HTTPClientTLSFromFile(caCert, tlsCert, tlsKey, sysCert)
	}
	keyDer, err := x509.DecryptPEMBlock(encKeyPemBlock, keyPass)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decrypt tls key")
	}
	encKeyPemBlock.Bytes = keyDer
	encKeyPemBlock.Headers = nil

	keyPem := pem.EncodeToMemory(encKeyPemBlock)

	certPem, err := ioutil.ReadFile(tlsCert)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load tls cert")
	}
	x509Cert, err := tls.X509KeyPair(certPem, keyPem)

	caCertPool, err := crypto.GetCertPoolFromPath(caCert, sysCert)
	if err != nil {
		return nil, err
	}
	return prepTLSConfig(x509Cert, caCertPool), nil
}

func prepTLSConfig(cert tls.Certificate, certPool *x509.CertPool) *http.Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}
