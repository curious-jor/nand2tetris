package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestSimpleAdd(t *testing.T) {
	inputPath := "StackArithmetic/SimpleAdd/SimpleAdd.vm"
	testScript := strings.Split(inputPath, ".")[0] + ".tst"

	err := translate(inputPath)
	if err != nil {
		panic(err)
	}

	runCPUEmulator := exec.Command("cmd", "/C", `..\..\tools\CPUEmulator.bat`, testScript)
	output, err := runCPUEmulator.Output()
	if err != nil {
		panic(err)
	}

	successMsg := "End of script - Comparison ended successfully"
	if strings.TrimSpace(string(output)) != successMsg {
		t.Errorf("%s", output)
	}

}
