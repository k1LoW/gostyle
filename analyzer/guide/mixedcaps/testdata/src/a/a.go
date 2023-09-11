package a

const MAX_LENGTH = 10 // want "MixedCaps"

func f_a() { // want "MixedCaps"
	var go_Pher int // want "MixedCaps"
	print(go_Pher)  // want "MixedCaps"
}

type t_a struct { // want "MixedCaps"
	foo_bar int // want "MixedCaps"
}

type i_a interface { // want "MixedCaps"
	Foo_Bar() // want "MixedCaps"
}
