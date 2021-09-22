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
	label      string // used as a prefix in the naming of static variables encountered in the file
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
	_, f := filepath.Split(cw.fileName)
	cw.label = strings.Split(f, ".")[0]
}

const unsupportedCmdString = "Unsupported Command: "
const newLineString = "\n"

// maps VM command names to their assembly equivalent
var cmdsWithAsm = map[string]string{
	"add": "D=D+M",
	"sub": "D=M-D",
	"and": "D=D&M",
	"or":  "D=D|M",
	"not": "D=!M",
	"neg": "D=-M",
}

// maps VM names for memory segments to their assembly equivalent
var segmentsAsm = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
}

// assembly to push contents of D register onto stack
var stackPushString = "@SP" + newLineString + "A=M" + newLineString + "M=D"

// assembly to pop top of stack into D register
var stackPopString = "@SP" + newLineString + "AM=M-1" + newLineString + "D=M"

// assembly to increment stack pointer (SP)
var incrementSPString = "@SP" + newLineString + "M=M+1"

// assembly to decrement stack pointer
var decrementSPString = "@SP" + newLineString + "AM=M-1"

func (cw *CodeWriter) WriteArithmetic(command string) error {
	var output asmBuilder
	outputList := []string{}
	commandUnsupported := false

	switch command {
	case "add", "sub", "and", "or", "eq", "gt", "lt":
		{
			outputList = append(outputList, cw.getBinaryCmdOutput(command))
			if command == "eq" || command == "gt" || command == "lt" {
				cw.eqCounter += 1
			}
		}
	case "neg", "not":
		{
			outputList = append(outputList, cw.getUnaryCmdOutput(command))
		}
	default:
		{
			outputList = append(outputList, unsupportedCmdString+command)
			commandUnsupported = true
		}
	}
	outputList = append(outputList, incrementSPString)

	for _, str := range outputList {
		if _, err := output.writeASMString(str); err != nil {
			return err
		}
	}

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
	outputList := []string{}

	// The assembly to load a constant into data memory is shared by every push/pop command
	loadIndex := joinASMStrings(
		fmt.Sprintf("@%d", index),
		"D=A",
	)

	if command == parser.C_PUSH {
		switch segment {
		case "constant":
			{
				loadConstant := joinASMStrings(
					loadIndex,
					stackPushString,
				)
				outputList = append(outputList, loadConstant)
			}
		case "local", "argument", "this", "that":
			{
				segmentName := segmentsAsm[segment]
				loadIndexOfSegment := joinASMStrings(
					fmt.Sprintf("@%s", segmentName),
					"A=D+M",
					"D=M",
				)
				push := stackPushString
				outputList = append(outputList, loadIndex, loadIndexOfSegment, push)
			}
		case "temp":
			{
				loadIndexOfSegment := joinASMStrings(
					"@R5",
					"A=D+A",
					"D=M",
				)
				push := stackPushString
				outputList = append(outputList, loadIndex, loadIndexOfSegment, push)
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
				outputList = append(outputList, loadAddress, push)
			}
		case "static":
			{
				symbolCmd := fmt.Sprintf("@%s.%d", cw.label, index)
				loadStatic := joinASMStrings(
					symbolCmd,
					"D=M",
				)
				push := stackPushString
				outputList = append(outputList, loadStatic, push)
			}
		}
		outputList = append(outputList, incrementSPString)
	}

	if command == parser.C_POP {
		switch segment {
		case "local", "argument", "this", "that":
			{
				segmentName := segmentsAsm[segment]
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
				outputList = append(outputList, loadIndex, loadIndexOfSegment, storeAddress, popFromStack, push)
			}
		case "temp":
			{
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
				outputList = append(outputList, loadIndex, loadIndexOfSegment, storeAddress, popFromStack, push)
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
				outputList = append(outputList, popFromStack, push)
			}
		case "static":
			{
				popFromStack := stackPopString
				symbolCmd := fmt.Sprintf("@%s.%d", cw.label, index)
				push := joinASMStrings(
					symbolCmd,
					"M=D",
				)
				outputList = append(outputList, popFromStack, push)
			}
		case "constant":
			return fmt.Errorf("attempted to write pop command with %q as segment and %d as index", segment, index)
		}
	}

	for _, str := range outputList {
		if _, err := output.writeASMString(str); err != nil {
			return err
		}
	}

	_, err := cw.outputFile.WriteString(output.string())
	if err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) WriteLabel(label string) error {
	output := joinASMStrings(
		fmt.Sprintf("(%s)", label),
		"",
	)
	_, err := cw.outputFile.WriteString(output)
	if err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteGoto(label string) error {
	output := joinASMStrings(
		fmt.Sprintf("@%s", label),
		"0;JMP",
		"",
	)
	_, err := cw.outputFile.WriteString(output)
	if err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteIf(label string) error {
	var output asmBuilder
	outputList := []string{}

	outputList = append(outputList, stackPopString)

	loadLabel := fmt.Sprintf("@%s", label)
	jump := joinASMStrings(
		loadLabel,
		"D;JNE",
	)

	outputList = append(outputList, jump)

	for _, str := range outputList {
		if _, err := output.writeASMString(str); err != nil {
			return err
		}
	}

	_, err := cw.outputFile.WriteString(output.string())
	if err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteReturn() error {
	output := joinASMStrings(
		"@LCL", // FRAME = LCL
		"D=M",
		"@FRAME",
		"M=D",
		"@5", // RET = *(FRAME-5)
		"D=D-A",
		"A=D",
		"D=M",
		"@RET",
		"M=D",
		stackPopString, // *ARG = pop()
		"@ARG",
		"A=M",
		"M=D",
		"@ARG", // SP = ARG + 1
		"D=M+1",
		"@SP",
		"M=D",
		"@FRAME", // THAT = *(FRAME-1)
		"D=M-1",
		"A=D",
		"D=M",
		"@THAT",
		"M=D",
		"@2", // THIS = *(FRAME-2)
		"D=A",
		"@FRAME",
		"D=M-D",
		"A=D",
		"D=M",
		"@THIS",
		"M=D",
		"@3", // ARG = *(FRAME-3)
		"D=A",
		"@FRAME",
		"D=M-D",
		"A=D",
		"D=M",
		"@ARG",
		"M=D",
		"@4", // LCL = *(FRAME-4)
		"D=A",
		"@FRAME",
		"D=M-D",
		"A=D",
		"D=M",
		"@LCL",
		"M=D",
		"@RET", // goto RET
		"A=M",
		"0;JMP",
		"",
	)

	_, err := cw.outputFile.WriteString(output)
	if err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteFunction(functionName string, numLocals int) error {
	if err := cw.WriteLabel(functionName); err != nil {
		return err
	}

	var output asmBuilder

	// initialize local variables to 0
	for i := 0; i < numLocals; i++ {
		initLCLVar := joinASMStrings(
			"@SP",
			"A=M",
			"M=0",
			incrementSPString,
		)
		if _, err := output.writeASMString(initLCLVar); err != nil {
			return err
		}
	}

	if _, err := cw.outputFile.WriteString(output.string()); err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) Close() error {
	return cw.outputFile.Close()
}

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
	case "eq", "gt", "lt": // all three equality checks use the same logic, but different jump mnemonics
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

			// Create a unique label branch for each instance of equality command encountered in
			// the file(s)
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
	)

	return output
}

// Wraps strings.Builder so we can add newline characters automatically when building string output
type asmBuilder struct {
	b strings.Builder
}

func (ab *asmBuilder) writeASMString(s string) (int, error) {
	return ab.b.WriteString(s + newLineString)
}

func (ab *asmBuilder) string() string {
	return ab.b.String()
}

// Makes it explicit that we're joining together assembly output. Adds newline characters automatically.
func joinASMStrings(strs ...string) string {
	strList := []string{}
	strList = append(strList, strs...)
	return strings.Join(strList, "\n")
}
