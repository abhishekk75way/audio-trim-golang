package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func CreateZip(zipPath string, files []string) error {
	z, _ := os.Create(zipPath)
	defer z.Close()

	w := zip.NewWriter(z)
	defer w.Close()

	for _, f := range files {
		file, _ := os.Open(f)
		defer file.Close()

		wr, _ := w.Create(filepath.Base(f))
		io.Copy(wr, file)
	}
	return nil
}
