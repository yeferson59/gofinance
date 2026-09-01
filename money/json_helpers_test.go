package money

import (
	"bytes"
	"io"
)

// bytesReader is the reader jsontext decoders are built over in these tests.
func bytesReader(s string) io.Reader { return bytes.NewReader([]byte(s)) }

// sliceWriter collects what an encoder writes so a test can assert on the
// whole document rather than on one value.
type sliceWriter struct{ out *[]byte }

func (w sliceWriter) Write(p []byte) (int, error) {
	*w.out = append(*w.out, p...)

	return len(p), nil
}
