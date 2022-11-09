package main

import (
	"fmt"
)

type Gate struct {
	Name  string
	State []bool
	Func  func([]bool) []bool
}

func (g *Gate) Eval(in []bool) (out []bool, changed bool) {
	nextState := g.Func(in)
	if len(g.State) != len(nextState) {
		changed = true
	} else {
		for i := range g.State {
			if nextState[i] != g.State[i] {
				changed = true
				break
			}
		}
	}
	copy(g.State, nextState)
	return nextState, changed
}

func BufferFunc(in []bool) (out []bool) {
	for _, v := range in {
		out = append(out, v)
	}
	return out
}

func NotFunc(in []bool) (out []bool) {
	for _, v := range in {
		out = append(out, !v)
	}
	return out
}

func OrFunc(in []bool) (out []bool) {
	res := in[0]
	for _, v := range in[1:] {
		res = res || v
	}
	return []bool{res}
}

func AndFunc(in []bool) (out []bool) {
	res := in[0]
	for _, v := range in[1:] {
		res = res && v
	}
	return []bool{res}
}

func main() {
	in := []bool{true, false}

	buf := Gate{Name: "BUF", Func: BufferFunc}
	not := Gate{Name: "NOT", Func: NotFunc}
	and := Gate{Name: "AND", Func: AndFunc}
	or := Gate{Name: "OR", Func: OrFunc}

	queue := []Gate{buf, not, and, or}
	for _, v := range queue {
		val, _ := v.Eval(in)
		fmt.Println("Gate:", v.Name, "in:", in, "out:", val)
	}
}
