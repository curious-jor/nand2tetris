package lexer

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
)

type Token int

const (
	EOF Token = iota
	ILLEGAL
	BACKSLASH
	NEWLINE

	COMMAND
	ARG
)

var tokens = []string{
	EOF:       "EOF",
	ILLEGAL:   "ILLEGAL",
	BACKSLASH: "/",
	NEWLINE:   "\\n",
	COMMAND:   "COMMAND",
	ARG:       "ARG",
}

func (t Token) String() string {
	return tokens[t]
}

type Lexer struct {
	r    *bufio.Reader
	prev Token
	fp   FilePosition
}

type FilePosition struct {
	Line int
	Col  int
}

func NewLexer(f *os.File) *Lexer {
	return &Lexer{
		r:    bufio.NewReader(f),
		prev: NEWLINE,
		fp:   FilePosition{Line: 1, Col: 1},
	}
}

const eofRune = rune(0)
const newlineRune = rune(10)

func (l *Lexer) getChar() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return eofRune
		} else {
			panic(err)
		}
	}
	return ch
}

func (l *Lexer) unread() error {
	return l.r.UnreadRune()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

type Lexeme struct {
	Token Token
	Value string
}

func (l *Lexer) NextToken() (FilePosition, *Lexeme) {
	var currChar = ' '
	for isWhitespace(currChar) {
		currChar = l.getChar()
		l.fp.Col += 1
	}

	if currChar == eofRune {
		l.prev = EOF
		return l.fp, &Lexeme{EOF, EOF.String()}
	}

	if currChar == newlineRune {
		l.prev = NEWLINE
		l.fp.Line += 1
		l.fp.Col = 1
		return l.NextToken()
	}

	if isLetter(currChar) || isDigit(currChar) {
		startingPos := FilePosition{Line: l.fp.Line, Col: l.fp.Col - 1}
		charSeq := []rune{currChar}
		for currChar = l.getChar(); !isWhitespace(currChar); {
			charSeq = append(charSeq, currChar)
			currChar = l.getChar()
			l.fp.Col += 1
		}

		if err := l.unread(); err != nil {
			panic(err)
		}

		if _, err := strconv.ParseInt(string(charSeq), 10, 16); err == nil {
			l.prev = ARG
			return startingPos, &Lexeme{ARG, string(charSeq)}
		}

		if l.prev == NEWLINE {
			l.prev = COMMAND
			return startingPos, &Lexeme{COMMAND, string(charSeq)}
		}

		if l.prev == COMMAND || l.prev == ARG {
			l.prev = ARG
			return startingPos, &Lexeme{ARG, string(charSeq)}
		}
	}

	if currChar == '/' {
		for currChar != eofRune && currChar != newlineRune && currChar != '\r' {
			currChar = l.getChar()
		}

		if currChar != eofRune {
			return l.NextToken()
		}
	}

	log.Printf("could not derive token from char: %q with prev token %s", currChar, l.prev.String())
	thisChar := currChar
	l.prev = ILLEGAL
	_ = l.getChar()
	l.fp.Col += 1
	return l.fp, &Lexeme{ILLEGAL, string(thisChar)}
}

func (l *Lexer) HasMoreTokens() bool {
	nextByte := 1
	_, err := l.r.Peek(nextByte)
	if err != nil && err != io.EOF {
		panic(err)
	}
	return err == nil
}

func (a *Lexeme) Equals(b *Lexeme) bool {
	return a.Token == b.Token && a.Value == b.Value
}
