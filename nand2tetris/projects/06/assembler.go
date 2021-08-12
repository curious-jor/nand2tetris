package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"assembler/code"
	"assembler/parser"
	"assembler/symboltable"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Assembler expects one argument: *filename*.asm")
		os.Exit(1)
	}
	filePath := os.Args[1]
	asmFile, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Could not open .asm file: %s\n", asmFile.Name())
		panic(err)
	}
	defer asmFile.Close()

	fmt.Printf("Assembling \"%s\"\n", filePath)
	p := parser.NewParser(asmFile)

	outputFileName := strings.Split(filePath, ".")[0] + ".hack"
	hackFile, err := os.Create(outputFileName)

	if err != nil {
		fmt.Printf("Error: Could not create file %s\n", outputFileName)
		panic(err)
	}

	const baseTwo = 2
	const baseTen = 10
	const sixteenBit = 16
	const newLine = "\n"

	// First pass to build symbol table
	romAddress := 0
	st := symboltable.NewSymbolTable()
	for p.HasMoreCommands() {
		if err := p.Advance(); err != nil {
			panic(err)
		}

		ct := p.CommandType()
		if ct == (parser.A_COMMAND{}) || ct == (parser.C_COMMAND{}) {
			romAddress += 1
		}
		if ct == (parser.L_COMMAND{}) {
			symbol, err := p.Symbol()
			if err != nil {
				panic(err)
			}
			st.AddEntry(symbol, romAddress)
		}
	}

	if err := p.Reset(); err != nil {
		fmt.Printf("could not reset parser after first pass got error: %v", err)
		os.Exit(1)
	}

	// Second pass
	ramAddress := 16
	for p.HasMoreCommands() {
		if err := p.Advance(); err != nil {
			panic(err)
		}

		if p.CommandType() == (parser.A_COMMAND{}) {
			symbol, err := p.Symbol()
			if err != nil {
				panic(err)
			}
			// Parse symbol into decimal representation
			symbolAsInt, err := strconv.ParseInt(symbol, baseTen, sixteenBit)

			if err != nil { // @Xxx is a symbol, not a decimal
				if address := st.GetAddress(symbol); address != -1 { // symbol is in table; replace with numeric meaning
					symbolAsInt = int64(st.GetAddress(symbol))
				} else { // symbol is a new variable
					st.AddEntry(symbol, ramAddress)
					symbolAsInt = int64(ramAddress)
					ramAddress += 1
				}
			}

			// Write A command as binary string to file
			symbolAsBinary := fmt.Sprintf("%016s", strconv.FormatInt(symbolAsInt, baseTwo))
			line := symbolAsBinary + newLine
			hackFile.WriteString(line)
		}
		if p.CommandType() == (parser.C_COMMAND{}) {
			dest, err := p.Dest()
			if err != nil {
				panic(err)
			}
			comp, err := p.Comp()
			if err != nil {
				panic(err)
			}
			jump, err := p.Jump()
			if err != nil {
				panic(err)
			}

			// convert mnemonics to bits
			destBits, compBits, jumpBits := code.Dest(dest), code.Comp(comp), code.Jump(jump)
			line := "111" + code.BytesToBitString(compBits) + code.BytesToBitString(destBits) + code.BytesToBitString(jumpBits) + newLine
			hackFile.WriteString(line)
		}
	}

	if err := hackFile.Close(); err != nil {
		fmt.Printf("Could not close .hack file: %s\n", hackFile.Name())
		panic(err)
	}
}
