package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Build the VMtranslator binary before running the test suite
	build := exec.Command("go", "build", ".")
	err := build.Run()
	if err != nil {
		fmt.Printf("could not build VMtranslator: %v", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestSimpleAdd(t *testing.T) {
	inputPath := "StackArithmetic/SimpleAdd/SimpleAdd.vm"
	testScript := strings.Split(inputPath, ".")[0] + ".tst"

	err := translate(inputPath)
	if err != nil {
		panic(err)
	}

	runCPUEmulator := exec.Command("cmd", "/C", `..\..\tools\CPUEmulator.bat`, testScript)
	output, err := runCPUEmulator.CombinedOutput()
	if err != nil {
		panic(err)
	}

	successMsg := "End of script - Comparison ended successfully"
	if strings.TrimSpace(string(output)) != successMsg {
		t.Errorf("%s", output)
	}

}

func TestStackTest(t *testing.T) {
	inputPath := "StackArithmetic/StackTest/StackTest.vm"
	testScript := strings.Split(inputPath, ".")[0] + ".tst"

	err := translate(inputPath)
	if err != nil {
		panic(err)
	}

	runCPUEmulator := exec.Command("cmd", "/C", `..\..\tools\CPUEmulator.bat`, testScript)
	output, err := runCPUEmulator.CombinedOutput()
	if err != nil {
		panic(err)
	}

	successMsg := "End of script - Comparison ended successfully"
	if strings.TrimSpace(string(output)) != successMsg {
		t.Errorf("%s", output)
	}

}
