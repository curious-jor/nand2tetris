package parser

import (
	"VMtranslator/lexer"
	"os"
	"testing"
)

func (a *Command) Equals(b *Command) bool {
	return a.ct == b.ct && a.arg1 == b.arg1 && a.arg2 == b.arg2
}

func TestInitialization(t *testing.T) {
	f, err := os.Open("../StackArithmetic/SimpleAdd/SimpleAdd.vm")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	parser := NewParser(f)

	// test parser initially loaded with first lexeme in file
	expected := &lexer.Lexeme{Token: lexer.COMMAND, Value: "push"}
	if !expected.Equals(parser.lexeme) {
		t.Errorf("expected %v got %v", expected, parser.lexeme)
	}

	// test parser initially has no command
	if !(parser.cmd == nil) {
		t.Errorf("expected initialized parser to have no command but got %v instead", parser.cmd)
	}
}

func TestAdvanceSimpleAdd(t *testing.T) {
	f, err := os.Open("../StackArithmetic/SimpleAdd/SimpleAdd.vm")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	parser := NewParser(f)
	parser.Advance()

	expected := &Command{ct: C_PUSH, arg1: "constant", arg2: 7}
	if !expected.Equals(parser.cmd) {
		t.Errorf("expected %v got %v", expected, parser.cmd)
	}

	parser.Advance()
	expected = &Command{ct: C_PUSH, arg1: "constant", arg2: 8}
	if !expected.Equals(parser.cmd) {
		t.Errorf("expected %v got %v", expected, parser.cmd)
	}

	parser.Advance()
	expected = &Command{ct: C_ARITHMETIC, arg1: "add", arg2: emptyArg2}
	if !expected.Equals(parser.cmd) {
		t.Errorf("expected %v got %v", expected, parser.cmd)
	}

}

func TestAdvanceStackTest(t *testing.T) {
	f, err := os.Open("../StackArithmetic/StackTest/Stacktest.vm")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	parser := NewParser(f)
	parser.Advance()

	linesAdvanced := 1
	for parser.HasMoreCommands() {
		parser.Advance()
		linesAdvanced += 1
		t.Log(parser.cmd)
	}

	if linesAdvanced != 38 {
		t.Errorf("expected 38 lines advanced but got %d", linesAdvanced)
	}

	if err := parser.Advance(); err == nil {
		t.Errorf("expected non-nil error when advancing parser with no more commands")
	}

}

func TestCommandType(t *testing.T) {
	parser := NewParser(nil)
	tests := []struct {
		name     string
		input    *Command
		expected CommandType
	}{
		{"add", &Command{ct: C_ARITHMETIC, arg1: "add"}, C_ARITHMETIC},
		{"sub", &Command{ct: C_ARITHMETIC, arg1: "add"}, C_ARITHMETIC},
		{"neg", &Command{ct: C_ARITHMETIC, arg1: "add"}, C_ARITHMETIC},
		{"eq", &Command{ct: C_ARITHMETIC, arg1: "add"}, C_ARITHMETIC},
		{"gt", &Command{ct: C_ARITHMETIC, arg1: "add"}, C_ARITHMETIC},
		{"lt", &Command{ct: C_ARITHMETIC, arg1: "add"}, C_ARITHMETIC},
		{"and", &Command{ct: C_ARITHMETIC, arg1: "add"}, C_ARITHMETIC},
		{"or", &Command{ct: C_ARITHMETIC, arg1: "add"}, C_ARITHMETIC},
		{"not", &Command{ct: C_ARITHMETIC, arg1: "add"}, C_ARITHMETIC},
		{"push", &Command{ct: C_PUSH, arg1: "constant", arg2: 7}, C_PUSH},
		{"pop", &Command{ct: C_POP, arg1: "local", arg2: 1}, C_POP},
		{"empty", nil, emptyCommandType},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			parser.cmd = test.input
			if actual := parser.CommandType(); test.expected != actual {
				t.Errorf("expected %s got %s", test.expected.String(), parser.CommandType().String())
			}
		})
	}

}

func TestArg1(t *testing.T) {
	parser := NewParser(nil)

	passingTests := []struct {
		name     string
		input    *Command
		expected string
	}{
		{"C_ARITHMETIC add", &Command{ct: C_ARITHMETIC, arg1: "add"}, "add"},
		{"C_ARITHMETIC sub", &Command{ct: C_ARITHMETIC, arg1: "add"}, "add"},
		{"C_ARITHMETIC neg", &Command{ct: C_ARITHMETIC, arg1: "add"}, "add"},
		{"C_ARITHMETIC eq", &Command{ct: C_ARITHMETIC, arg1: "add"}, "add"},
		{"C_ARITHMETIC gt", &Command{ct: C_ARITHMETIC, arg1: "add"}, "add"},
		{"C_ARITHMETIC lt", &Command{ct: C_ARITHMETIC, arg1: "add"}, "add"},
		{"C_ARITHMETIC and", &Command{ct: C_ARITHMETIC, arg1: "add"}, "add"},
		{"C_ARITHMETIC or", &Command{ct: C_ARITHMETIC, arg1: "add"}, "add"},
		{"C_ARITHMETIC not", &Command{ct: C_ARITHMETIC, arg1: "add"}, "add"},
		{"C_PUSH push constant 7", &Command{ct: C_PUSH, arg1: "constant", arg2: 7}, "constant"},
		{"C_POP pop local 1", &Command{ct: C_POP, arg1: "local", arg2: 1}, "local"},
		{"C_RETURN arg1", &Command{ct: C_RETURN, arg1: "return"}, ""},
	}

	for _, test := range passingTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			parser.cmd = test.input
			if actual := parser.Arg1(); test.expected != actual {
				t.Errorf("expected %q got %q", test.expected, actual)
			}
		})
	}

}

func TestArg2(t *testing.T) {
	parser := NewParser(nil)

	tests := []struct {
		name     string
		input    *Command
		expected int
	}{
		{"C_ARITHMETIC add", &Command{ct: C_ARITHMETIC, arg1: "add"}, emptyArg2},
		{"C_ARITHMETIC sub", &Command{ct: C_ARITHMETIC, arg1: "add"}, emptyArg2},
		{"C_ARITHMETIC neg", &Command{ct: C_ARITHMETIC, arg1: "add"}, emptyArg2},
		{"C_ARITHMETIC eq", &Command{ct: C_ARITHMETIC, arg1: "add"}, emptyArg2},
		{"C_ARITHMETIC gt", &Command{ct: C_ARITHMETIC, arg1: "add"}, emptyArg2},
		{"C_ARITHMETIC lt", &Command{ct: C_ARITHMETIC, arg1: "add"}, emptyArg2},
		{"C_ARITHMETIC and", &Command{ct: C_ARITHMETIC, arg1: "add"}, emptyArg2},
		{"C_ARITHMETIC or", &Command{ct: C_ARITHMETIC, arg1: "add"}, emptyArg2},
		{"C_ARITHMETIC not", &Command{ct: C_ARITHMETIC, arg1: "add"}, emptyArg2},
		{"C_PUSH push constant 7", &Command{ct: C_PUSH, arg1: "constant", arg2: 7}, 7},
		{"C_PUSH pop local 1", &Command{ct: C_POP, arg1: "local", arg2: 1}, 1},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			parser.cmd = test.input
			if actual := parser.Arg2(); test.expected != actual {
				t.Errorf("expected %q got %q", test.expected, actual)
			}
		})
	}

}
