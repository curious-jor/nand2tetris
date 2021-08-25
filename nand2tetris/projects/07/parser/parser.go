package parser

import (
	"VMtranslator/lexer"
	"bufio"
	"fmt"
	"log"
	"os"
)

type lexeme struct {
	token lexer.Token
	value string
}

type CommandType int

const (
	C_ARITHMETIC CommandType = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

type Command struct {
	ct   CommandType // the type of command
	arg1 string
	arg2 string
}

type Parser struct {
	f      *os.File
	rdr    *bufio.Reader
	lxr    *lexer.Lexer
	lexeme *lexeme
	cmd    *Command
}

func NewParser(file *os.File) *Parser {
	reader := bufio.NewReader(file)
	lexer := lexer.NewLexer(file)

	// Load the first token from the input
	tok, val := lexer.NextToken()
	initialLex := lexeme{token: tok, value: val}
	return &Parser{
		f:      file,
		rdr:    reader,
		lxr:    lexer,
		lexeme: &initialLex,
	}
}

func (p *Parser) HasMoreCommands() bool {
	return p.lxr.HasMoreTokens()
}

func isArithmeticCommand(command string) bool {
	switch command {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return true
	default:
		return false
	}
}

func (p *Parser) parsePushCommand() *Command {
	token, segment := p.lxr.NextToken() // consume segment
	if token != lexer.ARG {
		fmt.Printf("expected ARG token while parsing \"push\" command got %s instead\n", token)
		return &Command{ct: C_PUSH, arg1: segment, arg2: ""}
	}
	token, index := p.lxr.NextToken() // consume index
	if token != lexer.ARG {
		fmt.Printf("expected ARG token while parsing \"push\" command got %s instead\n", token)
		return &Command{ct: C_PUSH, arg1: segment, arg2: index}
	}

	return &Command{ct: C_PUSH, arg1: segment, arg2: index}
}

func (p *Parser) Advance() {
	if !p.HasMoreCommands() {
		log.Panicln("attempted to advance a parser with no commands left")
		return
	}

	switch p.lexeme.token {
	case lexer.COMMAND:
		{
			if isArithmeticCommand(p.lexeme.value) {
				p.cmd = &Command{ct: C_ARITHMETIC, arg1: p.lexeme.value, arg2: ""}
			}

			if p.lexeme.value == "push" {
				p.cmd = p.parsePushCommand()
			}
		}
	}

	// Update parser with next lexeme
	t, v := p.lxr.NextToken()
	p.lexeme = &lexeme{token: t, value: v}
}
