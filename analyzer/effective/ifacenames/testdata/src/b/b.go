package b

type Query interface { // want "-er suffix"
	Do() error
}

type Closer interface { // want "-er suffix"
	Do() error
}

type Writer interface {
	Write() error
}

type Add interface { // want "-er suffix"
	One() error
	Two() error
}
