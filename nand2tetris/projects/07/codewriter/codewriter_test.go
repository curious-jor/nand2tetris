package codewriter

import (
	"VMtranslator/parser"
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
		{"error empty", "", unsupportedCmdString},
		{"error nand", "nand", unsupportedCmdString},
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

			if !strings.HasPrefix(test.name, "error") {
				if err := cw.WriteArithmetic(test.input); err != nil {
					t.Fatal(err)
				}
			} else if err := cw.WriteArithmetic(test.input); err == nil {
				t.Fatalf("expected error from call to WriteArithmetic with input %s", test.input)
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

type pushPopInput struct {
	command parser.CommandType
	segment string
	index   int
}

func TestWritePushPop(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    pushPopInput
		expected []string
	}{
		{"push contant 1", pushPopInput{parser.C_PUSH, "constant", 1}, []string{"@1"}},
		{"push local 0", pushPopInput{parser.C_PUSH, "local", 0}, []string{"@LCL", "@0"}},
		{"push argument 2", pushPopInput{parser.C_PUSH, "argument", 2}, []string{"@ARG", "@2"}},
		{"push this 2", pushPopInput{parser.C_PUSH, "this", 2}, []string{"@THIS", "@2"}},
		{"push that 5", pushPopInput{parser.C_PUSH, "that", 5}, []string{"@THAT", "@5"}},
		{"push temp 6", pushPopInput{parser.C_PUSH, "temp", 6}, []string{"@R5", "@6"}},
		{"pop local 0", pushPopInput{parser.C_POP, "local", 0}, []string{"@LCL", "@0"}},
		{"pop argument 1", pushPopInput{parser.C_POP, "argument", 1}, []string{"@ARG", "@1"}},
		{"pop this 6", pushPopInput{parser.C_POP, "this", 6}, []string{"@THIS", "@6"}},
		{"pop that 5", pushPopInput{parser.C_POP, "that", 5}, []string{"@THAT", "@5"}},
		{"pop temp 6", pushPopInput{parser.C_POP, "temp", 6}, []string{"@R5", "@6"}},
		{"push pointer 0", pushPopInput{parser.C_PUSH, "pointer", 0}, []string{"@THIS"}},
		{"push pointer 1", pushPopInput{parser.C_PUSH, "pointer", 1}, []string{"@THAT"}},
		{"pop pointer 0", pushPopInput{parser.C_POP, "pointer", 0}, []string{"@THIS"}},
		{"pop pointer 1", pushPopInput{parser.C_POP, "pointer", 1}, []string{"@THAT"}},
		{"error pop constant 6", pushPopInput{parser.C_POP, "constant", 6}, []string{}},
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

			if !strings.HasPrefix(test.name, "error") {
				if err := cw.WritePushPop(test.input.command, test.input.segment, test.input.index); err != nil {
					t.Fatal(err)
				}
			} else if err := cw.WritePushPop(test.input.command, test.input.segment, test.input.index); err == nil {
				t.Fatalf("expected error from call to WritePushPop with input %v", test.input)
			}

			output, err := os.ReadFile(tempFile.Name())
			if err != nil {
				t.Fatal(err)
			}

			for _, expectedStr := range test.expected {
				if !bytes.Contains(output, []byte(expectedStr)) {
					t.Errorf("expected output file to contain %q with input %v", expectedStr, test.input)
				}
			}

			if err := cw.Close(); err != nil {
				t.Fatal(err)
			}

			if err := os.Remove(tempFile.Name()); err != nil {
				t.Fatal(err)
			}
		})
	}
}
