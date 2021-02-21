package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
)

func extract(source io.ReaderAt, length int64, target string) error {
	reader, err := zip.NewReader(source, length)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		if err := extractFile(file, target); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(file *zip.File, target string) (err error) {
	path, err := filepath.Abs(filepath.Join(target, file.Name))
	if err != nil {
		return err
	}

	if !strings.HasPrefix(path, target) {
		return errors.New("dangerous filename detected: " + path)
	}

	fileReader, err := file.Open()
	if err != nil {
		return err
	}

	defer func() {
		if cerr := fileReader.Close(); cerr != nil {
			err = cerr
		}
	}()

	err = os.MkdirAll(filepath.Dir(path), 0750)
	if err != nil {
		return err
	}

	targetFile, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := targetFile.Close(); cerr != nil {
			err = cerr
		}
	}()

	_, err = io.Copy(targetFile, fileReader)
	if err != nil {
		return err
	}

	return
}
