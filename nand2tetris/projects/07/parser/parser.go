package parser

import (
	"VMtranslator/lexer"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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

var commandTypes = []string{
	C_ARITHMETIC: "C_ARITHMETIC",
	C_PUSH:       "C_PUSH",
	C_POP:        "C_POP",
	C_LABEL:      "C_LABEL",
	C_GOTO:       "C_GOTO",
	C_IF:         "C_IF",
	C_FUNCTION:   "C_FUNCTION",
	C_RETURN:     "C_RETURN",
	C_CALL:       "C_CALL",
}

func (ct CommandType) String() string {
	return commandTypes[ct]
}

type Command struct {
	ct   CommandType // the type of command
	arg1 string
	arg2 int
}

var emptyArg2 = -1

type Parser struct {
	f      *os.File
	rdr    *bufio.Reader
	lxr    *lexer.Lexer
	lexeme *lexer.Lexeme
	cmd    *Command
}

func NewParser(file *os.File) *Parser {
	if file == nil {
		return &Parser{}
	}

	reader := bufio.NewReader(file)
	lexer := lexer.NewLexer(file)

	// Load the first token from the input
	initialLex := lexer.NextToken()
	return &Parser{
		f:      file,
		rdr:    reader,
		lxr:    lexer,
		lexeme: initialLex,
		cmd:    nil,
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

func (p *Parser) parseArithmeticCommand() *Command {
	return &Command{ct: C_ARITHMETIC, arg1: p.lexeme.Value, arg2: emptyArg2}
}

func (p *Parser) parsePushCommand() *Command {
	segment := p.lxr.NextToken() // consume segment
	if segment.Token != lexer.ARG {
		fmt.Printf("expected ARG token while parsing \"push\" command got %s instead\n", segment.Token)
		return &Command{ct: C_PUSH, arg1: segment.Value, arg2: emptyArg2}
	}

	index := p.lxr.NextToken() // consume index
	if index.Token != lexer.ARG {
		fmt.Printf("expected ARG token while parsing \"push\" command got %s instead\n", index.Token)
	}

	indexInt, err := strconv.Atoi(index.Value)
	if err != nil {
		fmt.Printf("could not convert %q to int while parsing push command (%s, %q)", index, index.Token.String(), index)
		return &Command{ct: C_PUSH, arg1: segment.Value, arg2: emptyArg2}
	}

	return &Command{ct: C_PUSH, arg1: segment.Value, arg2: indexInt}
}

func (p *Parser) Advance() {
	if !p.HasMoreCommands() {
		fmt.Println("attempted to advance a parser with no commands left")
		return
	}

	switch p.lexeme.Token {
	case lexer.COMMAND:
		{
			if isArithmeticCommand(p.lexeme.Value) {
				p.cmd = p.parseArithmeticCommand()
			}

			if p.lexeme.Value == "push" {
				p.cmd = p.parsePushCommand()
			}
		}
	}

	// Update parser with next lexeme

	p.lexeme = p.lxr.NextToken()
}

var emptyCommandType = -1

func (p *Parser) CommandType() CommandType {
	if p.cmd == nil {
		fmt.Println("attempted to get command type of empty command")
		return CommandType(emptyCommandType)
	}
	return p.cmd.ct
}

func (p *Parser) Arg1() string {
	if p.CommandType() == C_RETURN {
		fmt.Println("attempted to call Arg1 on a C_RETURN command")
		return ""
	}

	return p.cmd.arg1
}

func commandHasArg2(ct CommandType) bool {
	return ct == C_PUSH || ct == C_POP || ct == C_FUNCTION || ct == C_CALL
}

func (p *Parser) Arg2() int {
	if curr := p.CommandType(); !commandHasArg2(curr) {
		fmt.Printf("attempted to call Arg2 on a command with the incorrect type. got %#v\n", curr)
	}

	return p.cmd.arg2
}
