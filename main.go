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

var wireTable []Wire

func newWire(name string) (id int) {
	id = len(wireTable)
	wireTable = append(wireTable, Wire{
		Name:    name,
		State:   stUndefined,
		Outputs: []int{}})
	return id
}

// TrFunc gate transmission function
type TrFunc func([]int) int

func AndFunc(in []int) (out int) {
	out = wireTable[in[0]].State
	for i := range in[1:] {
		out &= wireTable[i].State
	}
	return out
}

func NotFunc(in []int) (out int) {
	out = wireTable[in[0]].State
	if out == stFalse {
		return stTrue
	} else if out == stTrue {
		return stFalse
	}
	return out
}

var gateTable []Gate

func newGate(name string, trfunc TrFunc) (id int) {
	id = len(gateTable)
	gateTable = append(gateTable, Gate{
		Name:    name,
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
			if wireTable[w].Outputs[j] == w {
				found = true
			}
		}
		if !found {
			wireTable[w].Outputs = append(wireTable[w].Outputs, g)
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

var trace bool = false
var vcdDump bool = true
var vcdFileName string = "output.vcd"
var vcdOutFile *os.File

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

func main() {
	wireTable = []Wire{}
	gateTable = []Gate{}

	A := newWire("A")
	B := newWire("B")
	C := newWire("C")

	D := newWire("D")

	G1 := newGate("AND0", AndFunc)

	G2 := newGate("NOT0", NotFunc)

	attach(G1, []int{A, B}, []int{C})
	attach(G2, []int{D}, []int{D})

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

	signalQueue = map[int]bool{}
	gateQueue = map[int]bool{}

	for tick = 0; tick < 10; tick++ {
		// generators section
		if tick == 0 {
			setSignal(A, stTrue)
			setSignal(D, stTrue)
			if trace {
				fmt.Println()
			}
		} else if tick == 1 {
			setSignal(B, stTrue)
			if trace {
				fmt.Println()
			}
		} else if tick == 2 {
			setSignal(A, stFalse)
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
