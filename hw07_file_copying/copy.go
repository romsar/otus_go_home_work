package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const (
	chunkSize int64 = 1
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	in, err := os.Open(fromPath)
	if err != nil {
		return wrapErr(err)
	}
	defer doClose(in)

	f, err := in.Stat()
	if err != nil {
		return wrapErr(err)
	}

	fromSize := f.Size()

	if fromSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > fromSize {
		return ErrOffsetExceedsFileSize
	}

	out, err := os.Create(toPath)
	if err != nil {
		return wrapErr(err)
	}
	defer doClose(out)

	if _, err = in.Seek(offset, io.SeekStart); err != nil {
		return wrapErr(err)
	}

	total := calculateTotal(fromSize, offset)

	if limit == 0 {
		limit = total
	}

	bar := pb.StartNew(int(limit))

	for {
		size := chunkSize
		if size > limit {
			size = limit
		}

		n, err := io.CopyN(out, in, size)

		eof := errors.Is(err, io.EOF)

		if err != nil && !eof {
			return wrapErr(err)
		}

		bar.Add(int(n))
		limit -= n

		if eof || limit <= 0 {
			break
		}
	}

	bar.Finish()

	return nil
}

func calculateTotal(size int64, offset int64) int64 {
	return size - offset
}

func wrapErr(err error) error {
	return fmt.Errorf("error while copy file: %w", err)
}

func doClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println("error while close: ", err)
	}
}
