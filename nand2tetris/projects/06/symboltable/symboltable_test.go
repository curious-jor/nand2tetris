package symboltable

import "testing"

func TestSymbolTableInitialization(t *testing.T) {
	predefined := map[string]int{
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

	st := NewSymbolTable()

	for symbol, expected := range predefined {
		val, ok := st.t[symbol]
		if !ok {
			t.Errorf("could not retrieve symbol %s from new SymbolTable", symbol)
		}
		if val != expected {
			t.Errorf("expected value %d for key %s but got %d", expected, symbol, val)
		}
	}
}

func TestAddEntry(t *testing.T) {
	st := NewSymbolTable()

	// test adding a label
	label, val := "_label.1$:", 0
	if err := st.AddEntry(label, val); err != nil {
		t.Errorf("could not add label, val (%q, %d) to SymbolTable %v", label, val, err)
	}

	// test adding constant returns error
	constant := "16"
	if err := st.AddEntry(constant, 16); err == nil {
		t.Errorf("adding constant %s to SymbolTable should return error", constant)
	}
}

func TestContains(t *testing.T) {
	st := NewSymbolTable()

	notIn := "NOTIN"
	if ok := st.Contains(notIn); ok {
		t.Errorf("new SymbolTable should not contain any user-defined labels")
	}

	loop, val := "LOOP", 5
	if err := st.AddEntry(loop, val); err != nil {
		t.Fatalf("error while adding (%q, %d) to SymbolTable: %v", loop, val, err)
	}
	if expected, ok := st.t[loop]; !ok {
		t.Errorf("key %s was not in SymbolTable", loop)
	} else if expected != val {
		t.Errorf("expected %d but got %d", expected, val)
	}

}
