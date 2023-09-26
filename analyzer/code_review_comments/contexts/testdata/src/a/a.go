package a

import (
	"context"
)

type A struct {
	v   int
	ctx context.Context // want "gostyle.contexts"
}

func f(hello string, ctx context.Context) { // want "gostyle.contexts"
	b := func(a int, ctx context.Context) { // want "gostyle.contexts"
		println(a)
	}
	b(1, ctx)
	println(hello)
}
