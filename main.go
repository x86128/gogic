package main

import (
	"fmt"
)

// Signal states
const (
	stFalse     = 0 // 0
	stInvalid   = 1
	stUndefined = 2 // X
	stTrue      = 3 // 1
)

// Wire element of wire table
type Wire struct {
	Name    string
	State   int
	Outputs map[int]bool // set of connected gate ids
}

// Gate element of gate table
type Gate struct {
	Name    string
	State   int
	Func    func(map[int]bool) int
	Inputs  map[int]bool // ids of input wires
	Outputs map[int]bool // ids of output wires
}

var wireTable map[int]*Wire

func newWire(name string) (id int) {
	id = len(wireTable)
	wireTable[id] = &Wire{
		Name:    name,
		State:   stUndefined,
		Outputs: make(map[int]bool)}
	return id
}

// TrFunc gate transmission function
type TrFunc func(map[int]bool) int

func AndFunc(in map[int]bool) (out int) {
	out = -1
	for i := range in {
		if out < 0 {
			out = wireTable[i].State
		} else {
			out &= wireTable[i].State
		}
	}
	return out
}

func NotFunc(in map[int]bool) (out int) {
	for i := range in {
		out = wireTable[i].State
		if out == stFalse {
			return stTrue
		} else if out == stTrue {
			return stFalse
		}
	}
	return out
}

var gateTable map[int]*Gate

func newGate(name string, trfunc TrFunc) (id int) {
	id = len(gateTable)
	gateTable[id] = &Gate{
		Name:    name,
		State:   stUndefined,
		Func:    trfunc,
		Inputs:  make(map[int]bool),
		Outputs: make(map[int]bool)}
	return id
}

// Attach gate input/outputs to wires
func attach(gate int, inputs []int, outputs []int) {
	g := gateTable[gate]
	for i := range inputs {
		g.Inputs[inputs[i]] = true
		wireTable[inputs[i]].Outputs[gate] = true
	}
	for i := range outputs {
		g.Outputs[outputs[i]] = true
	}
}

// func BufferFunc(in []bool) (out []bool) {
// 	for _, v := range in {
// 		out = append(out, v)
// 	}
// 	return out
// }

// func OrFunc(in []bool) (out []bool) {
// 	res := in[0]
// 	for _, v := range in[1:] {
// 		res = res || v
// 	}
// 	return []bool{res}
// }

func stToS(st int) string {
	stateStrings := []string{"F", "I", "X", "T"}
	return stateStrings[st&0x3]
}

func dumpWires() {
	fmt.Println("Wire table:")
	for k := range wireTable {
		fmt.Printf("id: %d name: %s state: %s", k, wireTable[k].Name, stToS(wireTable[k].State))
		if len(wireTable[k].Outputs) > 0 {
			fmt.Printf("\n  connected to:")
			for i := range wireTable[k].Outputs {
				fmt.Printf(" %s", gateTable[i].Name)
			}
			fmt.Println()
		}
	}
	fmt.Println()
}

func dumpGates() {
	fmt.Println("\nGate table:")
	for k := range gateTable {
		fmt.Printf("id: %d name: %s\n", k, gateTable[k].Name)
		if len(gateTable[k].Inputs) > 0 {
			fmt.Printf("  inputs:")
			for i := range gateTable[k].Inputs {
				fmt.Printf(" %s", wireTable[i].Name)
			}
			fmt.Println()
		}
		if len(gateTable[k].Outputs) > 0 {
			fmt.Printf("  outputs:")
			for i := range gateTable[k].Outputs {
				fmt.Printf(" %s", wireTable[i].Name)
			}
			fmt.Println()
		}
	}
	fmt.Println()
}

var signalQueue map[int]bool
var gateQueue map[int]bool

var tick int

func setSignal(wire int, nextState int) {
	if nextState != wireTable[wire].State {
		// changed
		signalQueue[wire] = true
		wireTable[wire].State = nextState
		fmt.Printf("tick: %d signal: %s changed to: %s\n", tick, wireTable[wire].Name, stToS(wireTable[wire].State))
	}
}

func main() {

	wireTable = make(map[int]*Wire)
	gateTable = make(map[int]*Gate)

	A := newWire("A")
	B := newWire("B")
	C := newWire("C")

	D := newWire("D")

	G1 := newGate("AND0", AndFunc)

	G2 := newGate("NOT0", NotFunc)

	attach(G1, []int{A, B}, []int{C})
	attach(G2, []int{D}, []int{D})

	dumpWires()
	dumpGates()

	signalQueue = map[int]bool{}
	gateQueue = map[int]bool{}

	for tick = 0; tick < 10; tick++ {
		if tick == 0 {
			setSignal(A, stTrue)
			setSignal(D, stTrue)
		} else if tick == 1 {
			setSignal(B, stTrue)
		} else if tick == 2 {
			setSignal(A, stFalse)
		}

		// form gate queue
		gateQueue = map[int]bool{}
		for i := range signalQueue {
			for g := range wireTable[i].Outputs {
				gateQueue[g] = true
				fmt.Printf("tick: %d gate: %s touched by %s\n", tick, gateTable[g].Name, wireTable[i].Name)
			}
			fmt.Println()
		}

		// simulating gates
		signalQueue = map[int]bool{}
		for g := range gateQueue {
			nextState := gateTable[g].Func(gateTable[g].Inputs)
			if gateTable[g].State != nextState {
				fmt.Printf("tick: %d gate: %s changed to: %s\n", tick, gateTable[g].Name, stToS(nextState))
				for i := range gateTable[g].Outputs {
					setSignal(i, nextState)
				}
			}
			gateTable[g].State = nextState
			fmt.Println()
		}
	}
}
