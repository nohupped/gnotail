package generic

import (
	"os"
	"io"
)
// Opener opens the path and returns an io.Reader.
func Opener(filename string) (io.ReadCloser, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	return f, err
}