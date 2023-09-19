package a

var Hello interface{} // want "gostyle.useany"

func f(b interface{}) { // want "gostyle.useany"
	print(b)
}

func f2() (interface{}, any, any) { // want "gostyle.useany"

	var s [][][]interface{}               // want "gostyle.useany"
	s2 := []map[int][]interface{}{}       // want "gostyle.useany"
	m := map[string]map[int]interface{}{} // want "gostyle.useany"

	return m, s, s2
}

func f3() (any, interface{}) { // want "gostyle.useany"
	s := []map[int][]any{}
	var s2 [][][]any
	return s, s2
}

func f4() any {
	return nil
}
