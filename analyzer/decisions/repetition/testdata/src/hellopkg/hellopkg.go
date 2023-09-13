package hellopkg

import "fmt"

const HelloTitle = "HELLO" // want "gostyle.repetition"

var GoHello = "Go" // want "gostyle.repetition"

var helloStr = "Hello, World!" // want "gostyle.repetition"

func HelloWorld() string { // want "gostyle.repetition"
	return helloStr
}

func MyHello() string { // want "gostyle.repetition"
	return helloStr
}

func HelloMe() string { //nostyle:repetition
	var tenInt = 10 // want "gostyle.repetition"
	return fmt.Sprintf("Hello %d", tenInt)
}
