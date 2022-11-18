package sim

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var trace bool

func SetTrace() {
	trace = true
}

func DumpDot(dotFileName string) {
	out, err := os.OpenFile(dotFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln("Cannot create dot file:", dotFileName, "err:", err)
	}
	defer out.Close()

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

var vcdDump bool
var vcdOutFile *os.File

func SetVCDOutput(name string) {
	vcdDump = true
	var err error
	vcdOutFile, err = os.OpenFile(name, os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln("Cannot create vcd file:", name, "err:", err)
	}
	writeVCDHeader(vcdOutFile)
}

func writeVCDHeader(out io.Writer) {
	fmt.Fprintln(out, "$date\n  ", time.Now().String(), "\n$end")
	fmt.Fprintln(out, "$version\n  Gogic logic simulator\n$end")
	fmt.Fprintln(out, "$timescale 1ns $end")
	fmt.Fprintln(out, "$scope module top $end")
	for _, w := range wireTable {
		fmt.Fprintf(out, "$var wire 1 w%s %s $end\n", w.Name, w.Name)
	}
	fmt.Fprintln(out, "$upscope $end\n$enddefinitions $end\n$dumpvars")
}

func vcdDumpVars(tick int, sq map[int]bool) {
	sTov := []string{"0", "i", "x", "1"}
	fmt.Fprintf(vcdOutFile, "#%d\n", tick)
	for w := range sq {
		fmt.Fprintf(vcdOutFile, "%sw%s\n", sTov[wireTable[w].State], wireTable[w].Name)
	}
}

func StToS(st int) string {
	stateStrings := []string{"F", "I", "X", "T"}
	return stateStrings[st&0x3]
}

func DumpWires() {
	fmt.Println("Wire table:")
	for i, w := range wireTable {
		fmt.Printf("id: %d name: %s state: %s\n", i, w.Name, StToS(w.State))
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

func DumpGates() {
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
