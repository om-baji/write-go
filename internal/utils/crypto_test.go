package utils

import (
	"testing"
)

func TestGenerateCRC(t *testing.T) {
	data := []byte("hello world")
	result := GenerateCRC(data)

	if result == 0 {
		t.Error("expected non-zero CRC for non-empty input")
	}
}

func TestGenerateCRCConsistent(t *testing.T) {
	data := []byte("consistent-test")

	a := GenerateCRC(data)
	b := GenerateCRC(data)

	if a != b {
		t.Errorf("expected consistent CRC, got %d and %d", a, b)
	}
}

func TestGenerateCRCDifferentInputs(t *testing.T) {
	a := GenerateCRC([]byte("alpha"))
	b := GenerateCRC([]byte("beta"))

	if a == b {
		t.Error("expected different CRCs for different inputs")
	}
}

func TestGenerateCRCEmpty(t *testing.T) {
	result := GenerateCRC([]byte{})

	_ = result
}
