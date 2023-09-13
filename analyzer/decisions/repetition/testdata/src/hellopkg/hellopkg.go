package hellopkg

const HelloTitle = "HELLO" // want "gostyle.repetition"

var GoHello = "Go" // want "gostyle.repetition"

var helloStr = "Hello, World!"

func HelloWorld() string { // want "gostyle.repetition"
	return helloStr
}

func MyHello() string { // want "gostyle.repetition"
	return helloStr
}

func HelloMe() string { //nostyle:repetition
	return helloStr
}
