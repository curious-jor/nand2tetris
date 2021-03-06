package main

import (
	"VMtranslator/codewriter"
	"VMtranslator/parser"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("VMTranslator expects a .vm file or dir containing .vm files")
		os.Exit(1)
	}
	srcPath := os.Args[1]

	fmt.Printf("Translating %s ...\n", srcPath)
	err := translate(srcPath)
	if err != nil {
		panic(err)
	}
}

func translate(path string) error {
	srcFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	p := parser.NewParser(srcFile)
	outputFilename := strings.Split(path, ".")[0] + ".asm"

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		panic(err)
	}

	codeWriter := codewriter.NewCodeWriter(outputFile)

	for p.HasMoreCommands() {
		p.Advance()

		if p.CommandType() == parser.C_ARITHMETIC {
			codeWriter.WriteArithmetic(p.Arg1())
		}

		if ct := p.CommandType(); ct == parser.C_PUSH || ct == parser.C_POP {
			codeWriter.WritePushPop(p.CommandType(), p.Arg1(), p.Arg2())
		}
	}

	if err := codeWriter.Close(); err != nil {
		return err
	}
	fmt.Printf("Created output file: %s", outputFilename)

	return nil
}
