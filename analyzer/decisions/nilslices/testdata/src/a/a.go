package a

var s = []string{} // want "gostyle.nilslices"

var s2 = []string{"a", "b", "c"}

func f() {
	s := []string{} // want "gostyle.nilslices"
	print(s)
	var s2 []string
	print(s2)
	s3 := []string{"a", "b", "c"}
	print(s3)
	s4 := []string{} //nostyle:nilslices
	print(s4)
}
