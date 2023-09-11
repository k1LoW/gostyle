package b

type Query interface { // want "gostyle.ifacenames"
	Do() error
}

type Closer interface { // want "gostyle.ifacenames"
	Do() error
}

type Writer interface {
	Write() error
}

type Add interface { // want "gostyle.ifacenames"
	One() error
	Two() error
}

type Sub interface { //nostyle:ifacenames
	One() error
	Two() error
}
