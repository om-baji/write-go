package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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
	err := EnsureFile(segment.Path)

	if err != nil {
		panic(err)
	}

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

	lm := os.Getenv("SEGMENT_LIMIT")

	println("Limit ", lm)

	if lm == "" {
		lm = "64"
	}

	limit, err := strconv.Atoi(lm)

	if err != nil {
		panic(err)
	}

	if limit*1024 <= int(fi.Size()) {
		segment.Id++
		segment.Path = "./data/wal_segment" + strconv.Itoa(segment.Id) + ".log"
		segment.Size = 0
	} else {
		segment.Size = int(fi.Size())
	}

	return segment
}
