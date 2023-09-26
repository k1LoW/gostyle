package a

func b(a, b, c, d int, // want "gostyle.funcfmt"
	e, f,
	g string) error {
	return nil
}

func c() error {
	return b(1, 2, 3, 4, // want "gostyle.funcfmt"
		"", "", "")
}
