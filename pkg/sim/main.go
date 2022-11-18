package sim

import (
	"fmt"
)

// Signal states
const (
	StFalse     = 0 // 0
	StInvalid   = 1
	StUndefined = 2 // X
	StTrue      = 3 // 1
)

// Table storing all wires
var wireTable []Wire

// Table storing all gates
var gateTable []Gate

var signalQueue map[int]bool
var gateQueue map[int]bool

var tick int

func SetSignal(wire int, nextState int) {
	if nextState != wireTable[wire].State {
		// changed
		signalQueue[wire] = true
		wireTable[wire].State = nextState
		if trace {
			fmt.Printf("tick: %d signal: %s changed to: %s\n", tick, wireTable[wire].Name, StToS(wireTable[wire].State))
		}
	}
}

func Init() {
	wireTable = []Wire{}
	gateTable = []Gate{}
}

func SimStep() {
	// dump vars to vcd
	if vcdDump {
		vcdDumpVars(tick, signalQueue)
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
				fmt.Printf("tick: %d gate: %s changed to: %s\n", tick, gateTable[g].Name, StToS(nextState))
			}
			for _, w := range gateTable[g].Outputs {
				SetSignal(w, nextState)
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

type StimulusFunc func(int)

func Simulate(ticksMax int, stimulus StimulusFunc) {
	signalQueue = map[int]bool{}
	gateQueue = map[int]bool{}

	// main loop
	for tick = 0; tick < ticksMax; tick++ {
		// call stimuli
		stimulus(tick)
		SimStep()
	}
}
