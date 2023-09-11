package a

type Query interface { // want "-er suffix"
	Do() error
}

type Closer interface { // want "-er suffix"
	Do() error
}

type Writer interface {
	Write() error
}

type Add interface {
	One() error
	Two() error
}

type ReadCloser interface { //nolint:all
	Do() error
}

type WriteCloser interface { //nostyle:ifacenames
	Do() error
}
