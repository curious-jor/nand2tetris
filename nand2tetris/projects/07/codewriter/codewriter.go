package codewriter

import (
	"VMtranslator/parser"
	"fmt"
	"os"
	"strings"
)

type CodeWriter struct {
	outputFile *os.File
	eqCounter  int // used to make unique label assembly commands for each vm equality command
}

func NewCodeWriter(outputFile *os.File) (*CodeWriter, error) {
	var cw = new(CodeWriter)
	cw.outputFile = outputFile
	cw.eqCounter = 1
	return cw, nil
}

func (cw *CodeWriter) SetFileName(fileName string) error {
	if err := os.Rename(cw.outputFile.Name(), fileName); err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) WriteArithmetic(command string) error {
	var output strings.Builder

	switch command {
	case "add":
		{
			loadArg1 := []string{
				"@SP",
				"AM=M-1",
				"D=M\n",
			}
			loadArg2 := []string{
				"@SP",
				"AM=M-1",
				"D=D+M\n",
			}
			pushResult := []string{
				"@SP",
				"A=M",
				"M=D\n",
			}
			output.WriteString(strings.Join(loadArg1, "\n"))
			output.WriteString(strings.Join(loadArg2, "\n"))
			output.WriteString(strings.Join(pushResult, "\n"))
		}
	case "sub":
		{

			loadArg1 := []string{
				"@SP",
				"AM=M-1",
				"D=M\n",
			}
			loadArg2 := []string{
				"@SP",
				"AM=M-1",
				"D=M-D\n",
			}
			pushResult := []string{
				"@SP",
				"A=M",
				"M=D\n",
			}

			output.WriteString(strings.Join(loadArg1, "\n"))
			output.WriteString(strings.Join(loadArg2, "\n"))
			output.WriteString(strings.Join(pushResult, "\n"))

		}
	case "neg":
		{
			loadArg1 := []string{
				"@SP",
				"AM=M-1",
				"D=-M\n",
			}
			pushResult := []string{
				"@SP",
				"A=M",
				"M=D\n",
			}

			output.WriteString(strings.Join(loadArg1, "\n"))
			output.WriteString(strings.Join(pushResult, "\n"))
		}
	case "eq", "gt", "lt": // all three equality checks use the same logic, but different jump mnemonics
		{
			var jumpMnemonic string
			switch command {
			case "eq":
				jumpMnemonic = "JEQ"
			case "gt":
				jumpMnemonic = "JGT"
			case "lt":
				jumpMnemonic = "JLT"
			}

			loadArg1 := []string{
				"@SP",
				"AM=M-1",
				"D=M\n",
			}
			loadArg2 := []string{
				"@SP",
				"AM=M-1",
				"D=M-D\n",
			}
			checkEquality := []string{
				fmt.Sprintf("@EQ%d", cw.eqCounter),
				fmt.Sprintf("D;%s", jumpMnemonic),
				"D=0",
				fmt.Sprintf("@PUSHEQ%d", cw.eqCounter),
				"0;JMP",
				fmt.Sprintf("(EQ%d)", cw.eqCounter),
				"D=-1\n",
			}
			pushResult := []string{
				fmt.Sprintf("(PUSHEQ%d)", cw.eqCounter),
				"@SP",
				"A=M",
				"M=D\n",
			}

			output.WriteString(strings.Join(loadArg1, "\n"))
			output.WriteString(strings.Join(loadArg2, "\n"))
			output.WriteString(strings.Join(checkEquality, "\n"))
			output.WriteString(strings.Join(pushResult, "\n"))
			cw.eqCounter += 1
		}
	case "and":
		{
			loadArg1 := []string{
				"@SP",
				"AM=M-1",
				"D=M\n",
			}
			loadArg2 := []string{
				"@SP",
				"AM=M-1",
				"D=D&M\n",
			}
			pushResult := []string{
				"@SP",
				"A=M",
				"M=D\n",
			}
			output.WriteString(strings.Join(loadArg1, "\n"))
			output.WriteString(strings.Join(loadArg2, "\n"))
			output.WriteString(strings.Join(pushResult, "\n"))
		}
	case "or":
		{
			loadArg1 := []string{
				"@SP",
				"AM=M-1",
				"D=M\n",
			}
			loadArg2 := []string{
				"@SP",
				"AM=M-1",
				"D=D|M\n",
			}
			pushResult := []string{
				"@SP",
				"A=M",
				"M=D\n",
			}
			output.WriteString(strings.Join(loadArg1, "\n"))
			output.WriteString(strings.Join(loadArg2, "\n"))
			output.WriteString(strings.Join(pushResult, "\n"))
		}
	case "not":
		{
			loadArg1 := []string{
				"@SP",
				"AM=M-1",
				"D=!M\n",
			}
			pushResult := []string{
				"@SP",
				"A=M",
				"M=D\n",
			}
			output.WriteString(strings.Join(loadArg1, "\n"))
			output.WriteString(strings.Join(pushResult, "\n"))
		}
	}

	n, err := cw.outputFile.WriteString(output.String())
	if err != nil {
		return err
	}
	if n < len(output.String()) {
		return fmt.Errorf("underwrote string with call to WriteArithmetic with arg: %q", command)
	}

	if err := cw.writeIncrementSP(); err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) writeIncrementSP() error {
	incrementSPString := "@SP\nM=M+1\n"
	l := len(incrementSPString)
	n, err := cw.outputFile.WriteString(incrementSPString)
	if err != nil {
		return err
	}
	if n < l {
		return fmt.Errorf("wrote %d chars but expected %d chars while writing increment SP output string", n, l)
	}

	return nil
}

func (cw *CodeWriter) WritePushPop(command parser.CommandType, segment string, index int) error {
	if !(command == parser.C_PUSH || command == parser.C_POP) {
		return fmt.Errorf("attempted to write %s as push or pop command. expected C_PUSH or C_POP", command.String())
	}

	if command == parser.C_PUSH {
		if segment == "constant" {
			var output strings.Builder
			loadConstant := []string{
				fmt.Sprintf("@%d", index),
				"D=A",
				"@SP",
				"A=M",
				"M=D\n",
			}
			output.WriteString(strings.Join(loadConstant, "\n"))
			n, err := cw.outputFile.WriteString(output.String())
			if err != nil {
				return err
			}
			if n < len(output.String()) {
				return fmt.Errorf("underwrote string from call to WritePushPop with args: %s, %q, %d", command.String(), segment, index)
			}

			if err := cw.writeIncrementSP(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cw *CodeWriter) Close() error {
	return cw.outputFile.Close()
}
