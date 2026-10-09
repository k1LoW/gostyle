package a

type UnnamedMap map[string]string

func (*UnnamedMap) Reset() { // want "gostyle.recvtype\\] If the receiver is a map.*: UnnamedMap$"
}

type UnnamedCh chan int

func (*UnnamedCh) Reset() { // want "gostyle.recvtype\\] If the receiver is a map.*: UnnamedCh$"
}

type UnnamedFn func()

func (*UnnamedFn) Reset() { // want "gostyle.recvtype\\] If the receiver is a map.*: UnnamedFn$"
}

type NamedMap map[string]string

func (m *NamedMap) Reset() { // want "gostyle.recvtype\\] If the receiver is a map.*: m$"
}
