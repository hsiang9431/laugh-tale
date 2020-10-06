package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

func Compress(src, dst string) error {
	srcB, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.Wrap(err, "failed to read input file")
	}
	d, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to open file at destination")
	}
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	_, err = gw.Write(srcB)
	if err != nil {
		return errors.Wrap(err, "failed to write zip file")
	}
	return nil
}

func Decompress(src, dst string) error {
	srcB, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.Wrap(err, "failed to read input file")
	}
	d, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to open file at destination")
	}
	defer d.Close()
	gr, err := gzip.NewReader(bytes.NewBuffer(srcB))
	if err != nil {
		return errors.Wrap(err, "failed to unzip input file")
	}
	defer gr.Close()
	_, err = io.Copy(d, gr)
	if err != nil {
		return errors.Wrap(err, "failed to write output file")
	}
	return nil
}
