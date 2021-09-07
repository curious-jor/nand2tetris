package codewriter

import (
	"os"
	"testing"
)

func TestWriteArithmetic(t *testing.T) {
	tempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "*.vm")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	cw := NewCodeWriter(tempFile)
	emptyCmdString := ""

	if err := cw.WriteArithmetic(emptyCmdString); err == nil {
		t.Errorf("codewriter.WriteArithmetic should return an error when called with empty string")
	}

	if err := tempFile.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(tempFile.Name()); err != nil {
		t.Fatal(err)
	}
}
