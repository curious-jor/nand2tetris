package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// linesEqual compares two hack binary files line-by-line after stripping whitespace.
// Returns true and 0 if they're equivalent. Else returns false and the line number where the comparison failed.
// Needed because the built-in assembler ends lines in \r\n but my implementation ends them in \n.
func linesEqual(a *os.File, b *os.File, t *testing.T) (bool, int) {
	r1, r2 := bufio.NewReader(a), bufio.NewReader(b)
	lineCounter := 0
	for {
		line1, err1 := r1.ReadString('\n')
		line2, err2 := r2.ReadString('\n')
		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true, 0
			} else if err1 == io.EOF || err2 == io.EOF {
				t.Log(line1, line2)
				return false, lineCounter
			} else {
				t.Fatal(err1, err2)
			}
		}

		line1 = strings.TrimSpace(line1)
		line2 = strings.TrimSpace(line2)
		if line1 != line2 {
			t.Log(line1, line2)
			return false, lineCounter
		}
		lineCounter += 1
	}
}

func TestMain(m *testing.M) {
	// build the assembler binary before running the tests.
	build := exec.Command("go", "build", ".")
	err := build.Run()
	if err != nil {
		fmt.Printf("could not build assembler %v", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func testHelper(path string, t *testing.T) (*os.File, func()) {
	dir := filepath.Dir(path)
	tempFile, err := os.CreateTemp(dir, "*.asm")
	if err != nil {
		t.Fatalf("could not create temp file %s: %v", path, err)
	}

	tearDown := func() {
		if err = tempFile.Close(); err != nil {
			t.Fatalf("could not close temp file %s: %v", tempFile.Name(), err)
		}
		if err = os.Remove(tempFile.Name()); err != nil {
			t.Fatalf("could not remove temp file %s: %v", tempFile.Name(), err)
		}
	}

	return tempFile, tearDown
}

func TestNoSymbols(t *testing.T) {
	tests := []struct {
		name            string
		fpath           string
		compareFilepath string
	}{
		{"Add", "add/Add.asm", "add/AddCompare.hack"},
		{"MaxL", "max/MaxL.asm", "max/MaxLCompare.hack"},
		{"PongL", "pong/PongL.asm", "pong/PongLCompare.hack"},
		{"RectL", "rect/RectL.asm", "rect/RectLCompare.hack"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			src, err := os.Open(test.fpath)
			if err != nil {
				t.Fatalf("could not open test file: %v", err)
			}
			defer src.Close()

			tempFile, tearDown := testHelper(test.fpath, t)
			defer tearDown()

			// Copy test file contents to the temp file for this test.
			io.Copy(tempFile, src)

			// Run assembler by calling main
			os.Args = []string{"./assembler", tempFile.Name()}
			out := bytes.NewBuffer(nil)
			main()
			t.Log(out)
			outputFilename := strings.Split(tempFile.Name(), ".")[0] + ".hack"

			// open temp file and compare file for reading line by line
			outputFile, err := os.Open(outputFilename)
			if err != nil {
				t.Fatalf("could not open %s %v", outputFile.Name(), err)
			}

			compareFile, err := os.Open(test.compareFilepath)
			if err != nil {
				t.Fatalf("could not open %s %v", test.compareFilepath, err)
			}
			defer compareFile.Close()

			if equal, n := linesEqual(outputFile, compareFile, t); !equal {
				t.Errorf("comparison failed at line %d", n)
			}

			if err := outputFile.Close(); err != nil {
				panic(err)
			}

			if err := os.Remove(outputFile.Name()); err != nil {
				panic(err)
			}
		})
	}
}

func TestSymbols(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		compareFile string
	}{
		{"Max", "max/Max.asm", "max/MaxCompare.hack"},
		{"Rect", "rect/Rect.asm", "rect/RectCompare.hack"},
		{"Pong", "pong/Pong.asm", "pong/PongCompare.hack"},
	}

	for _, test := range tests {
		src, err := os.Open(test.input)
		if err != nil {
			t.Fatalf("could not open test file: %v", err)
		}
		defer src.Close()

		tempFile, tearDown := testHelper(test.input, t)
		defer tearDown()

		// Copy test file contents to the temp file for this test.
		io.Copy(tempFile, src)

		// Run assembler by calling main
		os.Args = []string{"./assembler", tempFile.Name()}
		out := bytes.NewBuffer(nil)
		main()
		t.Log(out)

		outputFilename := strings.Split(tempFile.Name(), ".")[0] + ".hack"

		// open temp file and compare file for reading line by line
		outputFile, err := os.Open(outputFilename)
		if err != nil {
			t.Fatalf("could not open %s %v", outputFile.Name(), err)
		}

		compareFile, err := os.Open(test.compareFile)
		if err != nil {
			t.Fatalf("could not open %s %v", test.compareFile, err)
		}
		defer compareFile.Close()

		if equal, n := linesEqual(outputFile, compareFile, t); !equal {
			t.Errorf("comparison failed at line %d", n)
		}

		if err := outputFile.Close(); err != nil {
			panic(err)
		}

		if err := os.Remove(outputFile.Name()); err != nil {
			panic(err)
		}
	}
}
