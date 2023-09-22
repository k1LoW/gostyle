package a

import (
	"fmt"
	"os"
	"testing"
)

func TestA(t *testing.T) {
	_, _ = fmt.Fprint(os.Stderr, "hello")
}
