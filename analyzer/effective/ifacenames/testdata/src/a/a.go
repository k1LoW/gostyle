package a

type Query interface { // want "By convention, one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun."
	Do() error
}

type Closer interface { // want "By convention, one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun."
	Do() error
}

type Writer interface {
	Write() error
}
