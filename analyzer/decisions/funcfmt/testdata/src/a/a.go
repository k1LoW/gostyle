package a

func b( // want "gostyle.funcfmt"
	a, b, c, d int,
	e, f,
	g string) error {
	return nil
}

func c() error {
	return b(1, 2, 3, 4, // want "gostyle.funcfmt"
		"", "", "")
}
