// "Parser: Encapsulates access to the input code. Reads an assembly language
// command, parses it, and provides convenient access to the command’s components
// (fields and symbols). In addition, removes all white space and comments."
package parser

import (
	"assembler/lexer"
	"fmt"
	"log"
	"os"
)

type Lexeme struct {
	token lexer.Token
	value string
}

type Command interface {
	Type() Command
	Symbol() string
	Dest() string
	Comp() string
	Jump() string
}

type A_COMMAND struct {
	symbol string
}

func (c A_COMMAND) Type() Command {
	return A_COMMAND{}
}

func (c A_COMMAND) Symbol() string {
	return c.symbol
}

func (c A_COMMAND) Dest() string {
	return ""
}

func (c A_COMMAND) Comp() string {
	return ""
}

func (c A_COMMAND) Jump() string {
	return ""
}

type L_COMMAND struct {
	symbol string
}

func (c L_COMMAND) Type() Command {
	return L_COMMAND{}
}

func (c L_COMMAND) Symbol() string {
	return c.symbol
}

func (c L_COMMAND) Dest() string {
	return ""
}

func (c L_COMMAND) Comp() string {
	return ""
}

func (c L_COMMAND) Jump() string {
	return ""
}

type C_COMMAND struct {
	dest string
	comp string
	jump string
}

func (c C_COMMAND) Type() Command {
	return C_COMMAND{}
}

func (c C_COMMAND) Symbol() string {
	return ""
}

func (c C_COMMAND) Dest() string {
	return c.dest
}

func (c C_COMMAND) Comp() string {
	return c.comp
}

func (c C_COMMAND) Jump() string {
	return c.jump
}

type Parser struct {
	file     *os.File
	command  Command // The current command pointed to by the parser.
	lxr      *lexer.Lexer
	lexeme   *Lexeme
	tokenNum int
}

func NewParser(f *os.File) *Parser {
	p := new(Parser)
	p.file = f
	p.command = nil // Initially there is no command.
	p.lxr = lexer.NewLexer(f)

	t, v := p.lxr.NextToken()
	p.lexeme = &Lexeme{token: t, value: v}
	p.tokenNum = 0

	return p
}

func (p *Parser) nextToken() (lexer.Token, string) {
	tok, val := p.lxr.NextToken()
	return tok, val
}

func (p *Parser) HasMoreCommands() bool {
	return p.lxr.HasMoreTokens()
}

func (p *Parser) parseA_Command() (A_COMMAND, error) {
	token, val := p.lxr.NextToken()
	if token != lexer.CONSTANT && token != lexer.SYMBOL {
		return A_COMMAND{}, fmt.Errorf("expected CONSTANT or SYMBOL token while parsing A_COMMAND got %s", token.String())
	}
	return A_COMMAND{symbol: val}, nil

}

func (p *Parser) parseC_Command() (C_COMMAND, error) {
	if p.lexeme.token == lexer.DEST { // dest=comp
		dest := p.lexeme.value
		if t, _ := p.nextToken(); t != lexer.EQUALS { // consume '='
			return C_COMMAND{}, fmt.Errorf("expected EQUALS token got %s", t.String())
		}
		token, comp := p.nextToken() // consume comp

		if token != lexer.COMP {
			return C_COMMAND{}, fmt.Errorf("expected COMP token got %s", token)
		}

		return C_COMMAND{dest: dest, comp: comp, jump: "null"}, nil
	}
	if p.lexeme.token == lexer.COMP { // comp;jump
		comp := p.lexeme.value
		if t, _ := p.nextToken(); t != lexer.SEMICOLON { // consume ';'
			return C_COMMAND{}, fmt.Errorf("expected SEMICOLON token got %s", t.String())
		}
		token, jump := p.lxr.NextToken() // consume jump

		if token != lexer.JUMP {
			return C_COMMAND{}, fmt.Errorf("expected JUMP token got %s", token.String())
		}

		return C_COMMAND{dest: "null", comp: comp, jump: jump}, nil
	}

	return C_COMMAND{}, fmt.Errorf("attempted to parse invalid C_COMMAND format with token, val: %s, %s", p.lexeme.token.String(), p.lexeme.value)
}

func (p *Parser) parseL_Command() (L_COMMAND, error) {
	token, val := p.lxr.NextToken()
	if token != lexer.LABEL {
		return L_COMMAND{}, fmt.Errorf("expected LABEL token while parsing L_COMMAND got: %s", token.String())
	}

	if t, _ := p.lxr.NextToken(); t != lexer.RIGHT_PAREN { // consume ')'
		return L_COMMAND{}, fmt.Errorf("expected RIGHT_PAREN token while parsing L_COMMAND got: %s", t.String())
	}
	return L_COMMAND{symbol: val}, nil
}

func (p *Parser) Advance() error {
	// Reads the next command from the input and makes it the current command.
	// Should be called only if hasMoreCommands() is true.
	// Initially there is no current command.
	if !p.HasMoreCommands() {
		return fmt.Errorf("attempted to advance a Parser with no more commands")
	}

	var command Command
	var err error

	switch p.lexeme.token {
	case lexer.EOF:
		return nil
	case lexer.AT:
		{
			command, err = p.parseA_Command()
		}

	case lexer.DEST, lexer.COMP:
		{
			command, err = p.parseC_Command()
		}
	case lexer.LEFT_PAREN:
		{
			command, err = p.parseL_Command()
		}
	default:
		return fmt.Errorf("failed to parse token: %s as command", p.lexeme.token.String())
	}

	if err != nil {
		log.Printf("parsing failed at token %d (%v,%q) with error %v", p.tokenNum, p.lexeme.token, p.lexeme.value, err)
	}
	p.command = command

	// Update the parser with the next lexeme
	t, v := p.nextToken()
	p.lexeme = &Lexeme{token: t, value: v}
	p.tokenNum += 1

	return err
}

func (p *Parser) CommandType() Command {
	if p.command == nil {
		return nil
	}
	return p.command.Type()
}

func (p *Parser) Symbol() (string, error) {
	ct := p.CommandType()
	if ct != (A_COMMAND{}) && ct != (L_COMMAND{}) {
		return "", fmt.Errorf("cannot retrieve the symbol of a non-A or non-L Command")
	}
	return p.command.Symbol(), nil
}

func (p *Parser) Dest() (string, error) {
	if ct := p.CommandType(); ct != (C_COMMAND{}) {
		return "", fmt.Errorf("cannot retrieve dest mnemonic of non-C command")
	}
	return p.command.Dest(), nil
}

func (p *Parser) Comp() (string, error) {
	if ct := p.CommandType(); ct != (C_COMMAND{}) {
		return "", fmt.Errorf("cannot retrieve comp mnemonic of non-C command")
	}
	return p.command.Comp(), nil
}

func (p *Parser) Jump() (string, error) {
	if ct := p.CommandType(); ct != (C_COMMAND{}) {
		return "", fmt.Errorf("cannot retrieve jump mnemonic of non-C command")
	}
	return p.command.Jump(), nil

}

func (p *Parser) String() string {
	return "Parser for " + p.file.Name()
}

func (p *Parser) Reset() error {
	_, err := p.file.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	p.command = nil

	t, v := p.lxr.NextToken()
	p.lexeme = &Lexeme{token: t, value: v}
	p.tokenNum = 0

	return nil
}
