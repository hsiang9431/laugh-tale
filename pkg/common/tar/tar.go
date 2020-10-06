package tar

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// ref: https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

func TarDirToFile(inDir, outFile string) error {
	if _, err := os.Stat(inDir); err != nil {
		return errors.Wrap(err, "Unable to tar directory")
	}
	inDir, err := filepath.Abs(inDir)
	if err != nil {
		return errors.Wrap(err, "Unable to tar directory")
	}
	outFileW, err := os.Create(outFile)
	if err != nil {
		return errors.Wrapf(err, "Failed to create outpul file %s", outFile)
	}
	tw := tar.NewWriter(outFileW)
	defer tw.Close()
	return filepath.Walk(inDir, func(file string, fi os.FileInfo, err error) error {
		// return on any error
		if err != nil {
			return errors.Wrap(err, "Error with filepath.Walk")
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.Wrapf(err, "file info header error on file %s", file)
		}
		// change to relative path
		header.Name = strings.TrimPrefix(
			strings.Replace(file, inDir, "", -1),
			string(filepath.Separator))
		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return errors.Wrapf(err, "tar writer error on file %s", file)
		}
		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return errors.Wrapf(err, "cannot open file %s", file)
		}
		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return errors.Wrapf(err, "cannot copy file %s", file)
		}
		return f.Close()
	})
}

func UnTarFileToDir(filename, outDir string) error {
	inFile, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to open file %s", filename)
	}
	fi, err := os.Stat(outDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(outDir, 0755)
		}
	} else if fi.Mode().IsRegular() {
		return errors.New(outDir + " is not a directory")
	}
	tr := tar.NewReader(inFile)
	for {
		header, err := tr.Next()
		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil
		// return any other error
		case err != nil:
			return errors.Wrap(err, "Error reading tar file")
		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}
		// the target location where the dir/file should be created
		target := filepath.Join(outDir, header.Name)

		fileDir := filepath.Dir(target)
		if _, err := os.Stat(fileDir); err != nil {
			if err := os.MkdirAll(fileDir, 0755); err != nil {
				return errors.Wrapf(err, "Failed to create directory %s", target)
			}
		}
		f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
		if err != nil {
			return errors.Wrapf(err, "Failed to create file %s", target)
		}
		// copy over contents
		if _, err := io.Copy(f, tr); err != nil {
			return errors.Wrapf(err, "Failed to copy file %s", target)
		}
		f.Close()
	}
}
