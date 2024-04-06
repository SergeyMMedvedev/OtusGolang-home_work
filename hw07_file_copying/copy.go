package main

import (
	"errors"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}
	if limit > fileSize || limit == 0 {
		limit = fileSize
	}
	w, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer w.Close()
	r, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer r.Close()

	return nil
}
