package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/pkg/errors"
)

// ref: https://gist.github.com/samuel/8b500ddd3f6118d052b5e6bc16bc4c09

var rsaBits = 2048
var defaultTimeGood = 180 * 24 * time.Hour

func GenSelfSignCert(certTemplate *x509.Certificate, prvPwd []byte) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to generate rsa key pair")
	}
	privateKeyPem, _, err := RSAKeyPairToPem(privateKey, prvPwd)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to convert private key to pem")
	}
	if certTemplate == nil {
		certTemplate = GetTemplateFromCSR(nil, nil)
	}
	signedCert, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, privateKey.Public, privateKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create self signed certificate")
	}
	return signedCert, privateKeyPem, nil
}

func Sign(certTemplate *x509.Certificate, cert4SignPem []byte, pub2Sign crypto.PublicKey, signKeyPrvPem []byte, prvPwd []byte) ([]byte, error) {
	prv4Sign, err := PemToPrivateKey(signKeyPrvPem, prvPwd)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse signer private key")
	}
	cert4Sign, err := x509.ParseCertificate(cert4SignPem)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse signer certificate")
	}
	if certTemplate == nil {
		certTemplate = GetTemplateFromCSR(nil, nil)
	}
	signedCert, err := x509.CreateCertificate(rand.Reader, certTemplate, cert4Sign, pub2Sign, prv4Sign)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create self signed certificate")
	}
	return signedCert, nil
}

func GetTemplateFromCSR(csr *x509.CertificateRequest, caCert *x509.Certificate) *x509.Certificate {
	var issuer pkix.Name
	if caCert != nil {
		issuer = caCert.Subject
	}
	// default template
	if csr == nil {
		return &x509.Certificate{
			SerialNumber: big.NewInt(1),
			Issuer:       issuer,
			Subject: pkix.Name{
				Organization: []string{"default"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().Add(defaultTimeGood),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth | x509.ExtKeyUsageClientAuth},
			BasicConstraintsValid: true,
		}
	}
	return &x509.Certificate{
		Signature:          csr.Signature,
		SignatureAlgorithm: csr.SignatureAlgorithm,
		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
		PublicKey:          csr.PublicKey,
		SerialNumber:       big.NewInt(1),
		Issuer:             issuer,
		Subject:            csr.Subject,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(defaultTimeGood),
		KeyUsage:           x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth | x509.ExtKeyUsageClientAuth},
	}
}
