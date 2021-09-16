package parser

import (
	"VMtranslator/lexer"
	"bufio"
	"errors"
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

func (p *Parser) parseArithmeticCommand() (*Command, error) {
	return &Command{ct: C_ARITHMETIC, arg1: p.lexeme.Value, arg2: emptyArg2}, nil
}

func (p *Parser) parsePushPopCommand() (*Command, error) {
	var currCmdType CommandType
	if p.lexeme.Value == "push" {
		currCmdType = C_PUSH
	}
	if p.lexeme.Value == "pop" {
		currCmdType = C_POP
	}

	segment := p.lxr.NextToken() // consume segment
	if segment.Token != lexer.ARG {
		return &Command{
				ct: currCmdType, arg1: segment.Value,
				arg2: emptyArg2,
			}, &ParserCouldNotParseError{
				lxm: segment,
				msg: fmt.Sprintf("expected ARG token while parsing \"push\" command got %s instead\n", segment.Token),
			}
	}

	index := p.lxr.NextToken() // consume index
	if index.Token != lexer.ARG {
		fmt.Printf("expected ARG token while parsing \"push\" command got %s instead\n", index.Token)
	}

	indexInt, err := strconv.Atoi(index.Value)
	if err != nil {
		return &Command{
				ct:   currCmdType,
				arg1: segment.Value,
				arg2: emptyArg2,
			}, &ParserCouldNotParseError{
				lxm: index,
				msg: fmt.Sprintf("could not convert %q to int while parsing push command (%s, %q)", index, index.Token.String(), index),
			}
	}

	return &Command{ct: currCmdType, arg1: segment.Value, arg2: indexInt}, nil
}

var ErrParserNoMoreCommands = errors.New("parser has no more commands")

type ParserCouldNotParseError struct {
	lxm *lexer.Lexeme
	msg string
}

func (e *ParserCouldNotParseError) Error() string {
	return e.msg
}

func (p *Parser) Advance() error {
	if !p.HasMoreCommands() {
		return ErrParserNoMoreCommands
	}

	var parsedCmd *Command
	var err error
	switch p.lexeme.Token {
	case lexer.COMMAND:
		{
			switch p.lexeme.Value {
			case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
				{
					parsedCmd, err = p.parseArithmeticCommand()
				}
			case "push", "pop":
				{
					parsedCmd, err = p.parsePushPopCommand()
				}
			default:
				{
					parsedCmd = nil
					err = fmt.Errorf("attempted to parse unsupported command: %q as Commmand", p.lexeme.Value)
				}
			}
		}
	default:
		{
			err = fmt.Errorf("attempted to parse non-Command token (%s, %q) as a non-terminal", p.lexeme.Token.String(), p.lexeme.Value)
		}
	}
	p.cmd = parsedCmd

	// Update parser with next lexeme
	p.lexeme = p.lxr.NextToken()
	return err
}

var emptyCommandType = CommandType(-1)

func (p *Parser) CommandType() CommandType {
	if p.cmd == nil {
		return emptyCommandType
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
		return emptyArg2
	}

	return p.cmd.arg2
}
