package a

const longlonglonglonglonglonglong = 99 // want "gostyle.varnames"

var hellohellohellohellohello int // want "gostyle.varnames"

func small() {
	for iiii := 0; iiii < 10; iiii++ { // want "gostyle.varnames"
		print(iiii)
	}

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for kkkkkkk, v := range m { // want "gostyle.varnames"
		print(kkkkkkk, v)
	}

	for k, vvvvvvvv := range m { // want "gostyle.varnames"
		print(k, vvvvvvvv)
	}
}

func medium() {
	var gopher int // want "gostyle.varnames"
	print(0)
	print(1)
	print(2)
	print(3)
	print(gopher)
}

func exlude() {
	var thisIsExludeVar int
	print(thisIsExludeVar)
}
