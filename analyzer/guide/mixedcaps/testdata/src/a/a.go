package a

const MAX_LENGTH = 10 // want "gostyle.mixedcaps"

func f_a() { // want "gostyle.mixedcaps"
	var go_Pher int // want "gostyle.mixedcaps"
	print(go_Pher)  // want "gostyle.mixedcaps"
}

type T_a struct { // want "gostyle.mixedcaps"
	foo_bar int //nolint:all
}

type i_a interface { //nostyle:all
	Foo_Bar() //nostyle:mixedcaps
}
