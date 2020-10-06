package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io/ioutil"
	"laugh-tale/pkg/common/tar"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

const testFileDir = "../../test"
const testWorkDir = "../../test/test-temp"

func TestRSAAndPem(t *testing.T) {
	test_string := RandB64String(500)
	prvKeyPem, pubKeyPem, pwd, err := GenRSAKeyPairRandPem()
	if err != nil {
		t.Fatal("GenRSAKeyPairRandPem failed", err)
	}
	pwdByte, err := base64.RawURLEncoding.DecodeString(pwd)
	if err != nil {
		t.Fatal("base64.RawStdEncoding.Decode failed", err)
	}
	prvKey, err := PemToPrivateKey(prvKeyPem, pwdByte)
	if err != nil {
		t.Fatal("PemToPrivateKey failed", err)
	}
	pubKey, err := PemToPublicKey(pubKeyPem)
	if err != nil {
		t.Fatal("PemToPublicKey failed", err)
	}
	cipher, err := RSAEncrypt(test_string, pubKey)
	if err != nil {
		t.Fatal("RSAEncrypt failed", err)
	}
	text, err := RSADecrypt(cipher, prvKey)
	if err != nil {
		t.Fatal("RSADecrypt failed", err)
	}
	if text != test_string {
		t.Fatal("decrypted message does not match original")
	}
	t.Log(string(pubKeyPem))
	t.Log(string(prvKeyPem))
	t.Log(string(pwd))
}

func compareFile(f1, f2 string) error {
	text, err := ioutil.ReadFile(f1)
	if err != nil {
		return err
	}
	decText, err := ioutil.ReadFile(f2)
	if err != nil {
		return err
	}
	if res := bytes.Compare(text, decText); res != 0 {
		return errors.New("decrypted text does not match original file")
	}
	return nil
}
func TestAES(t *testing.T) {
	encIn := filepath.Join(testFileDir, "Island's Sunrise.txt")
	encOut := filepath.Join(testWorkDir, "winnie-cant-see.bin")
	pwd, err := AESEncFileRand(encIn, encOut)
	if err != nil {
		t.Fatal("EncFileRand failed", err)
	}
	t.Log(pwd)
	decIn := encOut
	decOut := filepath.Join(testWorkDir, "winnie-QQ.txt")
	if err := AESDecFile(decIn, decOut, pwd); err != nil {
		t.Fatal("DecFile failed", err)
	}
	if err := compareFile(encIn, decOut); err != nil {
		t.Fatal("compare failed", err)
	}
}

func TestEncTar(t *testing.T) {
	dirIn := filepath.Join(testFileDir, "tar_test")
	tarOut := filepath.Join(testWorkDir, "tar_test.tar")
	if err := tar.TarDirToFile(dirIn, tarOut); err != nil {
		t.Fatal("Tar test failed", err)
	}
	encIn := tarOut
	encOut := filepath.Join(testWorkDir, "enc_tar_test.bin")
	pwd, err := AESEncFileRand(encIn, encOut)
	if err != nil {
		t.Fatal("EncFileRand failed", err)
	}
	t.Log(pwd)
	decIn := encOut
	decOut := filepath.Join(testWorkDir, "dec_tar_test.tar")
	if err := AESDecFile(decIn, decOut, pwd); err != nil {
		t.Fatal("DecFile failed", err)
	}
	if err := compareFile(encIn, decOut); err != nil {
		t.Fatal("compare failed", err)
	}
}

func AESEncFileRand(fileIn, fileOut string) (string, error) {
	_, pwd := RandBytesAndB64(AesKeySize)
	return pwd, AESEncFile(fileIn, fileOut, pwd)
}

func GenRSAKeyPairRandPem() ([]byte, []byte, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, RSABitSize)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "Error generating rsa key pair")
	}
	pwdB, pwdStr := RandBytesAndB64(AesKeySize)
	prvPem, pubPem, err := RSAKeyPairToPem(key, pwdB)
	if err != nil {
		return nil, nil, "", err
	}
	return prvPem, pubPem, pwdStr, nil
}
