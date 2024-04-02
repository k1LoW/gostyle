package a

type Query interface { // want "gostyle.ifacenames"
	Do() error
}

type Closer interface { // want "gostyle.ifacenames"
	Do() error
}

type Writer interface {
	Write() error
}

type Validator interface {
	Validate() error
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
