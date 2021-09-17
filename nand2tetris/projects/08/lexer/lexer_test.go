package lexer

import (
	"os"
	"testing"
)

func TestSimpleAdd(t *testing.T) {
	f, err := os.Open("../StackArithmetic/SimpleAdd/SimpleAdd.vm")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lxr := NewLexer(f)
	expected := []struct {
		token Token
		value string
	}{
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "7"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "8"},
		{COMMAND, "add"},
	}

	i := 0
	for _, lexeme := lxr.NextToken(); lexeme.Token != EOF; {
		if expected[i].token != lexeme.Token || expected[i].value != lexeme.Value {
			t.Errorf("comparison failed at token %d expected (%s, %q) got (%s, %q)", i, expected[i].token.String(), expected[i].value, lexeme.Token.String(), lexeme.Value)
		}
		_, lexeme = lxr.NextToken()
		i += 1
	}
}

func TestStackTest(t *testing.T) {
	f, err := os.Open("../StackArithmetic/StackTest/StackTest.vm")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lxr := NewLexer(f)
	expected := []Lexeme{
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "17"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "17"},
		{COMMAND, "eq"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "17"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "16"},
		{COMMAND, "eq"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "16"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "17"},
		{COMMAND, "eq"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "892"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "891"},
		{COMMAND, "lt"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "891"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "892"},
		{COMMAND, "lt"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "891"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "891"},
		{COMMAND, "lt"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "32767"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "32766"},
		{COMMAND, "gt"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "32766"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "32767"},
		{COMMAND, "gt"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "32766"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "32766"},
		{COMMAND, "gt"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "57"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "31"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "53"},
		{COMMAND, "add"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "112"},
		//
		{COMMAND, "sub"},
		{COMMAND, "neg"},
		{COMMAND, "and"},
		//
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "82"},
		//
		{COMMAND, "or"},
		{COMMAND, "not"},
	}

	for i, expctd := range expected {
		_, actual := lxr.NextToken()

		if !expctd.Equals(actual) {
			t.Errorf("comparison failed at token %d. expected %#v got %#v", i, expctd, actual)
		}
	}
}

func TestSimpleAddFilePosition(t *testing.T) {
	f, err := os.Open("../StackArithmetic/SimpleAdd/SimpleAdd.vm")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lxr := NewLexer(f)

	expected := []struct {
		lexeme *Lexeme
		line   int
		col    int
	}{
		{&Lexeme{Token: COMMAND, Value: "push"}, 7, 1},
		{&Lexeme{Token: ARG, Value: "constant"}, 7, 6},
		{&Lexeme{Token: ARG, Value: "7"}, 7, 15},
		{&Lexeme{Token: COMMAND, Value: "push"}, 8, 1},
		{&Lexeme{Token: ARG, Value: "constant"}, 8, 6},
		{&Lexeme{Token: ARG, Value: "8"}, 8, 15},
		{&Lexeme{Token: COMMAND, Value: "add"}, 9, 1},
	}

	for _, exp := range expected {
		actualPos, actualLex := lxr.NextToken()
		if !(exp.lexeme.Equals(actualLex)) || exp.line != actualPos.Line || exp.col != actualPos.Col {
			t.Errorf("expected %v, Line: %d, Col: %d but got %v, Line: %d, Col: %d)",
				exp.lexeme, exp.line, exp.col, actualLex, actualPos.Line, actualPos.Col,
			)
		}
	}
}
