package a

type Header map[string][]string

func (h *Header) Add(key, value string) { // want "gostyle.recvtype"
	(*h)[key] = append((*h)[key], value)
}

type Ch chan int

func (ch *Ch) Close() { // want "gostyle.recvtype"
	close(*ch)
}

type Fn func()

func (fn *Fn) Call() { // want "gostyle.recvtype"
	(*fn)()
}

type S struct{}

func (s S) M() string { // want "gostyle.recvtype"
	return "foo"
}
