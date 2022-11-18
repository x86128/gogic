package sim

func NewMUX2(a, d0, d1, y int) {
	t1 := NewWire("")
	t2 := NewWire("")
	t3 := NewWire("")

	NewGate("not", NotFunc, []int{a}, []int{t1})
	NewGate("and", AndFunc, []int{d0, t1}, []int{t2})
	NewGate("and", AndFunc, []int{d1, a}, []int{t3})
	NewGate("or", OrFunc, []int{t2, t3}, []int{y})
}

// XOR element: c = a XOR b
func NewXOR(a, b, c int) {
	t1 := NewWire("")
	t2 := NewWire("")
	t3 := NewWire("")
	t4 := NewWire("")

	NewGate("not", NotFunc, []int{a}, []int{t1})
	NewGate("not", NotFunc, []int{b}, []int{t2})
	NewGate("and", AndFunc, []int{b, t1}, []int{t3})
	NewGate("and", AndFunc, []int{a, t2}, []int{t4})
	NewGate("or", OrFunc, []int{t3, t4}, []int{c})
}

// Half adder A+B -> SUM, CARRY
func NewHalfAdder(a, b, sum, carry int) {
	NewXOR(a, b, sum)
	NewGate("and", AndFunc, []int{a, b}, []int{carry})
}

// Full Adder
func NewFullAdder(a, b int, ci, sum, co int) {
	t1 := NewWire("")
	NewXOR(a, b, t1)
	NewXOR(ci, t1, sum)
	t2 := NewWire("")
	NewGate("and", AndFunc, []int{a, b}, []int{t2})
	t3 := NewWire("")
	NewGate("and", AndFunc, []int{ci, t1}, []int{t3})
	NewGate("or", OrFunc, []int{t2, t3}, []int{co})
}
