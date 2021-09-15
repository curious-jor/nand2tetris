package codewriter

import (
	"VMtranslator/parser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CodeWriter struct {
	outputFile *os.File
	eqCounter  int    // used to make unique label assembly commands for each vm equality command
	fileName   string // name of the file currently being translated
	label      string
}

type asmBuilder struct {
	b strings.Builder
}

func (ab *asmBuilder) writeString(s string) (int, error) {
	return ab.b.WriteString(s + newLineString)
}

func (ab *asmBuilder) string() string {
	return ab.b.String()
}

func joinASMStrings(strs ...string) string {
	strList := []string{}
	strList = append(strList, strs...)
	return strings.Join(strList, "\n")
}

func NewCodeWriter(outputFile *os.File) *CodeWriter {
	var cw = new(CodeWriter)
	cw.outputFile = outputFile
	cw.eqCounter = 1
	cw.fileName = outputFile.Name()

	_, f := filepath.Split(cw.fileName)
	cw.label = strings.Split(f, ".")[0]
	return cw
}

func (cw *CodeWriter) SetFileName(fileName string) {
	cw.fileName = fileName
}

const unsupportedCmdString = "Unsupported Command: "
const newLineString = "\n"

var cmdsWithAsm = map[string]string{
	"add": "D=D+M",
	"sub": "D=M-D",
	"and": "D=D&M",
	"or":  "D=D|M",
	"not": "D=!M",
	"neg": "D=-M",
}
var segmentsAsm = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
}
var stackPushString = "@SP" + newLineString + "A=M" + newLineString + "M=D"
var stackPopString = "@SP" + newLineString + "AM=M-1" + newLineString + "D=M"
var incrementSPString = "@SP" + newLineString + "M=M+1"
var decrementSPString = "@SP" + newLineString + "AM=M-1"

func (cw *CodeWriter) getBinaryCmdOutput(cmd string) string {
	var output string

	switch cmd {
	case "add", "sub", "and", "or":
		{
			computeInstruction := cmdsWithAsm[cmd]
			loadArg1 := joinASMStrings(
				decrementSPString,
				"D=M",
			)
			loadArg2 := joinASMStrings(
				decrementSPString,
				computeInstruction,
			)

			asm := joinASMStrings(
				loadArg1,
				loadArg2,
				stackPushString,
			)
			output = asm
		}
	case "eq", "gt", "lt":
		{
			var jumpMnemonic string
			if cmd == "eq" {
				jumpMnemonic = "JEQ"
			}
			if cmd == "gt" {
				jumpMnemonic = "JGT"
			}
			if cmd == "lt" {
				jumpMnemonic = "JLT"
			}

			loadArg1 := joinASMStrings(
				decrementSPString,
				"D=M",
			)
			loadArg2 := joinASMStrings(
				decrementSPString,
				"D=M-D",
			)
			checkEquality := joinASMStrings(
				fmt.Sprintf("@EQ%d", cw.eqCounter),
				fmt.Sprintf("D;%s", jumpMnemonic),
				"D=0",
				fmt.Sprintf("@PUSHEQ%d", cw.eqCounter),
				"0;JMP",
				fmt.Sprintf("(EQ%d)", cw.eqCounter),
				"D=-1",
			)
			pushResult := joinASMStrings(
				fmt.Sprintf("(PUSHEQ%d)", cw.eqCounter),
				"@SP",
				"A=M",
				"M=D",
			)

			asm := joinASMStrings(
				loadArg1,
				loadArg2,
				checkEquality,
				pushResult,
			)
			output = asm
		}
	}

	return output
}

func (cw *CodeWriter) getUnaryCmdOutput(cmd string) string {
	var output string

	computeInstruction := cmdsWithAsm[cmd]
	loadArg := joinASMStrings(
		decrementSPString,
		computeInstruction,
	)
	output = joinASMStrings(
		loadArg,
		stackPushString,
		"\n",
	)

	return output
}

func (cw *CodeWriter) WriteArithmetic(command string) error {
	var output asmBuilder
	commandUnsupported := false

	switch command {
	case "add", "sub", "and", "or":
		{
			_, err := output.writeString(cw.getBinaryCmdOutput(command))
			if err != nil {
				return err
			}
		}
	case "eq", "gt", "lt": // all three equality checks use the same logic, but different jump mnemonics
		{
			_, err := output.writeString(cw.getBinaryCmdOutput(command))
			if err != nil {
				return err
			}
			cw.eqCounter += 1
		}
	case "neg", "not":
		{
			_, err := output.writeString(cw.getUnaryCmdOutput(command))
			if err != nil {
				return err
			}
		}
	default:
		{
			output.writeString(unsupportedCmdString + command)
			commandUnsupported = true
		}
	}
	output.writeString(incrementSPString)
	_, err := cw.outputFile.WriteString(output.string())
	if err != nil {
		return err
	}

	if commandUnsupported {
		return fmt.Errorf("attempted to write unsupported arithmetic command: %q", command)
	}

	return nil
}

