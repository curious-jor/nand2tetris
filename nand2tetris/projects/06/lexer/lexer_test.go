package lexer

import (
	"os"
	"path/filepath"
	"testing"
)

type Lexeme struct {
	token Token
	value string
}

func TestNoSymbols(t *testing.T) {
	tests := []struct {
		name     string
		fpath    string
		expected []Lexeme
	}{
		{"Add", "../add/Add.asm", []Lexeme{
			{AT, "@"},
			{CONSTANT, "2"},
			{DEST, "D"},
			{EQUALS, "="},
			{COMP, "A"},
			{AT, "@"},
			{CONSTANT, "3"},
			{DEST, "D"},
			{EQUALS, "="},
			{COMP, "D+A"},
			{AT, "@"},
			{CONSTANT, "0"},
			{DEST, "M"},
			{EQUALS, "="},
			{COMP, "D"},
		}},

		{"MaxL", "../max/MaxL.asm", []Lexeme{
			{AT, "@"},
			{CONSTANT, "0"},
			{DEST, "D"},
			{EQUALS, "="},
			{COMP, "M"},
			{AT, "@"},
			{CONSTANT, "1"},
			{DEST, "D"},
			{EQUALS, "="},
			{COMP, "D-M"},
			{AT, "@"},
			{CONSTANT, "10"},
			{COMP, "D"},
			{SEMICOLON, ";"},
			{JUMP, "JGT"},
			{AT, "@"},
			{CONSTANT, "1"},
			{DEST, "D"},
			{EQUALS, "="},
			{COMP, "M"},
			{AT, "@"},
			{CONSTANT, "12"},
			{COMP, "0"},
			{SEMICOLON, ";"},
			{JUMP, "JMP"},
			{AT, "@"},
			{CONSTANT, "0"},
			{DEST, "D"},
			{EQUALS, "="},
			{COMP, "M"},
			{AT, "@"},
			{CONSTANT, "2"},
			{DEST, "M"},
			{EQUALS, "="},
			{COMP, "D"},
			{AT, "@"},
			{CONSTANT, "14"},
			{COMP, "0"},
			{SEMICOLON, ";"},
			{JUMP, "JMP"},
		}},

		{"RectL", "../rect/RectL.asm", []Lexeme{
			{AT, "@"},
			{CONSTANT, "0"},
			{DEST, "D"},
			{EQUALS, "="},
			{COMP, "M"},
			{AT, "@"},
			{CONSTANT, "23"},
			{COMP, "D"},
			{SEMICOLON, SEMICOLON.String()},
			{JUMP, "JLE"},
			{AT, AT.String()},
			{CONSTANT, "16"},
			{DEST, "M"},
			{EQUALS, EQUALS.String()},
			{COMP, "D"},
			{AT, AT.String()},
			{CONSTANT, "16384"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "A"},
			{AT, "@"},
			{CONSTANT, "17"},
			{DEST, "M"},
			{EQUALS, EQUALS.String()},
			{COMP, "D"},
			{AT, "@"},
			{CONSTANT, "17"},
			{DEST, "A"},
			{EQUALS, "="},
			{COMP, "M"},
			{DEST, "M"},
			{EQUALS, "="},
			{COMP, "-1"},
			{AT, AT.String()},
			{CONSTANT, "17"},
			{DEST, "D"},
			{EQUALS, "="},
			{COMP, "M"},
			{AT, AT.String()},
			{CONSTANT, "32"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "D+A"},
			{AT, AT.String()},
			{CONSTANT, "17"},
			{DEST, "M"},
			{EQUALS, EQUALS.String()},
			{COMP, "D"},
			{AT, AT.String()},
			{CONSTANT, "16"},
			{DEST, "MD"},
			{EQUALS, EQUALS.String()},
			{COMP, "M-1"},
			{AT, AT.String()},
			{CONSTANT, "10"},
			{COMP, "D"},
			{SEMICOLON, SEMICOLON.String()},
			{JUMP, "JGT"},
			{AT, AT.String()},
			{CONSTANT, "23"},
			{COMP, "0"},
			{SEMICOLON, SEMICOLON.String()},
			{JUMP, "JMP"},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFile, err := filepath.Abs(test.fpath)
			if err != nil {
				t.Fatal(err)
			}
			f, err := os.Open(testFile)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			l := NewLexer(f)
			i := 0
			for token, value := l.NextToken(); token != EOF; {
				if token != test.expected[i].token || value != test.expected[i].value {
					t.Errorf("comparison failed at token %d. expected token: %v value: %v got %v %v", i, test.expected[i].token, test.expected[i].value, token, value)
				}
				token, value = l.NextToken()
				i += 1
			}
		})
	}

}

