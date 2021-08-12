// "Keeps a correspondence between symbolic labels and numeric addresses."
package symboltable

import (
	"fmt"
	"strconv"
)

type SymbolTable struct {
	t map[string]int
}

var predefined = map[string]int{
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"R0":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SCREEN": 16384,
	"KBD":    24576,
}

func NewSymbolTable() *SymbolTable {
	st := new(SymbolTable)
	st.t = predefined
	return st
}

// "Adds the pair (symbol, address) to the table."
func (st *SymbolTable) AddEntry(symbol string, address int) error {
	const sixteenBit = 16
	const baseTen = 10
	_, err := strconv.ParseInt(symbol, baseTen, sixteenBit)
	if err == nil {
		return fmt.Errorf("attempted to add constant %s to SymbolTble", symbol)
	}

	st.t[symbol] = address

	return nil
}

// "Does the symbol table contain the given symbol?"
func (st *SymbolTable) Contains(symbol string) bool {
	_, ok := st.t[symbol]
	return ok
}

// "Returns the address associated with the symbol."
func (st *SymbolTable) GetAddress(symbol string) int {
	address, ok := st.t[symbol]
	if !ok {
		return -1
	}
	return address
}
