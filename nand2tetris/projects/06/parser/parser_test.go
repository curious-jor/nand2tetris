package parser

import (
	"bytes"
	"os"
	"testing"
)

func setup(t *testing.T) (*os.File, func()) {
	testFile, err := os.CreateTemp(".", "*.asm")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	teardownFn := func() {
		if err := testFile.Close(); err != nil {
			t.Fatalf("could not close temp file %v", err)
		}
		if err := os.Remove(testFile.Name()); err != nil {
			t.Fatalf("could not remove temp file %v", err)
		}
	}

	return testFile, teardownFn
}

func TestParserInitialization(t *testing.T) {
	testFile, tearDown := setup(t)
	defer tearDown()
	testFileContent := []byte("@2\nD=A\n@3\nD=D+A\n@0\nM=D\n")

	if _, err := testFile.Write(testFileContent); err != nil {
		t.Fatalf("could not write to test file %v", err)
	}

	testFile.Seek(0, 0)
	p := NewParser(testFile)
	if p.file != testFile {
		t.Errorf("parser initialized without setting file to test file.")
	}
	if p.command != nil {
		t.Errorf("parser initialized with non-nil command field value. Command was %#v", p.command)
	}

}

func TestHasMoreCommands(t *testing.T) {
	testFile, tearDown := setup(t)
	defer tearDown()
	testFileContent := []byte("@2\nD=A\n@3\nD=D+A\n@0\nM=D\n")

	if _, err := testFile.Write(testFileContent); err != nil {
		t.Fatalf("Could not write %b to test file: %s", testFileContent, testFile.Name())
	}

	// Reset file cursor to beginning before initializing parser.
	testFile.Seek(0, 0)
	p := NewParser(testFile)
	beginningBytes, err := p.lxr.Peek(2)
	if err != nil {
		t.Fatalf("Could not read bytes from beginning of file %v", err)
	}

	// Check that the function is true when called with the reader at the start of the file.
	startHasMore := p.HasMoreCommands()
	if !startHasMore {
		t.Errorf("Start of file should have more commands. Test returned %t\n", startHasMore)
	}

	// Check that HasMoreCommands didn't advance reader
	moreBeginningBytes, err := p.lxr.Peek(2)
	if err != nil {
		t.Fatalf("Could not read bytes from beginning of file.\n")
	}
	if !bytes.Equal(moreBeginningBytes, beginningBytes) {
		t.Errorf("Reader was advanced. Bytes from beginning were %b. But peeked %b after call to HasMoreCommands\n", beginningBytes, moreBeginningBytes)
	}

	testFileContentLength := len(bytes.Split(testFileContent, []byte{'\n'}))
	t.Logf("testFileContent has %d lines\n", testFileContentLength)

	// Check that HasMoreCommands is true when called from second to last line
	// Subtract 2 from the length of the test file to read the right amount of lines.
	// A newline at the end of a file counts as a line, so we need to subtract two lines
	// instead of one.
	for i := 0; i < testFileContentLength-2; i++ {
		t.Log(p.lxr.ReadBytes('\n'))
	}

	secondToLastHasMore := p.HasMoreCommands()
	if !secondToLastHasMore {
		t.Errorf("Second to last line should have more commands. Test returned %t\n", secondToLastHasMore)
	}

	// Should return false if called from last line
	lastCommand, err := p.lxr.ReadBytes('\n')
	if err != nil {
		t.Fatalf("Could not read last line of file. Got %b\n", lastCommand)
	}
	t.Logf("Last line: %b\n", lastCommand)
	lastHasNone := !p.HasMoreCommands()
	if !lastHasNone {
		t.Errorf("Last line should have no more commands. Test returned %t\n", lastHasNone)
	}
}

