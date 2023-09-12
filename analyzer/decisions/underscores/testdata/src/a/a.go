package a

import (
	_ "embed"
	"log"
	"os"
)

const MAX_LENGTH = 10 // want "gostyle.underscores"

func f_a(i i_a) { // want "gostyle.underscores"
	var go_Pher int // want "gostyle.underscores"
	print(go_Pher)
	i.Foo_Bar()
	d_d, _ := os.ReadDir("tmp") // want "gostyle.underscores"
	d_d, _ = os.ReadDir("tmp")
	log.Println(d_d)
}

type T_a struct { // want "gostyle.underscores"
	foo_bar int //nolint:all
}

type i_a interface { //nostyle:all
	Foo_Bar() //nostyle:underscores
}

type S struct{}

func (s_a *S) Foo() {} // want "gostyle.underscores"
