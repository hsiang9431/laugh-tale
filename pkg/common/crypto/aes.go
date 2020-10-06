package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/pkg/errors"
)

// ref: https://eli.thegreenplace.net/2019/aes-encryption-of-files-in-go/

const AesKeySize = 32

func AESEncFile(fileIn, fileOut, keyB64 string) error {
	key, err := base64.RawURLEncoding.DecodeString(keyB64)
	if err != nil {
		return errors.Wrap(err, "Failed to decode key")
	}
	fi, err := os.Stat(fileIn)
	if err != nil {
		return errors.Wrap(err, "Failed to load file for encryption")
	}
	// buffer size must be multiple of aes.BlockSize, which is 16
	bufSize := (fi.Size()>>4 + 1) << 4
	plainFileBuffer := make([]byte, bufSize)
	encFileBuffer := make([]byte, bufSize)
	// read file into buffer
	fIn, err := os.Open(fileIn)
	if _, err := fIn.Read(plainFileBuffer); err != nil {
		return errors.Wrap(err, "Failed to load file for encryption")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.Wrap(err, "Failed to create aes cipher")
	}
	// iv length = 16 byte, first 8 bytes are for original size
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)
	binary.LittleEndian.PutUint64(iv, uint64(fi.Size()))

	// encrypt the file
	enc := cipher.NewCBCEncrypter(block, iv)
	enc.CryptBlocks(encFileBuffer, plainFileBuffer)
	if err != nil {
		return errors.Wrapf(err, "Failed to create output file")
	}
	// write to output file
	f, err := os.OpenFile(fileOut, os.O_RDWR|os.O_CREATE, 0644)
	if _, err := f.Write(iv); err != nil {
		return errors.Wrapf(err, "Failed to write initial vector")
	}
	if _, err := f.Write(encFileBuffer); err != nil {
		return errors.Wrapf(err, "Failed to write encrypted file")
	}
	return nil
}

func AESDecFile(fileIn, fileOut, keyB64 string) error {
	key, err := base64.RawURLEncoding.DecodeString(keyB64)
	if err != nil {
		return errors.Wrap(err, "Failed to decode key")
	}
	fileBytes, err := ioutil.ReadFile(fileIn)
	if err != nil {
		return errors.Wrap(err, "Failed to load file for decryption")
	}
	iv := fileBytes[:aes.BlockSize]
	encFile := fileBytes[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	plainFile := make([]byte, len(encFile))
	dec := cipher.NewCBCDecrypter(block, iv)
	dec.CryptBlocks(plainFile, encFile)

	originalFileLen := binary.LittleEndian.Uint64(iv)
	if err := ioutil.WriteFile(fileOut, plainFile[:originalFileLen], 0644); err != nil {
		return errors.Wrap(err, "Failed to write decrypted file")
	}
	return nil
}
