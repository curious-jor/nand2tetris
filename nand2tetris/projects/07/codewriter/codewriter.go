package codewriter

import (
	"VMtranslator/parser"
	"fmt"
	"os"
	"strings"
)

type CodeWriter struct {
	outputFile *os.File
}

func NewCodeWriter(outputFile *os.File) (*CodeWriter, error) {
	var cw = new(CodeWriter)
	file, err := os.Open(outputFile.Name())
	if err != nil {
		return nil, err
	}

	cw.outputFile = file
	return cw, nil
}

func (cw *CodeWriter) SetFileName(fileName string) error {
	if err := os.Rename(cw.outputFile.Name(), fileName); err != nil {
		return err
	}

	return nil
}

// TODO: Implement cases for all nine arithmetic commands
func (cw *CodeWriter) WriteArithmetic(command string) error {
	if command == "add" {
		var output strings.Builder
		loadArg1 := []string{
			"@SP",
			"M=M-1",
			"A=M",
			"D=M\n",
		}
		loadArg2 := []string{
			"@SP",
			"M=M-1",
			"A=M",
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
			loadConstant := fmt.Sprintf(
				`@%d
				D=A
				@SP
				A=M
				M=D
				`, index)
			output.WriteString(loadConstant)
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
