package parser

import (
	"VMtranslator/lexer"
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

var emptyArg1 = ""
var emptyArg2 = -1

type Parser struct {
	f      *os.File
	lxr    *lexer.Lexer
	lexeme *lexer.Lexeme
	cmd    *Command
	fp     lexer.FilePosition
}

func NewParser(file *os.File) *Parser {
	if file == nil {
		return &Parser{}
	}
	lexer := lexer.NewLexer(file)

	// Load the first token from the input
	pos, initialLex := lexer.NextToken()
	return &Parser{
		f:      file,
		lxr:    lexer,
		lexeme: initialLex,
		cmd:    nil,
		fp:     pos,
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

	pos, segment := p.lxr.NextToken() // consume segment
	p.fp = pos
	if segment.Token != lexer.ARG {
		return &Command{
				ct: currCmdType, arg1: segment.Value,
				arg2: emptyArg2,
			}, &ParserError{
				line: p.fp.Line,
				col:  p.fp.Col,
				lxm:  segment,
				msg:  fmt.Sprintf("expected ARG token while parsing \"push\" command got %s instead\n", segment.Token),
			}
	}

	pos, index := p.lxr.NextToken() // consume index
	p.fp = pos
	if index.Token != lexer.ARG {
		fmt.Printf("expected ARG token while parsing \"push\" command got %s instead\n", index.Token)
	}

	indexInt, err := strconv.Atoi(index.Value)
	if err != nil {
		return &Command{
				ct:   currCmdType,
				arg1: segment.Value,
				arg2: emptyArg2,
			}, &ParserError{
				line: p.fp.Line,
				col:  p.fp.Col,
				lxm:  index,
				msg:  fmt.Sprintf("could not convert %q to int while parsing push command (%s, %q)", index, index.Token.String(), index),
			}
	}

	return &Command{ct: currCmdType, arg1: segment.Value, arg2: indexInt}, nil
}

var ErrParserNoMoreCommands = errors.New("parser has no more commands")

type ParserError struct {
	line int
	col  int
	lxm  *lexer.Lexeme
	msg  string
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("Error (Line: %d, Col: %d) - %s", e.line, e.col, e.msg)
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
			case "label":
				{
					// consume label
					pos, label := p.lxr.NextToken()
					p.fp = pos
					if label.Token != lexer.ARG {
						err = &ParserError{
							line: p.fp.Line,
							col:  p.fp.Col,
							lxm:  p.lexeme,
							msg:  fmt.Sprintf("expected ARG token while parsing \"label\" command but got %s", label.Token.String()),
						}
					}
					parsedCmd = &Command{ct: C_LABEL, arg1: label.Value, arg2: emptyArg2}
				}
			case "goto":
				{
					// consume label
					pos, label := p.lxr.NextToken()
					p.fp = pos
					if label.Token != lexer.ARG {
						err = &ParserError{
							line: p.fp.Line,
							col:  p.fp.Col,
							lxm:  p.lexeme,
							msg:  fmt.Sprintf("expected ARG token while parsing \"goto\" command but got %s", label.Token.String()),
						}
					}
					parsedCmd = &Command{ct: C_GOTO, arg1: label.Value, arg2: emptyArg2}
				}
			case "if-goto":
				{
					// consume label
					pos, label := p.lxr.NextToken()
					p.fp = pos
					if label.Token != lexer.ARG {
						err = &ParserError{
							line: p.fp.Line,
							col:  p.fp.Col,
							lxm:  p.lexeme,
							msg:  fmt.Sprintf("expected ARG token while parsing \"if-goto\" command but got %s", label.Token.String()),
						}
					}
					parsedCmd = &Command{ct: C_IF, arg1: label.Value, arg2: emptyArg2}
				}
			case "function":
				{
					// consume functionName
					pos, functionName := p.lxr.NextToken()
					p.fp = pos
					if functionName.Token != lexer.ARG {
						err = &ParserError{line: p.fp.Line, col: p.fp.Col, lxm: p.lexeme, msg: fmt.Sprintf("expected ARG token while parsing %q command but got %q)", p.lexeme.Value, functionName.Token.String())}
					}

					// consume numLocals
					pos, numLocals := p.lxr.NextToken()
					p.fp = pos
					if functionName.Token != lexer.ARG {
						err = &ParserError{line: p.fp.Line, col: p.fp.Col, lxm: p.lexeme, msg: fmt.Sprintf("expected ARG token while parsing %q command but got %q)", p.lexeme.Value, numLocals.Token.String())}
					}

					numLocalsInt, err := strconv.Atoi(numLocals.Value)
					if err != nil {
						parsedCmd = &Command{
							ct:   C_FUNCTION,
							arg1: functionName.Value,
							arg2: emptyArg2,
						}
						return &ParserError{
							line: p.fp.Line,
							col:  p.fp.Col,
							lxm:  p.lexeme,
							msg:  fmt.Sprintf("could not convert %q to int while parsing \"function\" %s %s", numLocals.Value, functionName.Value, numLocals.Value),
						}
					}
					parsedCmd = &Command{ct: C_FUNCTION, arg1: functionName.Value, arg2: numLocalsInt}
				}
			case "return":
				{
					parsedCmd = &Command{ct: C_RETURN, arg1: emptyArg1, arg2: emptyArg2}
				}
			default:
				{
					parsedCmd = nil
					err = &ParserError{line: p.fp.Line, col: p.fp.Col, msg: fmt.Sprintf("attempted to parse unsupported command: %q as Commmand", p.lexeme.Value)}
				}
			}
		}
	default:
		{
			err = &ParserError{line: p.fp.Line, col: p.fp.Col, msg: fmt.Sprintf("attempted to parse non-Command token (%s, %q) as a non-terminal", p.lexeme.Token.String(), p.lexeme.Value)}
		}
	}
	p.cmd = parsedCmd

	// Update parser with next lexeme
	p.fp, p.lexeme = p.lxr.NextToken()
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
		return emptyArg2
	}

	return p.cmd.arg2
}
