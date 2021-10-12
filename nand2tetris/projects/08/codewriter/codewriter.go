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
	retCounter int
	currFnName string
}

func NewCodeWriter(outputFile *os.File) *CodeWriter {
	var cw = new(CodeWriter)
	cw.outputFile = outputFile
	cw.eqCounter = 1
	cw.retCounter = 1
	cw.fileName = outputFile.Name()

	_, f := filepath.Split(cw.fileName)
	cw.label = strings.Split(f, ".")[0]
	return cw
}

func (cw *CodeWriter) SetFileName(fileName string) {
	cw.fileName = fileName
	cw.currFnName = ""
	cw.eqCounter = 1
	cw.retCounter = 1
	_, f := filepath.Split(cw.fileName)
	cw.label = strings.Split(f, ".")[0]
}

const unsupportedCmdString = "Unsupported Command: "

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

var cmdMnemonics = map[parser.CommandType]string{
	parser.C_PUSH: "push",
	parser.C_POP:  "pop",
}

// assembly to push contents of D register onto stack
var stackPushString = "@SP" + "\n\t" + "A=M" + "\n\t" + "M=D"

// assembly to pop top of stack into D register
var stackPopString = "@SP" + "\n\t" + "AM=M-1" + "\n\t" + "D=M"

// assembly to increment stack pointer (SP)
var incrementSPString = "@SP" + "\n\t" + "M=M+1"

// assembly to decrement stack pointer
var decrementSPString = "@SP" + "\n\t" + "AM=M-1"

