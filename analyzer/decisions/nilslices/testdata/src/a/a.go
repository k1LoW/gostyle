package a

import "net/http"

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

	if s == nil { // want "gostyle.nilslices"
		print("nil")
	}

	sa := "a"
	if sa == "b" || s2 != nil { // want "gostyle.nilslices"
		print("nil")
	}

	if nil != s3 { // want "gostyle.nilslices"
		print("nil")
	}

	if len(s4) == 0 {
		print("nil")
	}

	var h *http.Client
	if h == nil {
		print("nil")
	}
}
