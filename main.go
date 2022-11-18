package main

import "github.com/x86128/gogic/pkg/sim"

func main() {
	// init system
	sim.Init()

	// build scheme
	A := sim.NewWire("A")
	D0 := sim.NewWire("D0")
	D1 := sim.NewWire("D1")

	Y := sim.NewWire("Y")

	sim.NewMUX2(A, D0, D1, Y)

	sim.SetVCDOutput("output.vcd")
	sim.DumpWires()
	sim.DumpGates()
	sim.DumpDot("schematic.dot")

	// simulate 10 ticks
	sim.SetTrace()
	sim.Simulate(10, func(tick int) {
		// stimuli
		if tick == 0 {
			sim.SetSignal(A, sim.StFalse)
			sim.SetSignal(D0, sim.StFalse)
			sim.SetSignal(D1, sim.StFalse)
		} else if tick == 3 {
			sim.SetSignal(D0, sim.StTrue)
		} else if tick == 6 {
			sim.SetSignal(A, sim.StTrue)
		}
	})
}
