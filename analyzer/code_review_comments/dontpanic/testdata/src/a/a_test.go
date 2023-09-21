package a

import (
	"errors"
	"testing"
)

func TestA(t *testing.T) {
	if true {
		panic(errors.New("error"))
	}
}