func TestAdvance(t *testing.T) {
	testFile, tearDown := setup(t)
	defer tearDown()
	testFileContentWithComments :=
		`// This file is part of www.nand2tetris.org
	// and the book "The Elements of Computing Systems"
	// by Nisan and Schocken, MIT Press.
	// File name: projects/06/add/Add.asm
	
	// Computes R0 = 2 + 3  (R0 refers to RAM[0])
	
	@2
	D=A
	@3
	D=D+A
	@0
	M=D

	`

	if _, err := testFile.Write([]byte(testFileContentWithComments)); err != nil {
		t.Fatalf("could not write to test file %v", err)
	}

	testFile.Seek(0, 0)
	p := NewParser(testFile)

	// Advance should maintain the invariant that a parser is initialized with no current command.
	if p.command != nil {
		t.Errorf("parser has a non-nil initial value for command")
	}

	p.Advance()

	// First command after call to Advance() should not be nil
	if p.command == nil {
		t.Errorf("parser command was nil after call to Advance()")
	}

	// Advance to the end of the file.
	for p.HasMoreCommands() {
		t.Log(p.lexeme)
		p.Advance()
	}

	// Advance should return error after commands are exhausted.
	if err := p.Advance(); err == nil {
		t.Errorf("parser did not return error after call to Advance with no more commands.")
	}
}

func TestAdvanceParseA_COMMAND(t *testing.T) {
	tests := map[string]Command{
		"@0\n":             A_COMMAND{symbol: "0"},
		"@R0\n":            A_COMMAND{symbol: "R0"},
		"@_INNER_LOOP\n":   A_COMMAND{"_INNER_LOOP"},
		"@INFINITE_LOOP\n": A_COMMAND{"INFINITE_LOOP"},
		"@counter\n":       A_COMMAND{symbol: "counter"},
		"@address\n":       A_COMMAND{symbol: "address"},
	}

	for line, expected := range tests {
		testFile, tearDown := setup(t)
		defer tearDown()

		n, err := testFile.WriteString(line)
		if n < len(line) {
			t.Fatalf("wrote %d bytes to test file. should have written %d", n, len(line))
		}
		if err != nil {
			t.Fatalf("could not write to test file %v", err)
		}

		testFile.Seek(0, 0)
		p := NewParser(testFile)
		if err = p.Advance(); err != nil {
			t.Fatalf("could not advance parser for line: %s %v", line, err)
		}

		if p.command != expected {
			t.Errorf("expected command: %v but got %v", expected, p.command)
		}
	}
}

// Check that Advance() parses C Commands correctly
func TestAdvanceParseC_COMMAND(t *testing.T) {
	// Map c instructions to their expected command parses
	testCommands := map[string]Command{
		"0;JMP\n": C_COMMAND{dest: "null", comp: "0", jump: "JMP"},
		"D=A\n":   C_COMMAND{dest: "D", comp: "A", jump: "null"},
	}

	// Create a file for each test instruction. Then create a parser and attempt to parse the instruction
	// into a C_COMMAND
	for line, expectedCmd := range testCommands {
		testFile, tearDown := setup(t)
		defer tearDown()

		n, err := testFile.WriteString(line)
		if n < len(line) {
			t.Fatalf("wrote %d bytes to test file. should have written %d", n, len(line))
		}
		if err != nil {
			t.Fatalf("could not write to test file %v", err)
		}

		testFile.Seek(0, 0)
		p := NewParser(testFile)
		if err = p.Advance(); err != nil {
			t.Fatalf("could not advance parser for line: %s %v", line, err)
		}

		if p.command != expectedCmd {
			t.Errorf("expected command: %v but got %v", expectedCmd, p.command)
		}

	}
}

