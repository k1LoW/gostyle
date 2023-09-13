package a_test

import "testing"

func TestA(t *testing.T) {
	t.Error(1)
}

func Test_B(t *testing.T) {
	t.Error(1)
}
