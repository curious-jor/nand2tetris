package main

import (
	"VMtranslator/codewriter"
	"VMtranslator/parser"
	"fmt"
	"os"
	"strings"
)

func main() {
	var srcPath string
	if len(os.Args) != 2 {
		fmt.Println("VMTranslator expects a .vm file or dir containing .vm files")
		os.Exit(1)
	}
	srcPath = os.Args[1]

	srcFile, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}

	p := parser.NewParser(srcFile)
	outputFilename := strings.Split(srcPath, ".")[0] + ".asm"

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		panic(err)
	}

	codeWriter, err := codewriter.NewCodeWriter(outputFile)
	if err != nil {
		panic(err)
	}
	defer codeWriter.Close()

	for p.HasMoreCommands() {
		p.Advance()

		if p.CommandType() == parser.C_ARITHMETIC {
			codeWriter.WriteArithmetic(p.Arg1())
		}

		if p.CommandType() == parser.C_PUSH {
			codeWriter.WritePushPop(p.CommandType(), p.Arg1(), p.Arg2())
		}
	}

}