func TestAdvanceParseL_COMMAND(t *testing.T) {
	testCommands := map[string]Command{
		"(LOOP)\n":        L_COMMAND{symbol: "LOOP"},
		"(INFINITE_LOOP)": L_COMMAND{symbol: "INFINITE_LOOP"},
	}

	for line, expected := range testCommands {
		testFile, tearDown := setup(t)
		defer tearDown()

		n, err := testFile.WriteString(line)
		if n < len(line) {
			t.Fatalf("wrote %d bytes to test file. should have written %d", n, len(line))
		}
		if err != nil {
			t.Fatalf("could not write to test file %v", err)
		}

		testFile.Seek(0, 0)
		p := NewParser(testFile)

		if err = p.Advance(); err != nil {
			t.Fatalf("could not advance parser for line: %s %v", line, err)
		}

		if p.command != expected {
			t.Errorf("expected command: %v but got %v", expected, p.command)
		}
	}
}

func TestCommandType(t *testing.T) {
	testFile, tearDown := setup(t)
	defer tearDown()

	p := NewParser(testFile)

	// parser command should be empty initially
	if p.CommandType() != nil {
		t.Errorf("parser should have no command initially. Command was %#v", p.command)
	}

	p.command = A_COMMAND{}
	if p.CommandType() != (A_COMMAND{}) {
		t.Errorf("parser's command type should be A_COMMAND. Command type was %#v", p.CommandType())
	}

	p.command = A_COMMAND{symbol: "2"}
	if p.CommandType() != (A_COMMAND{}) {
		t.Errorf("parser's command type should be A_COMMAND. Command type was %#v", p.CommandType())
	}

	p.command = L_COMMAND{}
	if p.CommandType() != (L_COMMAND{}) {
		t.Errorf("parser's command type should be L_COMMAND. Command type was %#v", p.CommandType())
	}
	p.command = L_COMMAND{symbol: "2"}
	if p.CommandType() != (L_COMMAND{}) {
		t.Errorf("parser's command type should be L_COMMAND. Command type was %#v", p.CommandType())
	}

	p.command = C_COMMAND{}
	if p.CommandType() != (C_COMMAND{}) {
		t.Errorf("parser's command type should be C_COMMAND. Command type was %#v", p.CommandType())
	}

	p.command = C_COMMAND{dest: "D", comp: "D+A", jump: "JMP"}
	if p.CommandType() != (C_COMMAND{}) {
		t.Errorf("parser's command type should be C_COMMAND. Command type was %#v", p.CommandType())
	}
}

func TestSymbol(t *testing.T) {
	testFile, tearDown := setup(t)
	defer tearDown()
	p := NewParser(testFile)

	if ps, err := p.Symbol(); ps != "" || err == nil {
		t.Error("symbol returned non-zero values when called on newly initalized parser")
	}

	p.command = A_COMMAND{symbol: "2"}
	if ps, err := p.Symbol(); ps != "2" || err != nil {
		t.Errorf("incorrect symbol returned for parser with A_COMMAND. Expected 2; Got %s %v", ps, err)
	}

	p.command = L_COMMAND{symbol: "222"}
	if ps, err := p.Symbol(); ps != "222" || err != nil {
		t.Errorf("incorrect symbol returned for parser with L_COMMAND. Expected 222; Got %s %v", ps, err)
	}

	p.command = C_COMMAND{dest: "D", comp: "D+A", jump: "JMP"}
	if _, err := p.Symbol(); err == nil {
		t.Errorf("symbol called on C_COMMAND should return an error.")
	}
}

func TestDest(t *testing.T) {
	testFile, tearDown := setup(t)
	defer tearDown()
	p := NewParser(testFile)

	// Dest called on parser with empty command should return err
	if _, err := p.Dest(); err == nil {
		t.Errorf("dest called on parser with empty command should return non-nil error")
	}

	// dest should return error if called on parser with non-C command
	p.command = A_COMMAND{}
	if _, err := p.Dest(); err == nil {
		t.Errorf("dest called on parser with A command should return non-nil error")
	}

	p.command = L_COMMAND{}
	if _, err := p.Dest(); err == nil {
		t.Errorf("dest called on parser with L command should return non-nil error")
	}

	// C commands have 8 possibilities for the dest mnemonic. Test for the "null", single, and double mnemonics.
	p.command = C_COMMAND{dest: "D"}
	if dest, err := p.Dest(); dest != "D" || err != nil {
		t.Errorf("Expected D. Got %s %v", dest, err)
	}

	p.command = C_COMMAND{dest: "MD"}
	if dest, err := p.Dest(); dest != "MD" || err != nil {
		t.Errorf("Expected MD. Got %s %v", dest, err)
	}

	p.command = C_COMMAND{dest: "null"}
	if dest, err := p.Dest(); dest != "null" || err != nil {
		t.Errorf("Expected null. Got %s %v", dest, err)
	}
}

