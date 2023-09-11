package a

func f_a() { // want "The convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names."
	var go_Pher int // want "The convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names."
	print(go_Pher)  // want "The convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names."
}

type t_a struct { // want "The convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names."
	foo_bar int // want "The convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names."
}

type i_a interface { // want "The convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names."
	Foo_Bar() // want "The convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names."
}
