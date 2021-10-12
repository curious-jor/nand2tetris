package main

import (
	"VMtranslator/codewriter"
	"VMtranslator/parser"
	"fmt"
	"os"
	"path/filepath"
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
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
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
		fmt.Printf("%s is a file.\nCreating output: %s\n", path, path)

		codeWriter := codewriter.NewCodeWriter(outputFile)

		for p.HasMoreCommands() {
			err := p.Advance()
			if err != nil {
				fmt.Println(err)
			}

			ct := p.CommandType()
			switch ct {
			case parser.C_ARITHMETIC:
				codeWriter.WriteArithmetic(p.Arg1())

			case parser.C_PUSH, parser.C_POP:
				codeWriter.WritePushPop(p.CommandType(), p.Arg1(), p.Arg2())

			case parser.C_IF:
				codeWriter.WriteIf(p.Arg1())

			case parser.C_LABEL:
				codeWriter.WriteLabel(p.Arg1())

			case parser.C_GOTO:
				codeWriter.WriteGoto(p.Arg1())

			case parser.C_CALL:
				codeWriter.WriteCall(p.Arg1(), p.Arg2())

			case parser.C_FUNCTION:
				codeWriter.WriteFunction(p.Arg1(), p.Arg2())

			case parser.C_RETURN:
				codeWriter.WriteReturn()
			}
		}

		if err := codeWriter.Close(); err != nil {
			return err
		}
		fmt.Printf("Created output file: %s\n", outputFilename)
	}

	if fi.IsDir() {
		outputFilename := path + filepath.Base(path) + ".asm"
		outputFile, err := os.Create(outputFilename)
		if err != nil {
			return err
		}
		cw := codewriter.NewCodeWriter(outputFile)

		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		fmt.Printf("%s is a directory\n", path)

		for _, file := range files {
			fname := file.Name()
			if filepath.Ext(fname) == ".vm" {
				fmt.Printf("Found vm file %s. Translating...\n", fname)
				cw.SetFileName(fname)

				f, err := os.Open(path + fname)
				if err != nil {
					return err
				}
				cw.WriteInit()

				p := parser.NewParser(f)
				for p.HasMoreCommands() {
					err := p.Advance()
					if err != nil {
						fmt.Println(err)
					}

					ct := p.CommandType()
					switch ct {
					case parser.C_ARITHMETIC:
						cw.WriteArithmetic(p.Arg1())

					case parser.C_PUSH, parser.C_POP:
						cw.WritePushPop(p.CommandType(), p.Arg1(), p.Arg2())

					case parser.C_IF:
						cw.WriteIf(p.Arg1())

					case parser.C_LABEL:
						cw.WriteLabel(p.Arg1())

					case parser.C_GOTO:
						cw.WriteGoto(p.Arg1())

					case parser.C_FUNCTION:
						cw.WriteFunction(p.Arg1(), p.Arg2())

					case parser.C_CALL:
						cw.WriteCall(p.Arg1(), p.Arg2())

					case parser.C_RETURN:
						cw.WriteReturn()
					}
				}

				if err := f.Close(); err != nil {
					return err
				}
			}
		}

		if err := cw.Close(); err != nil {
			return err
		}

		fmt.Printf("Created output file: %s\n", outputFilename)
	}

	return nil
}
