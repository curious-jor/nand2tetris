package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var srcPath string
	if len(os.Args) != 2 {
		fmt.Println("VMTranslator expects a .vm file or dir containing .vm files")
		os.Exit(1)
	}
	srcPath = os.Args[1]

	srcFile, err := os.Stat(srcPath)
	if err != nil {
		panic(err)
	}

	switch mode := srcFile.Mode(); {
	case mode.IsDir():
		{
			fmt.Printf("Translating files in %q\n", srcPath)
			dirEntries, err := os.ReadDir(srcPath)
			if err != nil {
				panic(err)
			}
			for _, entry := range dirEntries {
				if filepath.Ext(entry.Name()) == ".vm" {
					fmt.Println("Found", entry.Name())
				}

			}
		}
	case mode.IsRegular():
		{
			fmt.Printf("Translating %q\n", srcPath)
		}
	}

}
