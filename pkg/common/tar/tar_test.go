package tar

import (
	"path/filepath"
	"testing"
	"time"
)

const testWorkDir = "../../test/test-temp"

func TestTarUntar(t *testing.T) {
	dirIn := "../../test/tar_test"
	tarOut := filepath.Join(testWorkDir, "tar_test.tar")
	if err := TarDirToFile(dirIn, tarOut); err != nil {
		t.Fatal("Tar test failed", err)
	}
	time.Sleep(100)
	tarIn := filepath.Join(testWorkDir, "tar_test.tar")
	dirOut := filepath.Join(testWorkDir, "tar_test_untar")
	if err := UnTarFileToDir(tarIn, dirOut); err != nil {
		t.Fatal("Untar test failed", err)
	}
	// Add cleanup after upgrade to Go 1.14
}
