// "Code: Translates Hack assembly language mnemonics into binary codes."
package code

import (
	"strconv"
	"strings"
)

func Dest(mnemonic string) []byte {
	var bits []byte
	switch mnemonic {
	case "null":
		bits = []byte{0, 0, 0}
	case "M":
		bits = []byte{0, 0, 1}
	case "D":
		bits = []byte{0, 1, 0}
	case "MD":
		bits = []byte{0, 1, 1}
	case "A":
		bits = []byte{1, 0, 0}
	case "AM":
		bits = []byte{1, 0, 1}
	case "AD":
		bits = []byte{1, 1, 0}
	case "AMD":
		bits = []byte{1, 1, 1}
	case "":
		bits = []byte{}
	default:
		bits = []byte{}
	}
	return bits
}

func Comp(mnemonic string) []byte {
	var bits []byte

	switch mnemonic {
	case "0":
		bits = []byte{0, 1, 0, 1, 0, 1, 0}
	case "1":
		bits = []byte{0, 1, 1, 1, 1, 1, 1}
	case "-1":
		bits = []byte{0, 1, 1, 1, 0, 1, 0}
	case "D":
		bits = []byte{0, 0, 0, 1, 1, 0, 0}
	case "A":
		bits = []byte{0, 1, 1, 0, 0, 0, 0}
	case "M":
		bits = []byte{1, 1, 1, 0, 0, 0, 0}
	case "!D":
		bits = []byte{0, 0, 0, 1, 1, 0, 1}
	case "!A":
		bits = []byte{0, 1, 1, 0, 0, 0, 1}
	case "!M":
		bits = []byte{1, 1, 1, 0, 0, 0, 1}
	case "-D":
		bits = []byte{0, 0, 0, 1, 1, 1, 1}
	case "-A":
		bits = []byte{0, 1, 1, 0, 0, 1, 1}
	case "-M":
		bits = []byte{1, 1, 1, 0, 0, 1, 1}
	case "D+1":
		bits = []byte{0, 0, 1, 1, 1, 1, 1}
	case "A+1":
		bits = []byte{0, 1, 1, 0, 1, 1, 1}
	case "M+1":
		bits = []byte{1, 1, 1, 0, 1, 1, 1}
	case "D-1":
		bits = []byte{0, 0, 0, 1, 1, 1, 0}
	case "A-1":
		bits = []byte{0, 1, 1, 0, 0, 1, 0}
	case "M-1":
		bits = []byte{1, 1, 1, 0, 0, 1, 0}
	case "D+A":
		bits = []byte{0, 0, 0, 0, 0, 1, 0}
	case "D+M":
		bits = []byte{1, 0, 0, 0, 0, 1, 0}
	case "D-A":
		bits = []byte{0, 0, 1, 0, 0, 1, 1}
	case "D-M":
		bits = []byte{1, 0, 1, 0, 0, 1, 1}
	case "A-D":
		bits = []byte{0, 0, 0, 0, 1, 1, 1}
	case "M-D":
		bits = []byte{1, 0, 0, 0, 1, 1, 1}
	case "D&A":
		bits = []byte{0, 0, 0, 0, 0, 0, 0}
	case "D&M":
		bits = []byte{1, 0, 0, 0, 0, 0, 0}
	case "D|A":
		bits = []byte{0, 0, 1, 0, 1, 0, 1}
	case "D|M":
		bits = []byte{1, 0, 1, 0, 1, 0, 1}
	case "":
		bits = []byte{}
	default:
		bits = []byte{}
	}

	return bits
}

func Jump(mnemonic string) []byte {
	var bits []byte
	switch mnemonic {
	case "null":
		bits = []byte{0, 0, 0}
	case "JGT":
		bits = []byte{0, 0, 1}
	case "JEQ":
		bits = []byte{0, 1, 0}
	case "JGE":
		bits = []byte{0, 1, 1}
	case "JLT":
		bits = []byte{1, 0, 0}
	case "JNE":
		bits = []byte{1, 0, 1}
	case "JLE":
		bits = []byte{1, 1, 0}
	case "JMP":
		bits = []byte{1, 1, 1}
	case "":
		bits = []byte{}
	default:
		bits = []byte{}
	}
	return bits
}

func BytesToBitString(b []byte) string {
	s := make([]string, len(b))
	for i, bit := range b {
		s[i] = strconv.Itoa(int(bit))
	}
	return strings.Join(s, "")
}
