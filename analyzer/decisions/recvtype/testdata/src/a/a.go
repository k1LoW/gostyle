package a

type Header map[string][]string

func (h *Header) Add(key, value string) { // want "gostyle.recvtype"
	(*h)[key] = append((*h)[key], value)
}

func (h Header) Del(key, value string) {
	delete(h, key)
}

type Ch chan int

func (ch *Ch) Close() { // want "gostyle.recvtype"
	close(*ch)
}

func (ch Ch) Done() {
	<-ch
}

type Fn func()

func (fn *Fn) Call() { // want "gostyle.recvtype"
	(*fn)()
}

func (fn Fn) Fn() {
	fn()
}

type S struct{}

func (s S) M() string { // want "gostyle.recvtype"
	return "foo"
}