func (cw *CodeWriter) WriteArithmetic(command string) error {
	var output strings.Builder
	outputList := []string{}
	commandUnsupported := false

	switch command {
	case "add", "sub", "and", "or", "eq", "gt", "lt":
		{
			asm := cw.getBinaryCmdOutput(command)
			if command == "eq" || command == "gt" || command == "lt" {
				outputList = append(outputList, fmt.Sprintf("// %s %d", command, cw.eqCounter))
				cw.eqCounter += 1
			} else {
				outputList = append(outputList, fmt.Sprintf("// %s", command))
			}
			outputList = append(outputList, asm)
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
		if _, err := output.WriteString("\t" + str + "\n"); err != nil {
			return err
		}
	}

	_, err := cw.outputFile.WriteString(output.String())
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

	var output strings.Builder
	outputList := []string{}

	outputList = append(outputList, fmt.Sprintf("// %s %s %d", cmdMnemonics[command], segment, index))

	// The assembly to load a constant into data memory is shared by every push/pop command
	loadIndex := strings.Join([]string{
		fmt.Sprintf("@%d", index),
		"D=A",
	}, "\n\t")

	if command == parser.C_PUSH {
		switch segment {
		case "constant":
			{
				loadConstant := strings.Join([]string{
					loadIndex,
					stackPushString,
				}, "\n\t")
				outputList = append(outputList, loadConstant)
			}
		case "local", "argument", "this", "that":
			{
				segmentName := segmentsAsm[segment]
				loadIndexOfSegment := strings.Join([]string{
					fmt.Sprintf("@%s", segmentName),
					"A=D+M",
					"D=M",
				}, "\n\t")
				push := stackPushString
				outputList = append(outputList, loadIndex, loadIndexOfSegment, push)
			}
		case "temp":
			{
				loadIndexOfSegment := strings.Join([]string{
					"@R5",
					"A=D+A",
					"D=M",
				}, "\n\t")
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
				loadAddress := strings.Join([]string{
					fmt.Sprintf("@%s", entry),
					"D=M",
				}, "\n\t")
				push := stackPushString
				outputList = append(outputList, loadAddress, push)
			}
		case "static":
			{
				symbolCmd := fmt.Sprintf("@%s.%d", cw.label, index)
				loadStatic := strings.Join([]string{
					symbolCmd,
					"D=M",
				}, "\n\t")
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
				loadIndexOfSegment := strings.Join([]string{
					fmt.Sprintf("@%s", segmentName),
					"D=D+M",
				}, "\n\t")
				storeAddress := strings.Join([]string{
					"@R13",
					"M=D",
				}, "\n\t")
				popFromStack := stackPopString
				push := strings.Join([]string{
					"@R13",
					"A=M",
					"M=D",
				}, "\n\t")
				outputList = append(outputList, loadIndex, loadIndexOfSegment, storeAddress, popFromStack, push)
			}
		case "temp":
			{
				loadIndexOfSegment := strings.Join([]string{
					"@R5",
					"D=D+A",
				}, "\n\t")
				storeAddress := strings.Join([]string{
					"@R13",
					"M=D",
				}, "\n\t")
				popFromStack := stackPopString
				push := strings.Join([]string{
					"@R13",
					"A=M",
					"M=D",
				}, "\n\t")
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
				push := strings.Join([]string{
					fmt.Sprintf("@%s", entry),
					"M=D",
				}, "\n\t")
				outputList = append(outputList, popFromStack, push)
			}
		case "static":
			{
				popFromStack := stackPopString
				symbolCmd := fmt.Sprintf("@%s.%d", cw.label, index)
				push := strings.Join([]string{
					symbolCmd,
					"M=D",
				}, "\n\t")
				outputList = append(outputList, popFromStack, push)
			}
		case "constant":
			return fmt.Errorf("attempted to write pop command with %q as segment and %d as index", segment, index)
		}
	}

	for _, str := range outputList {
		if _, err := output.WriteString("\t" + str + "\n"); err != nil {
			return err
		}
	}

	_, err := cw.outputFile.WriteString(output.String())
	if err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) WriteLabel(label string) error {
	var output string
	var labelGen string
	if cw.currFnName != "" {
		labelGen = fmt.Sprintf("(%s$%s)", cw.currFnName, label)
	} else {
		labelGen = fmt.Sprintf("(%s)", label)
	}

	output = strings.Join([]string{
		fmt.Sprintf("// label %s", labelGen),
		labelGen,
	}, "\n\t")

	if _, err := cw.outputFile.WriteString(output); err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteGoto(label string) error {
	var loadLabel string
	if cw.currFnName != "" {
		loadLabel = fmt.Sprintf("@%s$%s", cw.currFnName, label)
	} else {
		loadLabel = fmt.Sprintf("@%s", label)
	}
	output := strings.Join([]string{
		fmt.Sprintf("// goto %s", label),
		loadLabel,
		"0;JMP",
		""}, "\n\t")
	_, err := cw.outputFile.WriteString(output)
	if err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteIf(label string) error {
	var loadLabel string
	if cw.currFnName != "" {
		loadLabel = fmt.Sprintf("@%s$%s", cw.currFnName, label)
	} else {
		loadLabel = fmt.Sprintf("@%s", label)
	}
	output := strings.Join([]string{
		fmt.Sprintf("// if-goto %s", label),
		stackPopString,
		// load label
		loadLabel,
		"D;JNE",
	}, "\n\t")

	_, err := cw.outputFile.WriteString("\t" + output + "\n")
	if err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteInit() error {
	initSP := strings.Join([]string{
		"// Initialize the stack pointer to 0x0100",
		"@256",
		"D=A",
		"@SP",
		"M=D",
		"// Start executing the translated code of Sys.init",
	}, "\n\t")
	if _, err := cw.outputFile.WriteString("\t" + initSP + "\n"); err != nil {
		return err
	}
	cw.currFnName = "Sys.init"
	return cw.WriteCall("Sys.init", 0)
}

func (cw *CodeWriter) WriteCall(functionName string, numArgs int) error {
	var retAddrLabel string
	if cw.currFnName != "" {
		retAddrLabel = fmt.Sprintf("%s$ret.%d", cw.currFnName, cw.retCounter)
	} else {
		retAddrLabel = fmt.Sprintf("ret.%d", cw.retCounter)
	}
	output := strings.Join([]string{
		fmt.Sprintf("// call %s %d", functionName, numArgs),
		// push return-address
		fmt.Sprintf("// push RET%d", cw.retCounter),
		fmt.Sprintf("@%s", retAddrLabel),
		"D=A",
		stackPushString,
		incrementSPString,
		// push LCL
		"// push LCL",
		"@LCL",
		"D=M",
		stackPushString,
		incrementSPString,
		// push ARG
		"// push ARG",
		"@ARG",
		"D=M",
		stackPushString,
		incrementSPString,
		// push THIS
		"// push THIS",
		"@THIS",
		"D=M",
		stackPushString,
		incrementSPString,
		// push THAT
		"// push THAT",
		"@THAT",
		"D=M",
		stackPushString,
		incrementSPString,
		// ARG = SP-n-5
		"// ARG = SP-n-5",
		"@SP",
		"D=M",
		fmt.Sprintf("@%d", numArgs),
		"D=D-A",
		"@5",
		"D=D-A",
		"@ARG",
		"M=D",
		// LCL = SP
		"// LCL = SP",
		"@SP",
		"D=M",
		"@LCL",
		"M=D",
		// goto f
		fmt.Sprintf("// goto %s", functionName),
		fmt.Sprintf("@%s", functionName),
		"0;JMP",
		fmt.Sprintf("(%s)", retAddrLabel),
	}, "\n\t")
	cw.retCounter += 1

	_, err := cw.outputFile.WriteString("\t" + output + "\n")
	if err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteReturn() error {
	output := strings.Join([]string{
		"// return",
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
	}, "\n\t")

	_, err := cw.outputFile.WriteString("\t" + output + "\n")
	if err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteFunction(functionName string, numLocals int) error {
	if _, err := cw.outputFile.WriteString(fmt.Sprintf("// function %s %d\n", functionName, numLocals)); err != nil {
		return err
	}
	if _, err := cw.outputFile.WriteString(fmt.Sprintf("(%s)\n", functionName)); err != nil {
		return err
	}
	cw.currFnName = functionName
	cw.retCounter = 1
	cw.eqCounter = 1

	var output strings.Builder

	// initialize local variables to 0
	for i := 0; i < numLocals; i++ {
		initLCLVar := strings.Join([]string{
			"@SP",
			"A=M",
			"M=0",
			incrementSPString,
		}, "\n\t")
		if _, err := output.WriteString("\t" + initLCLVar + "\n"); err != nil {
			return err
		}
	}

	if _, err := cw.outputFile.WriteString(output.String()); err != nil {
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
			output = strings.Join([]string{
				// loadArg1
				decrementSPString,
				"D=M",
				// loadArg2
				decrementSPString,
				computeInstruction,
				stackPushString,
			}, "\n\t")
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

			output = strings.Join([]string{
				// loadArg1
				decrementSPString,
				"D=M",
				// loadArg2
				decrementSPString,
				"D=M-D",
				// Create a unique label branch for each instance of equality command encountered in
				// the file(s)
				fmt.Sprintf("@EQ%d", cw.eqCounter),
				fmt.Sprintf("D;%s", jumpMnemonic),
				"D=0",
				fmt.Sprintf("@PUSHEQ%d", cw.eqCounter),
				"0;JMP",
				fmt.Sprintf("(EQ%d)", cw.eqCounter),
				"D=-1",
				// pushResult
				fmt.Sprintf("(PUSHEQ%d)", cw.eqCounter),
				"@SP",
				"A=M",
				"M=D",
			}, "\n\t")
		}
	}

	return output
}

func (cw *CodeWriter) getUnaryCmdOutput(cmd string) string {
	computeInstruction := cmdsWithAsm[cmd]
	output := strings.Join([]string{
		decrementSPString,
		computeInstruction,
		stackPushString,
	}, "\n\t")

	return output
}