func TestSymbols(t *testing.T) {
	tests := []struct {
		name     string
		testFile string
		expected []Lexeme
	}{
		{"Max", "../max/Max.asm", []Lexeme{
			{AT, AT.String()},
			{SYMBOL, "R0"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "M"},
			{AT, AT.String()},
			{SYMBOL, "R1"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "D-M"},
			{AT, AT.String()},
			{SYMBOL, "OUTPUT_FIRST"},
			{COMP, "D"},
			{SEMICOLON, SEMICOLON.String()},
			{JUMP, "JGT"},
			{AT, AT.String()},
			{SYMBOL, "R1"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "M"},
			{AT, AT.String()},
			{SYMBOL, "OUTPUT_D"},
			{COMP, "0"},
			{SEMICOLON, SEMICOLON.String()},
			{JUMP, "JMP"},
			{LEFT_PAREN, LEFT_PAREN.String()},
			{LABEL, "OUTPUT_FIRST"},
			{RIGHT_PAREN, RIGHT_PAREN.String()},
			{AT, AT.String()},
			{SYMBOL, "R0"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "M"},
			{LEFT_PAREN, LEFT_PAREN.String()},
			{LABEL, "OUTPUT_D"},
			{RIGHT_PAREN, RIGHT_PAREN.String()},
			{AT, AT.String()},
			{SYMBOL, "R2"},
			{DEST, "M"},
			{EQUALS, EQUALS.String()},
			{COMP, "D"},
			{LEFT_PAREN, LEFT_PAREN.String()},
			{LABEL, "INFINITE_LOOP"},
			{RIGHT_PAREN, RIGHT_PAREN.String()},
			{AT, AT.String()},
			{SYMBOL, "INFINITE_LOOP"},
			{COMP, "0"},
			{SEMICOLON, SEMICOLON.String()},
			{JUMP, "JMP"},
		}},
		{"Rect", "../rect/Rect.asm", []Lexeme{
			{AT, AT.String()},
			{CONSTANT, "0"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "M"},
			{AT, AT.String()},
			{SYMBOL, "INFINITE_LOOP"},
			{COMP, "D"},
			{SEMICOLON, SEMICOLON.String()},
			{JUMP, "JLE"},
			{AT, AT.String()},
			{SYMBOL, "counter"},
			{DEST, "M"},
			{EQUALS, EQUALS.String()},
			{COMP, "D"},
			{AT, AT.String()},
			{SYMBOL, "SCREEN"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "A"},
			{AT, AT.String()},
			{SYMBOL, "address"},
			{DEST, "M"},
			{EQUALS, EQUALS.String()},
			{COMP, "D"},

			{LEFT_PAREN, LEFT_PAREN.String()},
			{LABEL, "LOOP"},
			{RIGHT_PAREN, RIGHT_PAREN.String()},

			{AT, AT.String()},
			{SYMBOL, "address"},
			{DEST, "A"},
			{EQUALS, EQUALS.String()},
			{COMP, "M"},
			{DEST, "M"},
			{EQUALS, EQUALS.String()},
			{COMP, "-1"},
			{AT, AT.String()},
			{SYMBOL, "address"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "M"},
			{AT, AT.String()},
			{CONSTANT, "32"},
			{DEST, "D"},
			{EQUALS, EQUALS.String()},
			{COMP, "D+A"},
			{AT, AT.String()},
			{SYMBOL, "address"},
			{DEST, "M"},
			{EQUALS, EQUALS.String()},
			{COMP, "D"},
			{AT, AT.String()},
			{SYMBOL, "counter"},
			{DEST, "MD"},
			{EQUALS, EQUALS.String()},
			{COMP, "M-1"},
			{AT, AT.String()},
			{SYMBOL, "LOOP"},
			{COMP, "D"},
			{SEMICOLON, SEMICOLON.String()},
			{JUMP, "JGT"},
			{LEFT_PAREN, LEFT_PAREN.String()},
			{LABEL, "INFINITE_LOOP"},
			{RIGHT_PAREN, RIGHT_PAREN.String()},
			{AT, AT.String()},
			{SYMBOL, "INFINITE_LOOP"},
			{COMP, "0"},
			{SEMICOLON, SEMICOLON.String()},
			{JUMP, "JMP"},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFile, err := filepath.Abs(test.testFile)
			if err != nil {
				t.Fatal(err)
			}

			f, err := os.Open(testFile)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			l := NewLexer(f)
			i := 0
			for token, value := l.NextToken(); token != EOF; {
				if token != test.expected[i].token || value != test.expected[i].value {
					t.Errorf("comparison failed at token %d. expected token: %v value: %v got %v %v", i, test.expected[i].token, test.expected[i].value, token, value)
				}
				token, value = l.NextToken()
				i += 1
			}
		})
	}
}
