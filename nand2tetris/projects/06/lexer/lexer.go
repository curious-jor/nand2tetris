package lexer

import (
	"bufio"
	"io"
	"log"
	"os"
)

type Token int

const (
	// Special characters
	EOF         Token = iota
	AT                // @
	EQUALS            // =
	SEMICOLON         // ;
	BACKSLASH         // /
	LEFT_PAREN        // (
	RIGHT_PAREN       // )

	VALUE
	CONSTANT
	DEST
	COMP
	JUMP
	SYMBOL
	LABEL
)

var tokens = []string{
	EOF:       "EOF",
	AT:        "@",
	EQUALS:    "=",
	SEMICOLON: ";",
	VALUE:     "VALUE",

	CONSTANT:    "CONSTANT",
	DEST:        "DEST",
	COMP:        "COMP",
	JUMP:        "JUMP",
	LEFT_PAREN:  "(",
	RIGHT_PAREN: ")",
	SYMBOL:      "SYMBOL",
	LABEL:       "LABEL",
}

func (t Token) String() string {
	return tokens[t]
}

type Lexer struct {
	r    *bufio.Reader
	prev Token
}

func NewLexer(f *os.File) *Lexer {
	return &Lexer{
		r: bufio.NewReader(f),
	}
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\r' || ch == '\n' || ch == '\t'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isSymbolOnlyChar(ch rune) bool {
	return ch == '_' || ch == '.' || ch == '$' || ch == ':'
}

func isCompOnlyChar(ch rune) bool {
	return ch == '!' || ch == '-'
}

var eofRune = rune(0)

func (l *Lexer) getChar() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return eofRune
	}
	return ch
}

func (l *Lexer) unread() {
	l.r.UnreadRune()
}

// NextToken returns the next token and its value as a string from the input stream
func (l *Lexer) NextToken() (Token, string) {
	var lastChar = ' '

	// Skip whitespace
	for isWhiteSpace(lastChar) {
		lastChar = l.getChar()
	}

	switch lastChar {
	case '@':
		l.prev = AT
		return AT, AT.String()
	case '=':
		l.prev = EQUALS
		return EQUALS, EQUALS.String()
	case ';':
		l.prev = SEMICOLON
		return SEMICOLON, SEMICOLON.String()
	case '(':
		l.prev = LEFT_PAREN
		return LEFT_PAREN, LEFT_PAREN.String()
	case ')':
		l.prev = RIGHT_PAREN
		return RIGHT_PAREN, RIGHT_PAREN.String()
	}

	// If the beginning character is a digit, attempt to tokenize as a constant or comp mnemonic
	if isDigit(lastChar) {
		charSeq := []rune{lastChar}
		for lastChar = l.getChar(); lastChar != ';' && !isWhiteSpace(lastChar); {
			charSeq = append(charSeq, lastChar)
			lastChar = l.getChar()
		}

		l.unread()

		if l.prev == AT {
			l.prev = CONSTANT
			return CONSTANT, string(charSeq)
		}

		if lastChar == ';' || l.prev == EQUALS {
			l.prev = COMP
			return COMP, string(charSeq)
		}

		log.Printf("could not tokenize sequence %q as constant or comp with prev char %q", string(charSeq), l.prev.String())
		return VALUE, string(charSeq)

	} else if isCompOnlyChar(lastChar) {
		charSeq := []rune{lastChar}
		for lastChar = l.getChar(); lastChar != ';' && !isWhiteSpace(lastChar); {
			charSeq = append(charSeq, lastChar)
			lastChar = l.getChar()
		}

		l.unread()

		if lastChar == ';' || l.prev == EQUALS {
			l.prev = COMP
			return COMP, string(charSeq)
		}

		log.Printf("could not tokenize sequence %q as comp", string(charSeq))
		l.prev = VALUE
		return VALUE, string(charSeq)
	} else if isLetter(lastChar) || isSymbolOnlyChar(lastChar) {
		// any char sequence that doesn't begin with a digit can be a symbol, label, or dest/comp/jump mnemonic
		charSeq := []rune{lastChar}
		for lastChar = l.getChar(); lastChar != ';' && lastChar != '=' && lastChar != ')' && !isWhiteSpace(lastChar); {
			charSeq = append(charSeq, lastChar)
			lastChar = l.getChar()
		}

		// Put right parenthesis or whitespace back into input stream for a separate token
		l.unread()

		if l.prev == AT {
			l.prev = SYMBOL
			return SYMBOL, string(charSeq)
		}

		if lastChar == ')' {
			l.prev = LABEL
			return LABEL, string(charSeq)
		}

		if lastChar == '=' {
			l.prev = DEST
			return DEST, string(charSeq)
		}

		if lastChar == ';' || l.prev == EQUALS { // dest=comp or comp;jump
			l.prev = COMP
			return COMP, string(charSeq)
		}

		if l.prev == SEMICOLON { // comp;jump
			l.prev = JUMP
			return JUMP, string(charSeq)
		}

		log.Printf("could not tokenize sequence %q as symbol, label, or mnemonic", string(charSeq))
		l.prev = VALUE
		return VALUE, string(charSeq)

	}

	if lastChar == eofRune {
		l.prev = EOF
		return EOF, EOF.String()
	}

	if lastChar == '/' {
		// Skip comment and return next token
		for lastChar != eofRune && lastChar != '\n' && lastChar != '\r' {
			lastChar = l.getChar()
		}

		if lastChar != eofRune {
			return l.NextToken()
		}
	}

	log.Printf("could not derive current token %q", lastChar)
	thisChar := lastChar
	l.prev = Token(thisChar)
	_ = l.getChar()
	return Token(thisChar), string(thisChar)
}

func (l *Lexer) HasMoreTokens() bool {
	_, err := l.r.Peek(1)
	if err != nil && err != io.EOF {
		panic(err)
	}
	return err == nil
}

func (l *Lexer) Peek(n int) ([]byte, error) {
	return l.r.Peek(n)
}

func (l *Lexer) ReadBytes(delim byte) ([]byte, error) {
	return l.r.ReadBytes(delim)
}