func (cw *CodeWriter) WritePushPop(command parser.CommandType, segment string, index int) error {
	if !(command == parser.C_PUSH || command == parser.C_POP) {
		return fmt.Errorf("attempted to write %s as push or pop command. expected C_PUSH or C_POP", command.String())
	}

	var output asmBuilder
	if command == parser.C_PUSH {
		switch segment {
		case "constant":
			{
				loadConstant := joinASMStrings(
					fmt.Sprintf("@%d", index),
					"D=A",
					"@SP",
					"A=M",
					"M=D",
				)
				output.writeString(loadConstant)
			}
		case "local", "argument", "this", "that":
			{
				segmentName := segmentsAsm[segment]
				loadIndex := joinASMStrings(
					fmt.Sprintf("@%d", index),
					"D=A",
				)
				loadIndexOfSegment := joinASMStrings(
					fmt.Sprintf("@%s", segmentName),
					"A=D+M",
					"D=M",
				)
				push := stackPushString

				output.writeString(loadIndex)
				output.writeString(loadIndexOfSegment)
				output.writeString(push)
			}
		case "temp":
			{
				loadIndex := joinASMStrings(
					fmt.Sprintf("@%d", index),
					"D=A",
				)
				loadIndexOfSegment := joinASMStrings(
					"@R5",
					"A=D+A",
					"D=M",
				)
				push := stackPushString
				output.writeString(joinASMStrings(loadIndex))
				output.writeString(joinASMStrings(loadIndexOfSegment))
				output.writeString(joinASMStrings(push))
			}
		case "pointer":
			{
				var entry string
				if index == 0 {
					entry = "THIS"
				}
				if index == 1 {
					entry = "THAT"
				}

				loadAddress := joinASMStrings(
					fmt.Sprintf("@%s", entry),
					"D=M",
				)
				push := stackPushString
				output.writeString(loadAddress)
				output.writeString(push)
			}
		case "static":
			{
				symbolCmd := fmt.Sprintf("@%s.%d", cw.label, index)
				loadStatic := joinASMStrings(
					symbolCmd,
					"D=M",
				)
				push := stackPushString
				output.writeString(loadStatic)
				output.writeString(push)
			}
		}

		output.writeString(incrementSPString)
		_, err := cw.outputFile.WriteString(output.string())
		if err != nil {
			return err
		}
	}

	if command == parser.C_POP {
		switch segment {
		case "local", "argument", "this", "that":
			{
				segmentName := segmentsAsm[segment]
				loadIndex := joinASMStrings(
					fmt.Sprintf("@%d", index),
					"D=A",
				)
				loadIndexOfSegment := joinASMStrings(
					fmt.Sprintf("@%s", segmentName),
					"D=D+M",
				)
				storeAddress := joinASMStrings(
					"@R13",
					"M=D",
				)
				popFromStack := stackPopString
				push := joinASMStrings(
					"@R13",
					"A=M",
					"M=D",
				)

				output.writeString(loadIndex)
				output.writeString(loadIndexOfSegment)
				output.writeString(storeAddress)
				output.writeString(popFromStack)
				output.writeString(push)
			}
		case "temp":
			{
				loadIndex := joinASMStrings(
					fmt.Sprintf("@%d", index),
					"D=A",
				)
				loadIndexOfSegment := joinASMStrings(
					"@R5",
					"D=D+A",
				)
				storeAddress := joinASMStrings(
					"@R13",
					"M=D",
				)
				popFromStack := stackPopString
				push := joinASMStrings(
					"@R13",
					"A=M",
					"M=D",
				)

				output.writeString(loadIndex)
				output.writeString(loadIndexOfSegment)
				output.writeString(storeAddress)
				output.writeString(popFromStack)
				output.writeString(push)
			}
		case "pointer":
			{
				var entry string
				if index == 0 {
					entry = "THIS"
				}
				if index == 1 {
					entry = "THAT"
				}

				popFromStack := stackPopString
				push := joinASMStrings(
					fmt.Sprintf("@%s", entry),
					"M=D",
				)
				output.writeString(popFromStack)
				output.writeString(push)
			}
		case "static":
			{
				popFromStack := stackPopString
				symbolCmd := fmt.Sprintf("@%s.%d", cw.label, index)
				push := joinASMStrings(
					symbolCmd,
					"M=D",
				)

				output.writeString(popFromStack)
				output.writeString(push)
			}
		case "constant":
			return fmt.Errorf("attempted to write pop command with %q as segment and %d as index", segment, index)
		}

		_, err := cw.outputFile.WriteString(output.string())
		if err != nil {
			return err
		}
	}

	return nil
}

func (cw *CodeWriter) Close() error {
	return cw.outputFile.Close()
}
