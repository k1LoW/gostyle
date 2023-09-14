package a

func getSome() { // want "gostyle.getters"
	getFoo := func() {} // want "gostyle.getters"
	getFoo()

	var getBar = func() {} // want "gostyle.getters"
	getBar()

	getBaz := "getget"
	println(getBaz)
}

type Getted struct {
	geted int
}

type Getter interface {
	GetSomething() //  want "gostyle.getters"
}

type S struct{}

func (s *S) GetS() {} // want "gostyle.getters"
