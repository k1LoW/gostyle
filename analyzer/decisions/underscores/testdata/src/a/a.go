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

	for i_i := 0; i_i < 10; i_i++ { // want "gostyle.underscores"
		print(i_i)
	}

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k_k, v := range m { // want "gostyle.underscores"
		print(k_k, v)
	}
	for k, v_v := range m { // want "gostyle.underscores"
		print(k, v_v)
	}
}

type T_a struct { // want "gostyle.underscores"
	foo_bar int //nolint:all
}

type i_a interface { //nostyle:all
	Foo_Bar() //nostyle:underscores
}

type S struct{}

func (s_a *S) Foo() {} // want "gostyle.underscores"
