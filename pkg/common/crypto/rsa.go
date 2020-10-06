package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"

	"github.com/pkg/errors"
)

// ref: https://gist.github.com/sdorra/1c95de8cb80da31610d2ad767cd6f251
// ref: https://stackoverflow.com/questions/37316370/how-create-rsa-private-key-with-passphrase-in-golang
// ref: https://gist.github.com/jshap70/259a87a7146393aab5819873a193b88c

const RSABitSize = 4096

func RSAEncrypt(text string, rsaPub *rsa.PublicKey) (string, error) {
	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(text))
	if err != nil {
		return "", errors.Wrap(err, "Failed to encrypt message")
	}
	return base64.RawURLEncoding.EncodeToString(cipher), nil
}

func RSADecrypt(cipher string, rsaPrv *rsa.PrivateKey) (string, error) {
	cipherBytes, err := base64.RawURLEncoding.DecodeString(cipher)
	if err != nil {
		return "", errors.Wrap(err, "Base64 decode on cipher failed")
	}
	text, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrv, cipherBytes)
	if err != nil {
		return "", errors.Wrap(err, "Failed to decrypt cipher")
	}
	return string(text), nil
}
