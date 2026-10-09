package a

import (
	"context"
	"io"
)

type EmbeddedContext struct {
	context.Context // want "gostyle.contexts\\] Don't add a Context member.*: Context$"
}

type NamedContext struct {
	request context.Context // want "gostyle.contexts\\] Don't add a Context member.*: request$"
}

type MultipleContexts struct {
	first, second context.Context // want "gostyle.contexts\\] Don't add a Context member.*: first$"
}

type EmbeddedReader struct {
	io.Reader
}

type EmbeddedLocal struct {
	EmbeddedReader
}

func contextFirst(ctx context.Context, value string) {}

func contextSecond(value string, ctx context.Context) { // want "gostyle.contexts\\] Most functions.*: contextSecond$"
}
