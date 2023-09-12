package a

import (
	"log"
	"os"

	_ "embed"
)

const MAX_LENGTH = 10 // want "gostyle.mixedcaps"

func f_a(i i_a) { // want "gostyle.mixedcaps"
	var go_Pher int // want "gostyle.mixedcaps"
	print(go_Pher)
	i.Foo_Bar()
	d_d, _ := os.ReadDir("tmp") // want "gostyle.mixedcaps"
	d_d, _ = os.ReadDir("tmp")
	log.Println(d_d)
}

type T_a struct { // want "gostyle.mixedcaps"
	foo_bar int //nolint:all
}

type i_a interface { //nostyle:all
	Foo_Bar() //nostyle:mixedcaps
}

type S struct{}

func (s_a *S) Foo() {} // want "gostyle.mixedcaps"
