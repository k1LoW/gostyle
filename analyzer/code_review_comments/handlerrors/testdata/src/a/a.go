package a

import (
	"fmt"
	"os"
)

func f() {
	_, _ = fmt.Fprint(os.Stderr, "hello") // want "gostyle.handlerrors"
}
