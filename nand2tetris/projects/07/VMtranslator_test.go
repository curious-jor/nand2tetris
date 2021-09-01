package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type FileComparison struct {
	equal          bool
	differingLines []string
	lineNumber     int
}

func filesEqual(f1 *os.File, f2 *os.File) *FileComparison {
	var comparison *FileComparison = new(FileComparison)
	s1, s2 := bufio.NewScanner(f1), bufio.NewScanner(f2)
	lineCounter := 1
	for {
		more1, more2 := s1.Scan(), s2.Scan()
		if !more1 && !more2 {
			err1, err2 := s1.Err(), s2.Err()
			if err1 == nil && err2 == nil { // both EOF
				comparison.equal = true
				comparison.differingLines = []string{}
				comparison.lineNumber = -1
				return comparison
			} else {
				comparison.equal = false
				comparison.differingLines = []string{s1.Text(), s2.Text()}
				comparison.lineNumber = lineCounter
				return comparison
			}
		}

		line1, line2 := s1.Text(), s2.Text()
		if line1 != line2 {
			comparison.equal = false
			comparison.differingLines = []string{s1.Text(), s2.Text()}
			comparison.lineNumber = lineCounter
			return comparison
		}

		lineCounter += 1
	}
}

func TestSimpleAdd(t *testing.T) {
	inputFile, err := os.Open("StackArithmetic/SimpleAdd/SimpleAdd.vm")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	tempFile, err := os.CreateTemp(filepath.Dir(inputFile.Name()), "*.vm")
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(tempFile, inputFile); err != nil {
		panic(err)
	}

	if err := translate(tempFile.Name()); err != nil {
		panic(err)
	}
	outputFilename := strings.Split(tempFile.Name(), ".")[0] + ".asm"

	outputFile, err := os.Open(outputFilename)
	if err != nil {
		panic(err)
	}

	compareFile, err := os.Open("StackArithmetic/SimpleAdd/SimpleAddManual.asm")
	if err != nil {
		panic(err)
	}
	defer compareFile.Close()

	if compare := filesEqual(compareFile, outputFile); !compare.equal {
		t.Errorf("comparison failed at line %d. %q != %q", compare.lineNumber, compare.differingLines[0], compare.differingLines[1])
	}

	if err := outputFile.Close(); err != nil {
		panic(err)
	}

	if err := os.Remove(outputFile.Name()); err != nil {
		panic(err)
	}

	if err := tempFile.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(tempFile.Name()); err != nil {
		panic(err)
	}
}
