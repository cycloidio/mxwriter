package mxwriter

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

// Separator is the key used to separate when writing
// the length of the key from the content
var Separator []byte = []byte(";")

// mux is the internal implementation of the multiplexer
type mux struct {
	buffers map[string]*bytes.Buffer
	keys    []string
}

// NewMux returns a Multiplexer io.ReadWriter which
// can be used to write different streams of data
// and later on read them separately by key.
// To read from an specific key use the NewDemux
// function to get the Demux
func NewMux() io.ReadWriter {
	return &mux{
		buffers: make(map[string]*bytes.Buffer),
	}
}

// Write writes the content of p to the internal buffer.
// The format has to be `<length-key>;<key><content>`
func (w *mux) Write(p []byte) (int, error) {
	lenidx := bytes.Index(p, Separator)
	if lenidx == -1 {
		return 0, errors.New("invalid write format")
	}
	l, err := strconv.Atoi(string(p[:lenidx]))
	if err != nil {
		return 0, err
	}
	key := string(p[lenidx+1 : lenidx+l+1])
	p = p[lenidx+l+1:]

	if buff, ok := w.buffers[key]; ok {
		n, err := io.Copy(buff, bytes.NewReader(p))
		return int(n), err
	}

	w.buffers[key] = bytes.NewBuffer(p)
	w.keys = append(w.keys, key)
	return len(p), nil
}

// Read will basically read everything written, the
// order of content is based on the order of the Keys
// written using w.Write
func (w *mux) Read(p []byte) (int, error) {
	readers := make([]io.Reader, 0)
	for _, k := range w.keys {
		readers = append(readers, w.buffers[k])
	}

	return io.MultiReader(readers...).Read(p)
}

// Write is a helper that automatically writes with the expected format to the w
func Write(w io.Writer, key string, content []byte) (int, error) {
	l := len(key)
	content = append(bytes.Join([][]byte{[]byte(strconv.Itoa(l)), Separator, []byte(key)}, nil), content...)
	return w.Write(content)
}
