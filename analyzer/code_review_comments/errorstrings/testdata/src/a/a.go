package a

import (
	"errors"
	"fmt"
)

func f() {
	e := fmt.Errorf("This is %s", "world") // want "gostyle.errorstrings"
	print(e.Error())
	e2 := fmt.Errorf("this is %s.", "world") // want "gostyle.errorstrings"
	print(e2.Error())
	var e3 = errors.New("This is world") // want "gostyle.errorstrings"
	print(e3.Error())
	var e4 = errors.New("this is world.") // want "gostyle.errorstrings"
	print(e4.Error())
}
