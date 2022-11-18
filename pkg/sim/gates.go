package sim

import "fmt"

// Wire element of wire table
type Wire struct {
	Name    string
	State   int
	Inputs  []int // connected gate list
	Outputs []int // connected gate list
}

// Gate element of gate table
type Gate struct {
	Name    string
	State   int
	Func    func([]int) int
	Inputs  []int // ids of input wires
	Outputs []int // ids of output wires
}

func NewWire(name string) (id int) {
	id = len(wireTable)
	if len(name) == 0 {
		name = fmt.Sprintf("_t%d", id)
	}
	wireTable = append(wireTable, Wire{
		Name:    name,
		State:   StUndefined,
		Outputs: []int{}})
	return id
}

// TrFunc gate transmission function
type TrFunc func([]int) int

// AND transmission function
func AndFunc(in []int) (out int) {
	out = wireTable[in[0]].State
	for _, i := range in[1:] {
		out &= wireTable[i].State
	}
	return out
}

// OR transmission function
func OrFunc(in []int) (out int) {
	out = wireTable[in[0]].State
	for _, i := range in[1:] {
		out |= wireTable[i].State
	}
	return out
}

// Buffer adds delay to signal
func BufFunc(in []int) (out int) {
	out = wireTable[in[0]].State
	return out
}

// NOT transmission function
func NotFunc(in []int) (out int) {
	out = wireTable[in[0]].State
	if out == StFalse {
		return StTrue
	} else if out == StTrue {
		return StFalse
	}
	return out
}

// Attach gate input/outputs to wires
func attach(g int, inputs []int, outputs []int) {
	// attach inputs/outputs wires to gate
	gateTable[g].Inputs = inputs
	gateTable[g].Outputs = outputs
	// attach gate to inputs/output wires
	for _, w := range inputs {
		found := false
		for j := range wireTable[w].Outputs {
			if wireTable[w].Outputs[j] == g {
				found = true
			}
		}
		if !found {
			wireTable[w].Outputs = append(wireTable[w].Outputs, g)
		}
	}
	for _, w := range outputs {
		found := false
		for j := range wireTable[w].Inputs {
			if wireTable[w].Inputs[j] == g {
				found = true
			}
		}
		if !found {
			wireTable[w].Inputs = append(wireTable[w].Inputs, g)
		}
	}
}

func NewGate(name string, trfunc TrFunc, inputs []int, outputs []int) (id int) {
	id = len(gateTable)
	if len(name) == 0 {
		name = "_g"
	}
	gateTable = append(gateTable, Gate{
		Name:  fmt.Sprintf("%s%d", name, id),
		State: StUndefined,
		Func:  trfunc})
	attach(id, inputs, outputs)
	return id
}
