package a

const MAX_LENGTH = 10 // want "MixedCaps"

func f_a() { // want "MixedCaps"
	var go_Pher int // want "MixedCaps"
	print(go_Pher)  // want "MixedCaps"
}

type T_a struct { // want "MixedCaps"
	foo_bar int //nolint:all
}

type i_a interface { //nostyle:all
	Foo_Bar() //nostyle:mixedcaps
}
