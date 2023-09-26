package a

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestA(t *testing.T) {
	ctx := context.Background()
	b := func(a int, ctx context.Context) {
		println(a)
	}
	b(1, ctx)
	_, _ = fmt.Fprint(os.Stderr, "hello")
}