func TestComp(t *testing.T) {
	testFile, tearDown := setup(t)
	defer tearDown()
	p := NewParser(testFile)

	// Comp called on parser with empty command should return err
	if _, err := p.Comp(); err == nil {
		t.Errorf("comp called on parser with empty command should return non-nil error")
	}

	// comp should return error if called on parser with non-C command
	p.command = A_COMMAND{}
	if _, err := p.Comp(); err == nil {
		t.Errorf("comp called on parser with A command should return non-nil error")
	}

	p.command = L_COMMAND{}
	if _, err := p.Comp(); err == nil {
		t.Errorf("comp called on parser with L command should return non-nil error")
	}

	// Test a memory comp field
	p.command = C_COMMAND{comp: "A"}
	if comp, err := p.Comp(); comp != "A" || err != nil {
		t.Errorf("expected A got %s %v", comp, err)
	}

	// Test addition
	p.command = C_COMMAND{comp: "A+1"}
	if comp, err := p.Comp(); comp != "A+1" || err != nil {
		t.Errorf("expected A+1 got %s %v", comp, err)
	}

	// Test subtraction
	p.command = C_COMMAND{comp: "M-1"}
	if comp, err := p.Comp(); comp != "M-1" || err != nil {
		t.Errorf("expected M-1 got %s %v", comp, err)
	}

	// Test bitwise AND
	p.command = C_COMMAND{comp: "D&M"}
	if comp, err := p.Comp(); comp != "D&M" || err != nil {
		t.Errorf("expected D&M got %s %v", comp, err)
	}
	// Test bitwise OR
	p.command = C_COMMAND{comp: "D|A"}
	if comp, err := p.Comp(); comp != "D|A" || err != nil {
		t.Errorf("expected M-1 got %s %v", comp, err)
	}

	// Test blank comp field
	p.command = C_COMMAND{comp: ""}
	if comp, err := p.Comp(); comp != "" || err != nil {
		t.Errorf("expected \"\" got %s %v", comp, err)
	}
}

func TestJump(t *testing.T) {
	testFile, tearDown := setup(t)
	defer tearDown()
	p := NewParser(testFile)

	// jump should return error if called on parser with no command
	if _, err := p.Jump(); err == nil {
		t.Errorf("jump called on parser with empty command should return non-nil error")
	}

	// jump should return error if called on parser with non-C command
	p.command = A_COMMAND{}
	if _, err := p.Jump(); err == nil {
		t.Errorf("jump called on parser with A command should return non-nil error")
	}

	p.command = L_COMMAND{}
	if _, err := p.Jump(); err == nil {
		t.Errorf("jump called on parser with L command should return non-nil error")
	}

	// test jump greater than
	p.command = C_COMMAND{jump: "JGT"}
	if jump, err := p.Jump(); jump != "JGT" || err != nil {
		t.Errorf("Expected JGT got %s %v", jump, err)
	}

	// test unconditional jump
	p.command = C_COMMAND{jump: "JMP"}
	if jump, err := p.Jump(); jump != "JMP" || err != nil {
		t.Errorf("Expected JMP got %s %v", jump, err)
	}

	// test null jump
	p.command = C_COMMAND{jump: "null"}
	if jump, err := p.Jump(); jump != "null" || err != nil {
		t.Errorf("Expected null got %s %v", jump, err)
	}
}
