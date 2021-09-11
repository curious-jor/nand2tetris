package codewriter

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestWriteArithmetic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"add", "add", "D+M"},
		{"sub", "sub", "M-D"},
		{"neg", "neg", "-M"},
		{"eq", "eq", "D;JEQ"},
		{"gt", "gt", "D;JGT"},
		{"lt", "lt", "D;JLT"},
		{"and", "and", "D&M"},
		{"or", "or", "D|M"},
		{"not", "not", "!M"},
		{"empty", "", unsupportedCmdString},
		{"invalid nand", "nand", unsupportedCmdString},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()
			tempFile, err := os.CreateTemp(tempDir, "*.asm")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			cw := NewCodeWriter(tempFile)

			if strings.HasPrefix(test.name, "invalid") || test.name == "empty" {
				if err := cw.WriteArithmetic(test.input); err == nil {
					t.Fatalf("expected error for WriteArithmetic called with %q", test.input)
				}
			} else if err := cw.WriteArithmetic(test.input); err != nil {
				t.Fatal(err)
			}

			output, err := os.ReadFile(tempFile.Name())
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Contains(output, []byte(test.expected)) {
				t.Errorf("expected output file to contain %q with input %q", test.expected, test.input)
			}

			if err := tempFile.Close(); err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(tempFile.Name()); err != nil {
				t.Fatal(err)
			}
		})
	}
}
