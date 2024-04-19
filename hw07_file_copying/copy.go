package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"math"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	fileSize := fileInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}
	if limit > fileSize || limit == 0 {
		limit = fileSize
	}
	chunk := int64(128)
	if limit < chunk {
		chunk = limit
	}
	tail := fileSize - offset
	if limit < tail {
		tail = limit
	}
	bar := pb.StartNew(int(math.Ceil(float64(tail) / float64(chunk))))
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
	_, err = r.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	seek := offset
	for seek < limit+offset {
		seek += chunk
		if seek >= limit+offset {
			chunk = limit + offset + chunk - seek
		}
		_, err = io.CopyN(w, r, chunk)
		bar.Increment()
		if err != nil {
			if errors.Is(err, io.EOF) {
				bar.Finish()
				return nil
			}
			return err
		}
		_, err = r.Seek(seek, io.SeekStart)
		if err != nil {
			return err
		}
	}
	bar.Finish()
	return nil
}
