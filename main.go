package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
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
	Inputs  []int // connected gate list
	Outputs []int // connected gate list
}

// Table storing all wires
var wireTable []Wire

// Gate element of gate table
type Gate struct {
	Name    string
	State   int
	Func    func([]int) int
	Inputs  []int // ids of input wires
	Outputs []int // ids of output wires
}

// Table storing all gates
var gateTable []Gate

func newWire(name string) (id int) {
	id = len(wireTable)
	if len(name) == 0 {
		name = fmt.Sprintf("_t%d", id)
	}
	wireTable = append(wireTable, Wire{
		Name:    name,
		State:   stUndefined,
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
	if out == stFalse {
		return stTrue
	} else if out == stTrue {
		return stFalse
	}
	return out
}

func newGate(name string, trfunc TrFunc) (id int) {
	id = len(gateTable)
	if len(name) == 0 {
		name = "_g"
	}
	gateTable = append(gateTable, Gate{
		Name:    fmt.Sprintf("%s%d", name, id),
		State:   stUndefined,
		Func:    trfunc,
		Inputs:  []int{},
		Outputs: []int{}})
	return id
}

// Attach gate input/outputs to wires
func attach(g int, inputs []int, outputs []int) {
	gateTable[g].Inputs = inputs
	gateTable[g].Outputs = outputs
	// for each input wire W add this gate to output fanout list
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

func stToS(st int) string {
	stateStrings := []string{"F", "I", "X", "T"}
	return stateStrings[st&0x3]
}

func dumpWires() {
	fmt.Println("Wire table:")
	for i, w := range wireTable {
		fmt.Printf("id: %d name: %s state: %s\n", i, w.Name, stToS(w.State))
		if len(w.Outputs) > 0 {
			fmt.Printf("  connected to:")
			for _, g := range w.Outputs {
				fmt.Printf(" %s", gateTable[g].Name)
			}
			fmt.Println()
		}
	}
	fmt.Println()
}

func dumpGates() {
	fmt.Println("Gate table:")
	for i, g := range gateTable {
		fmt.Printf("id: %d name: %s\n", i, g.Name)
		if len(g.Inputs) > 0 {
			fmt.Printf("  inputs:")
			for _, w := range g.Inputs {
				fmt.Printf(" %s", wireTable[w].Name)
			}
			fmt.Println()
		}
		if len(g.Outputs) > 0 {
			fmt.Printf("  outputs:")
			for _, w := range g.Outputs {
				fmt.Printf(" %s", wireTable[w].Name)
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
		if trace {
			fmt.Printf("tick: %d signal: %s changed to: %s\n", tick, wireTable[wire].Name, stToS(wireTable[wire].State))
		}
	}
}

var trace bool = true

var vcdDump bool = true
var vcdFileName string = "output.vcd"
var vcdOutFile *os.File

var dotDump bool = true
var dotFileName string = "schematic.dot"

func dumpDot() {
	out, err := os.OpenFile(dotFileName, os.O_CREATE|os.O_RDWR, 0644)
	defer out.Close()
	if err != nil {
		log.Fatalln("Cannot create dot file:", dotFileName, "err:", err)
	}

	fmt.Fprintln(out, "digraph top {\nrankdir=\"LR\"")
	for _, g := range gateTable {

		for _, w := range g.Inputs {
			fmt.Fprintf(out, "%s -> %s;\n", wireTable[w].Name, g.Name)
		}
		for _, w := range g.Outputs {
			fmt.Fprintf(out, "%s -> %s;\n", g.Name, wireTable[w].Name)
		}
		fmt.Fprintf(out, "%s [shape=box];\n", g.Name)
	}

	for _, w := range wireTable {
		if len(w.Inputs) == 0 {
			fmt.Fprintf(out, "%s [shape=circle];\n", w.Name)
		} else if len(w.Outputs) == 0 {
			fmt.Fprintf(out, "%s [shape=doublecircle];\n", w.Name)
		} else {
			fmt.Fprintf(out, "%s [shape=octagon];\n", w.Name)
		}
	}

	fmt.Fprintln(out, "}")
}

func printVCDHeader(out io.Writer) {
	fmt.Fprintln(out, "$date\n  ", time.Now().String(), "\n$end")
	fmt.Fprintln(out, "$version\n  Gogic logic simulator\n$end")
	fmt.Fprintln(out, "$timescale 1ns $end")
	fmt.Fprintln(out, "$scope module top $end")
	for _, w := range wireTable {
		fmt.Fprintf(out, "$var wire 1 w%s %s $end\n", w.Name, w.Name)
	}
	fmt.Fprintln(out, "$upscope $end\n$enddefinitions $end\n$dumpvars")
}

func vcdDumpVars(out io.Writer, tick int, sq map[int]bool) {
	sTov := []string{"0", "i", "x", "1"}
	fmt.Fprintf(out, "#%d\n", tick)
	for w := range sq {
		fmt.Fprintf(out, "%sw%s\n", sTov[wireTable[w].State], wireTable[w].Name)
	}
}

func newMUX2(a, d0, d1, y int) {
	t1 := newWire("")
	t2 := newWire("")
	t3 := newWire("")

	not := newGate("not", NotFunc)
	and1 := newGate("and", AndFunc)
	and2 := newGate("and", AndFunc)
	or := newGate("or", OrFunc)

	attach(not, []int{a}, []int{t1})
	attach(and1, []int{d0, t1}, []int{t2})
	attach(and2, []int{d1, a}, []int{t3})

	attach(or, []int{t2, t3}, []int{y})
}

// XOR element: c = a XOR b
func newXOR(a, b, c int) {
	t1 := newWire("")
	t2 := newWire("")
	t3 := newWire("")
	t4 := newWire("")

	not1 := newGate("not", NotFunc)
	not2 := newGate("not", NotFunc)

	and1 := newGate("and", AndFunc)
	and2 := newGate("and", AndFunc)

	or1 := newGate("or", OrFunc)

	attach(not1, []int{a}, []int{t1})
	attach(not2, []int{b}, []int{t2})

	attach(and1, []int{b, t1}, []int{t3})
	attach(and2, []int{a, t2}, []int{t4})

	attach(or1, []int{t3, t4}, []int{c})
}

// Half summator A+B -> SUM, CARRY
func newHalfSum(a, b, sum, carry int) {
	newXOR(a, b, sum)
	and := newGate("and", AndFunc)
	attach(and, []int{a, b}, []int{carry})
}

func newFullAdder(a, b int, ci, sum, co int) {
	t1 := newWire("")
	newXOR(a, b, t1)
	newXOR(ci, t1, sum)
	and1 := newGate("and", AndFunc)

	t2 := newWire("")
	attach(and1, []int{a, b}, []int{t2})

	t3 := newWire("")
	and2 := newGate("and", AndFunc)
	attach(and2, []int{ci, t1}, []int{t3})

	or := newGate("or", OrFunc)
	attach(or, []int{t2, t3}, []int{co})
}

func main() {
	wireTable = []Wire{}
	gateTable = []Gate{}

	A := newWire("A")
	D0 := newWire("D0")
	D1 := newWire("D1")

	Y := newWire("Y")

	newMUX2(A, D0, D1, Y)

	if trace {
		dumpWires()
		dumpGates()
	}

	if vcdDump {
		var err error
		vcdOutFile, err = os.OpenFile(vcdFileName, os.O_CREATE|os.O_RDWR, 0644)
		defer vcdOutFile.Close()
		if err != nil {
			log.Fatalln("Cannot create vcd file:", vcdFileName, "err:", err)
		}
		printVCDHeader(vcdOutFile)
	}

	if dotDump {
		dumpDot()
	}

	signalQueue = map[int]bool{}
	gateQueue = map[int]bool{}

	// main loop
	for tick = 0; tick < 20; tick++ {
		// generators section
		if tick == 0 {
			setSignal(A, stFalse)
			setSignal(D0, stFalse)
			setSignal(D1, stFalse)
			if trace {
				fmt.Println()
			}
		} else if tick == 3 {
			setSignal(D0, stTrue)
			if trace {
				fmt.Println()
			}
		} else if tick == 6 {
			setSignal(A, stTrue)
			if trace {
				fmt.Println()
			}
		}
		// dump vars to vcd
		if vcdDump {
			vcdDumpVars(vcdOutFile, tick, signalQueue)
		}

		// form gate queue
		gateQueue = map[int]bool{}
		for w := range signalQueue {
			for _, g := range wireTable[w].Outputs {
				gateQueue[g] = true
				if trace {
					fmt.Printf("tick: %d gate: %s touched by %s\n", tick, gateTable[g].Name, wireTable[w].Name)
				}
			}
		}
		if trace {
			if len(signalQueue) > 0 {
				fmt.Println()
			}
		}

		// simulating gates and form signals queue
		signalQueue = map[int]bool{}
		for g := range gateQueue {
			nextState := gateTable[g].Func(gateTable[g].Inputs)
			if gateTable[g].State != nextState {
				if trace {
					fmt.Printf("tick: %d gate: %s changed to: %s\n", tick, gateTable[g].Name, stToS(nextState))
				}
				for _, w := range gateTable[g].Outputs {
					setSignal(w, nextState)
				}
			}
			gateTable[g].State = nextState
		}
		if trace {
			if len(gateQueue) > 0 {
				fmt.Println()
			}
		}
	}
}
