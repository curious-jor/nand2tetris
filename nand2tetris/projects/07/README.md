# Virtual Machine I: Stack Arithmetic
## About
This is my Go code for project 7. It's split into 3 modules and the driver program: `codewriter`, `lexer`, `parser`, and `VMtranslator.go`. I decided to add a separate lexer module just to get experience writing one.

## Requirements
The code requires that you have Go version 1.16 installed. It also assumes that you're on Windows.
The tests in `VMtranslator_test.go` will fail if you're not on Windows because they invoke the built-in `CPUEmulator.bat` script.

## Setup and Run
1. Build the program by running 
``` 
go build .
```
2. Run the program using
```
.\VMtranslator source
```
where source is the name of a Hack VM program. 

Ex: `StackArithmetic\SimpleAdd\SimpleAdd.vm`
