package a

import "fmt"

func f() {
	var s = fmt.Sprintf("hello %s", "world")
	print(s)
	var s2 = fmt.Sprintf("this is '%s'", "world") // want "gostyle.useq"
	print(s2)
	var s3 = fmt.Sprintf(`this is "%s"`, "world") // want "gostyle.useq"
	print(s3)
	var s4 = fmt.Sprintf("this is \"%s\"", "world") // want "gostyle.useq"
	print(s4)
	var s5 = fmt.Sprintf("%s", "world")
	print(s5)
}
