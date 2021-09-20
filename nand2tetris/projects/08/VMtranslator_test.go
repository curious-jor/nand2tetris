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

func TestStackArithmetic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
	}{
		// Stack Arithmetic tests from book
		{"SimpleAdd", "StackArithmetic/SimpleAdd/SimpleAdd.vm"},
		{"StackTest", "StackArithmetic/StackTest/StackTest.vm"},

		// Personal tests
		{"SimpleSub", "StackArithmetic/SimpleSub/SimpleSub.vm"},
		{"SimpleNeg", "StackArithmetic/SimpleNeg/SimpleNeg.vm"},
		{"SimpleEq", "StackArithmetic/SimpleEq/SimpleEq.vm"},
		{"SimpleGt", "StackArithmetic/SimpleGt/SimpleGt.vm"},
		{"SimpleLt", "StackArithmetic/SimpleLt/SimpleLt.vm"},
		{"SimpleAnd", "StackArithmetic/SimpleAnd/SimpleAnd.vm"},
		{"SimpleOr", "StackArithmetic/SimpleOr/SimpleOr.vm"},
		{"SimpleNot", "StackArithmetic/SimpleNot/SimpleNot.vm"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testScript := strings.Split(test.input, ".")[0] + ".tst"
			err := translate(test.input)
			if err != nil {
				t.Fatal(err)
			}

			runCPUEmulator := exec.Command("cmd", "/C", `..\..\tools\CPUEmulator.bat`, testScript)
			output, err := runCPUEmulator.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			successMsg := "End of script - Comparison ended successfully"
			if strings.TrimSpace(string(output)) != successMsg {
				t.Errorf("%s", output)
			}
		})
	}
}

// Memory access tests from book
func TestMemoryAccess(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
	}{
		// Memory Access tests from book
		{"BasicTest", "MemoryAccess/BasicTest/BasicTest.vm"},
		{"PointerTest", "MemoryAccess/PointerTest/PointerTest.vm"},
		{"StaticTest", "MemoryAccess/StaticTest/StaticTest.vm"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testScript := strings.Split(test.input, ".")[0] + ".tst"
			err := translate(test.input)
			if err != nil {
				t.Fatal(err)
			}

			runCPUEmulator := exec.Command("cmd", "/C", `..\..\tools\CPUEmulator.bat`, testScript)
			output, err := runCPUEmulator.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			successMsg := "End of script - Comparison ended successfully"
			if strings.TrimSpace(string(output)) != successMsg {
				t.Errorf("%s", output)
			}
		})
	}
}

func TestProgramFlow(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
	}{
		// Program Flow tests from book
		{"BasicLoop", "ProgramFlow/BasicLoop/BasicLoop.vm"},
		{"FibonacciSeries", "ProgramFlow/FibonacciSeries/FibonacciSeries.vm"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testScript := strings.Split(test.input, ".")[0] + ".tst"
			err := translate(test.input)
			if err != nil {
				t.Fatal(err)
			}

			runCPUEmulator := exec.Command("cmd", "/C", `..\..\tools\CPUEmulator.bat`, testScript)
			output, err := runCPUEmulator.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			successMsg := "End of script - Comparison ended successfully"
			if strings.TrimSpace(string(output)) != successMsg {
				t.Errorf("%s", output)
			}
		})
	}
}
