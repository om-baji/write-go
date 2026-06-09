package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/om-baji/write-go/internal"
)

func EnsureFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(
		path,
		os.O_CREATE,
		0644,
	)
	if err != nil {
		return err
	}
	return f.Close()
}

func AppendFlush(segment internal.Segment, message string) internal.Segment {
	f, err := os.OpenFile(
		segment.Path,
		os.O_APPEND|os.O_WRONLY,
		os.ModeAppend,
	)

	HandlExp(err)
	defer f.Close()

	if _, err = fmt.Fprintln(f, message); err != nil {
		panic(err)
	}

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	segment.Size = int(fi.Size())

	return segment
}
