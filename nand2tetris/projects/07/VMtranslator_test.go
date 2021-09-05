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

// Stack Arithmetic tests from book
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

// Personal stack arithmetic tests
func TestSimpleSub(t *testing.T) {
	inputPath := "StackArithmetic/SimpleSub/SimpleSub.vm"
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

func TestSimpleNeg(t *testing.T) {
	inputPath := "StackArithmetic/SimpleNeg/SimpleNeg.vm"
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

func TestSimpleEq(t *testing.T) {
	inputPath := "StackArithmetic/SimpleEq/SimpleEq.vm"
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

func TestSimpleLt(t *testing.T) {
	inputPath := "StackArithmetic/SimpleLt/SimpleLt.vm"
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

func TestSimpleAnd(t *testing.T) {
	inputPath := "StackArithmetic/SimpleAnd/SimpleAnd.vm"
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

func TestSimpleOr(t *testing.T) {
	inputPath := "StackArithmetic/SimpleOr/SimpleOr.vm"
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

func TestSimpleNot(t *testing.T) {
	inputPath := "StackArithmetic/SimpleNot/SimpleNot.vm"
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
