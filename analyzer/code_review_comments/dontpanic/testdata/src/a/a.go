package a

import (
	"errors"
)

func f() {
	if true {
		panic(errors.New("error")) // want "gostyle.dontpanic"
	}
}
