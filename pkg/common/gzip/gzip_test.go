package gzip

import (
	"os"
	"path/filepath"
	"testing"

	"time"
)

const testFile = "../../../test/Island's Sunrise.txt"
const testWorkDir = "../../../test/test-temp"

func TestGzip(t *testing.T) {
	os.MkdirAll(testWorkDir, 0777)
	zipOut := filepath.Join(testWorkDir, "Island's Sunrise.zip")
	if err := Compress(testFile, zipOut); err != nil {
		t.Fatal("Zip test failed", err)
	}
	time.Sleep(100)
	zipIn := filepath.Join(testWorkDir, "Island's Sunrise.zip")
	dirOut := filepath.Join(testWorkDir, "Island's Sunrise-unzipped.txt")
	if err := Decompress(zipIn, dirOut); err != nil {
		t.Fatal("Zip test failed", err)
	}
}
