package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/om-baji/write-go/internal"
)

func TestEnsureFileCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "newfile.log")

	err := EnsureFile(path)
	if err != nil {
		t.Fatalf("EnsureFile failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("expected empty file, got size %d", info.Size())
	}
}

func TestEnsureFileCreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "sub", "file.log")

	err := EnsureFile(path)
	if err != nil {
		t.Fatalf("EnsureFile failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file should exist: %v", err)
	}
}

func TestEnsureFileIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "idem.log")

	if err := EnsureFile(path); err != nil {
		t.Fatalf("first EnsureFile failed: %v", err)
	}

	if err := EnsureFile(path); err != nil {
		t.Fatalf("second EnsureFile failed: %v", err)
	}
}

func TestAppendFlush(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "segment.log")

	if err := EnsureFile(path); err != nil {
		t.Fatalf("EnsureFile failed: %v", err)
	}

	seg := internal.Segment{
		Id:   7,
		Size: 0,
		Path: path,
	}

	var err error
	seg, err = AppendFlush(seg, "first log line")
	if err != nil {
		t.Fatalf("AppendFlush failed: %v", err)
	}

	if seg.Size == 0 {
		t.Error("expected size > 0 after flush")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	if len(content) == 0 {
		t.Error("expected content in file")
	}

	seg, err = AppendFlush(seg, "second log line")
	if err != nil {
		t.Fatalf("second AppendFlush failed: %v", err)
	}
	if seg.Size == 0 {
		t.Error("expected size > 0 after second flush")
	}
}
