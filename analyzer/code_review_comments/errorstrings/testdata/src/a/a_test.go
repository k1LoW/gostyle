package a

import "testing"

func TestA(t *testing.T) {
	t.Errorf("This is %s", "world")
}
