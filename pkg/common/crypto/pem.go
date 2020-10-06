package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

func CertificateToPem(cert *x509.Certificate) ([]byte, error) {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}), nil
}

func RSAPublicKeyToPem(key *rsa.PublicKey) ([]byte, error) {
	pubPKCS1Bytes := x509.MarshalPKCS1PublicKey(key)
	if pubPKCS1Bytes == nil {
		return nil, errors.New("Error marshaling public key to pkcs1 block")
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubPKCS1Bytes,
	}), nil
}

func RSAKeyPairToPem(key *rsa.PrivateKey, pwd []byte) ([]byte, []byte, error) {
	pubPem, err := RSAPublicKeyToPem(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	prvPKCS1Bytes := x509.MarshalPKCS1PrivateKey(key)
	if prvPKCS1Bytes == nil {
		return nil, nil, errors.New("Error marshaling private key to pkcs1 block")
	}
	prvPemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: prvPKCS1Bytes,
	}
	if pwd != nil {
		prvPemBlock, err = x509.EncryptPEMBlock(rand.Reader, prvPemBlock.Type, prvPemBlock.Bytes, pwd, x509.PEMCipherAES256)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to encrypt private key")
		}
	}
	return pem.EncodeToMemory(prvPemBlock), pubPem, nil
}

func PemToPublicKey(pubKeyPem []byte) (*rsa.PublicKey, error) {
	keyPemBlock, _ := pem.Decode(pubKeyPem)
	if keyPemBlock == nil {
		return nil, errors.New("Failed to decode public key pem")
	}
	if keyPemBlock.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("pem is not RSA public key")
	}
	parsedKey, err := x509.ParsePKCS1PublicKey(keyPemBlock.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse public key pem")
	}
	return parsedKey, nil
}

func PemToPrivateKey(prvKeyPem []byte, pwd []byte) (*rsa.PrivateKey, error) {
	keyPemBlock, _ := pem.Decode(prvKeyPem)
	if keyPemBlock == nil {
		return nil, errors.New("Failed to decode private key pem")
	}
	// decode with passphrase
	if x509.IsEncryptedPEMBlock(keyPemBlock) {
		keyDer, err := x509.DecryptPEMBlock(keyPemBlock, pwd)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to decrypt private key")
		}
		keyPemBlock.Bytes = keyDer
		keyPemBlock.Headers = nil
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(keyPemBlock.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse private key PEM")
	}
	return parsedKey, nil
}
