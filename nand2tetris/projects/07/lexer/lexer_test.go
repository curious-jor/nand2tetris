package lexer

import (
	"os"
	"testing"
)

func TestSimpleAdd(t *testing.T) {
	f, err := os.Open("../StackArithmetic/SimpleAdd/SimpleAdd.vm")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lxr := NewLexer(f)
	expected := []struct {
		token Token
		value string
	}{
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "7"},
		{COMMAND, "push"},
		{ARG, "constant"},
		{ARG, "8"},
		{COMMAND, "add"},
	}

	i := 0
	for tok, val := lxr.NextToken(); tok != EOF; {
		if expected[i].token != tok || expected[i].value != val {
			t.Errorf("comparison failed at token %d expected (%s, %q) got (%s, %q)", i, expected[i].token.String(), expected[i].value, tok.String(), val)
		}
		tok, val = lxr.NextToken()
		i += 1
	}
}
