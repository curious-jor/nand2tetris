package code

import (
	"bytes"
	"testing"
)

func TestDest(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Not legal mnemonic", "B", ""},
		{"Null string", "null", "000"},
		{"Dest M", "M", "001"},
		{"Dest D", "D", "010"},
		{"Dest MD", "MD", "011"},
		{"Dest A", "A", "100"},
		{"Dest AM", "AM", "101"},
		{"Dest AD", "AD", "110"},
		{"Dest AMD", "AMD", "111"},
	}

	for _, test := range tests {
		actual := BytesToBitString(Dest(test.input))
		if test.expected != actual {
			t.Errorf("expected bit string %q for input %q but got %q", test.expected, test.input, actual)
		}
	}
}

func TestComp(t *testing.T) {
	// Map of mnemonics mapped to their pre-defined bits.
	// Use to test function one by one in a loop.
	mnemonicBits := map[string][]byte{
		"0":   {0, 1, 0, 1, 0, 1, 0},
		"1":   {0, 1, 1, 1, 1, 1, 1},
		"-1":  {0, 1, 1, 1, 0, 1, 0},
		"D":   {0, 0, 0, 1, 1, 0, 0},
		"A":   {0, 1, 1, 0, 0, 0, 0},
		"M":   {1, 1, 1, 0, 0, 0, 0},
		"!D":  {0, 0, 0, 1, 1, 0, 1},
		"!A":  {0, 1, 1, 0, 0, 0, 1},
		"!M":  {1, 1, 1, 0, 0, 0, 1},
		"-D":  {0, 0, 0, 1, 1, 1, 1},
		"-A":  {0, 1, 1, 0, 0, 1, 1},
		"-M":  {1, 1, 1, 0, 0, 1, 1},
		"D+1": {0, 0, 1, 1, 1, 1, 1},
		"A+1": {0, 1, 1, 0, 1, 1, 1},
		"M+1": {1, 1, 1, 0, 1, 1, 1},
		"D-1": {0, 0, 0, 1, 1, 1, 0},
		"A-1": {0, 1, 1, 0, 0, 1, 0},
		"M-1": {1, 1, 1, 0, 0, 1, 0},
		"D+A": {0, 0, 0, 0, 0, 1, 0},
		"D+M": {1, 0, 0, 0, 0, 1, 0},
		"D-A": {0, 0, 1, 0, 0, 1, 1},
		"D-M": {1, 0, 1, 0, 0, 1, 1},
		"A-D": {0, 0, 0, 0, 1, 1, 1},
		"M-D": {1, 0, 0, 0, 1, 1, 1},
		"D&A": {0, 0, 0, 0, 0, 0, 0},
		"D&M": {1, 0, 0, 0, 0, 0, 0},
		"D|A": {0, 0, 1, 0, 1, 0, 1},
		"D|M": {1, 0, 1, 0, 1, 0, 1},
	}

	for mnemonic, expectedBits := range mnemonicBits {
		comp := Comp(mnemonic)
		compString := BytesToBitString(comp)
		expectedString := BytesToBitString(expectedBits)
		if !bytes.Equal(comp, expectedBits) {
			t.Errorf("got %s bits for mnemonic %s but expected %s", compString, mnemonic, expectedString)
		}
	}
}

func TestJump(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Not legal mnemonic", "B", ""},
		{"Null string", "null", "000"},
		{"JGT", "JGT", "001"},
		{"JEQ", "JEQ", "010"},
		{"JGE", "JGE", "011"},
		{"JLT", "JLT", "100"},
		{"JNE", "JNE", "101"},
		{"JLE", "JLE", "110"},
		{"JMP", "JMP", "111"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			actual := BytesToBitString(Jump(test.input))
			if test.expected != actual {
				t.Errorf("expected %q for input %q but got %q", test.expected, test.input, actual)
			}
		})
	}
}
